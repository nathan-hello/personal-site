package auth

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nathan-hello/personal-site/db"
	"golang.org/x/crypto/bcrypt"
)

func dbInsertNewToken(t string, jwt_type string) error {
	_, claims, err := ParseToken(t)
	if err != nil {
		return err
	}
	ctx := context.Background()

	token, err := db.Db().InsertToken(ctx, db.InsertTokenParams{
		JwtType:   jwt_type,
		Jwt:       t,
		Valid:     true,
		Family:    claims.Family,
		ExpiresAt: time.Now().Unix(),
	})

	if err != nil {
		return ErrDbInsertToken
	}

	err = db.Db().InsertUsersTokens(ctx, db.InsertUsersTokensParams{
		UserID:  claims.UserId,
		TokenID: token.ID,
	})
	if err != nil {
		return ErrDbInsertUsersToken
	}

	return nil
}

func dbValidateJwt(t string, user string) (*db.Profile, error) {
	ctx := context.Background()

	token, err := db.Db().SelectTokenFromJwtString(ctx, t)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrJwtNotInDb
		}
		return nil, ErrDbSelectUserFromJwt
	}
	if !token.Valid {
		return nil, ErrJwtInvalidInDb
	}

	userObj, err := db.Db().SelectUserProfileById(context.Background(), user)

	if err != nil {
		return nil, err
	}

	return &userObj, nil
}

func dbInvalidateJwtFamily(t string) error {

	ctx := context.Background()

	token, err := db.Db().SelectTokenFromJwtString(ctx, t)
	if err != nil {
		return ErrDbSelectUserFromJwt
	}

	_, claims, err := ParseToken(token.Jwt)
	if err != nil {
		return err
	}
	err = db.Db().UpdateTokensFamilyInvalid(ctx, claims.Family)
	if err != nil {
		return ErrDbUpdateTokensInvalid
	}

	return nil
}

func ForceLogout(ctx context.Context, userId string) error {
	tokens, err := db.Db().SelectUserActiveTokens(ctx, userId)
	if err != nil {
		return err
	}

	for _, token := range tokens {
		err = db.Db().UpdateTokenInvalid(ctx, token.ID)
	}
	return nil
}

// HandlePasswordChange invalidates all sessions for a user when their password changes
func HandlePasswordChange(ctx context.Context, userId string, password string) error {
	err := ForceLogout(ctx, userId)
	if err != nil {
		return err
	}

	salt := uuid.NewString()[:8]
	pass, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = db.Db().UpdateAuthUserPassword(ctx, db.UpdateAuthUserPasswordParams{
		EncryptedPassword: string(pass),
		PasswordSalt:      salt,
		ID:                userId,
	})
	if err != nil {
		return err
	}

	return nil
}

// ListUserSessions returns all active sessions for a user
func ListUserSessions(ctx context.Context, userId string) ([]*CustomClaims, error) {
	tokens, err := db.Db().SelectUserActiveTokens(ctx, userId)
	if err != nil {
		return nil, err
	}

	var sessions []*CustomClaims
	for _, token := range tokens {
		_, claims, err := ParseToken(token.Jwt)
		if err != nil {
			log.Printf("failed to parse token: %v", err)
			continue
		}
		sessions = append(sessions, claims)
	}

	return sessions, nil
}
