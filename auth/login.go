package auth

import (
	"context"
	"database/sql"
	"log"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nathan-hello/personal-site/db"
	"github.com/nathan-hello/personal-site/utils"
	"golang.org/x/crypto/bcrypt"
)

// Username validation rules
const (
	MinUsernameLength = 3
	MaxUsernameLength = 30
	UsernamePattern   = `^[a-zA-Z0-9_-]+$` // Only alphanumeric, underscore, and hyphen allowed
)

var usernameRegex = regexp.MustCompile(UsernamePattern)

// sanitizeInput removes potentially dangerous characters and normalizes input
func sanitizeInput(input string) string {
	// Remove HTML tags and special characters
	input = strings.TrimSpace(input)
	input = strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1 // Remove control characters
		}
		return r
	}, input)
	return input
}

func validateUsername(username string) []AuthError {
	var errs []AuthError

	username = sanitizeInput(username)

	if len(username) < MinUsernameLength {
		errs = append(errs, ErrUsernameTooShort)
	}
	if len(username) > MaxUsernameLength {
		errs = append(errs, ErrUsernameTooLong)
	}
	if !usernameRegex.MatchString(username) {
		errs = append(errs, ErrUsernameInvalidFormat)
	}

	return errs
}

type AuthResult struct {
	User   db.Profile
	Errors []AuthError
}

type SignUpData struct {
	Email    string
	Username string
	Password string
	PassConf string
}

type SignInData struct {
	UserOrEmail string
	Password    string
}

func passwordCheck(p string) bool {
	return len(p) > 7
}

func validateSignUpStrings(data SignUpData) []AuthError {
	var errs []AuthError

	// Sanitize inputs
	data.Username = sanitizeInput(data.Username)
	data.Email = sanitizeInput(data.Email)

	// Validate username
	usernameErrs := validateUsername(data.Username)
	errs = append(errs, usernameErrs...)

	// Validate email if provided
	if data.Email != "" {
		if _, err := mail.ParseAddress(data.Email); err != nil {
			errs = append(errs, ErrEmailInvalid)
		}
	}

	if data.Username == "" && data.Email == "" {
		errs = append(errs, ErrEmailOrUsernameReq)
	}

	if !passwordCheck(data.Password) {
		errs = append(errs, ErrPasswordInvalid)
	}
	if data.Password != data.PassConf {
		errs = append(errs, ErrPassNoMatch)
	}

	return errs
}

func SignUp(data SignUpData) AuthResult {
	errs := validateSignUpStrings(data)
	if len(errs) > 0 {
		return AuthResult{Errors: errs}
	}

	ctx := context.Background()

	// Check if email is taken
	if data.Email != "" {
		_, err := db.Db().SelectAuthUserByEmail(ctx, data.Email)
		if err != sql.ErrNoRows {
			return AuthResult{Errors: []AuthError{ErrEmailTaken}}
		}
	}

	// Check if username is taken
	if data.Username != "" {
		_, err := db.Db().SelectUserProfileByUsername(ctx, data.Username)
		if err != sql.ErrNoRows {
			return AuthResult{Errors: []AuthError{ErrUsernameTaken}}
		}
	}

	userId := uuid.NewString()
	salt := uuid.NewString()[:8]
	pass, err := bcrypt.GenerateFromPassword([]byte(data.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		return AuthResult{Errors: []AuthError{{
			Code:     ErrHashPassword.Code,
			Message:  ErrHashPassword.Message,
			Category: "system",
			Err:      err,
		}}}
	}

	err = db.Db().InsertAuthUser(
		ctx,
		db.InsertAuthUserParams{
			ID:                userId,
			Email:             data.Email,
			EncryptedPassword: string(pass),
			PasswordSalt:      salt,
			PasswordCreatedAt: time.Now(),
		})

	if err != nil {
		log.Printf("failed to insert user: %v", err)
		return AuthResult{Errors: []AuthError{{
			Code:     ErrDbInsertUser.Code,
			Message:  ErrDbInsertUser.Message,
			Category: "system",
			Err:      err,
		}}}
	}

	newUser, err := db.Db().SelectUserProfileById(ctx, userId)
	if err != nil {
		log.Printf("failed to retrieve user after insert: %v", err)
		return AuthResult{Errors: []AuthError{{
			Code:     ErrDbSelectAfterInsert.Code,
			Message:  ErrDbSelectAfterInsert.Message,
			Category: "system",
			Err:      err,
		}}}
	}

	return AuthResult{
		User: newUser,
	}
}

func SignIn(data SignInData) AuthResult {
	// Sanitize inputs
	data.UserOrEmail = sanitizeInput(data.UserOrEmail)
	data.Password = sanitizeInput(data.Password)

	if data.UserOrEmail == "" || data.Password == "" {
		return AuthResult{Errors: []AuthError{ErrBadLogin}}
	}

	var user db.AuthUser
	ctx := context.Background()

	// Try to find user by email or username
	_, err := mail.ParseAddress(data.UserOrEmail)
	if err == nil {
		user, err = db.Db().SelectAuthUserByEmail(ctx, data.UserOrEmail)
		if err != nil {
			return AuthResult{Errors: []AuthError{ErrBadLogin}}
		}
	} else {
		user, err = db.Db().SelectAuthUserByEmail(ctx, data.UserOrEmail)
		if err != nil {
			return AuthResult{Errors: []AuthError{ErrBadLogin}}
		}
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(data.Password+user.PasswordSalt))
	if err != nil {
		return AuthResult{Errors: []AuthError{ErrBadLogin}}
	}

	// Get user without password
	dbUser, err := db.Db().SelectUserProfileById(ctx, user.ID)
	if err != nil {
		log.Printf("failed to retrieve user after login: %v", err)
		return AuthResult{Errors: []AuthError{{
			Code:     ErrDbSelectAfterInsert.Code,
			Message:  ErrDbSelectAfterInsert.Message,
			Category: "system",
			Err:      err,
		}}}
	}

	accessToken, refreshToken, err := NewTokenPair(&JwtParams{
		UserId:   user.ID,
		Username: dbUser.Username,
		Family:   uuid.NewString(),
	})
	if err != nil {
		return AuthResult{Errors: []AuthError{ErrBadLogin}}
	}

	accessTokenDb, err := db.Db().InsertToken(ctx, db.InsertTokenParams{
		JwtType:   "access",
		Jwt:       accessToken,
		Valid:     true,
		Family:    uuid.NewString(),
		ExpiresAt: time.Now().Add(utils.Env().ACCESS_EXPIRY_TIME).Unix(),
	})
	if err != nil {
		return AuthResult{Errors: []AuthError{ErrBadLogin}}
	}

	refreshTokenDb, err := db.Db().InsertToken(ctx, db.InsertTokenParams{
		JwtType:   "refresh",
		Jwt:       refreshToken,
		Valid:     true,
		Family:    uuid.NewString(),
		ExpiresAt: time.Now().Add(utils.Env().REFRESH_EXPIRY_TIME).Unix(),
	})
	if err != nil {
		return AuthResult{Errors: []AuthError{ErrBadLogin}}
	}

	err = db.Db().InsertUsersTokens(ctx, db.InsertUsersTokensParams{
		UserID:  user.ID,
		TokenID: accessTokenDb.ID,
	})
	if err != nil {
		return AuthResult{Errors: []AuthError{ErrBadLogin}}
	}

	err = db.Db().InsertUsersTokens(ctx, db.InsertUsersTokensParams{
		UserID:  user.ID,
		TokenID: refreshTokenDb.ID,
	})
	if err != nil {
		return AuthResult{Errors: []AuthError{ErrBadLogin}}
	}

	return AuthResult{
		User: dbUser,
	}
}
