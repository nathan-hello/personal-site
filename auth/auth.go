package auth

import (
	"context"
	"database/sql"
	"log"

	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/nathan-hello/personal-site/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                string
	Username          string
        Email string
	GlobalChatColor   string
}

func passwordCheck(p string) bool {
        return len(p) > 7
}

type AuthHandler interface {
	RenderErrs() []string
}

type SignUp struct {
	Email       string
	Username    string
	Password    string
	PassConf    string
	UsernameErr string
	EmailErr    string
	PassErr     string
	PassConfErr string
	MiscErrs    []string
}

func (a *SignUp) RenderErrs() []string {
	errs := []string{a.UsernameErr, a.EmailErr, a.PassErr, a.PassConfErr}
	errs = append(errs, a.MiscErrs...)
	rendered := []string{}
	for _, v := range errs {
		if v != "" {
			rendered = append(rendered, v)
		}
	}
	return rendered
}

func (a *SignUp) FlushPasswords() {
	a.Password = ""
	a.PassConf = ""
}

func (a *SignUp) validateStrings() bool {
	if len(a.Username) < 3 {
			a.UsernameErr = ErrUsernameTooShort.Error()
	}

	if a.Username == "" && a.Email == "" {
		a.MiscErrs = append(a.MiscErrs, ErrEmailOrUsernameReq.Error())
	}

	if !passwordCheck(a.Password) {
		a.PassErr = ErrPasswordInvalid.Error()
	}
	if a.Password != a.PassConf {
		a.PassConfErr = ErrPassNoMatch.Error()
	}

	ok := len(a.RenderErrs()) == 0
	if !ok {
		a.FlushPasswords()
		return false
	}

	return true
}

func (a *SignUp) SignUp() *User {
	ok := a.validateStrings()
	if !ok {
		return nil
	}
	ctx := context.Background()

	if a.Email != "" {
		_, err := db.Db().SelectUserByEmail(ctx, a.Email)
		if err != sql.ErrNoRows {
			a.EmailErr = ErrEmailTaken.Error()
		}
	}
	if a.Username != "" {
		_, err := db.Db().SelectUserByUsername(ctx, a.Username)
		if err != sql.ErrNoRows {
			a.UsernameErr = ErrUsernameTaken.Error()
		}
	}
	ok = len(a.RenderErrs()) == 0
	if !ok {
		a.FlushPasswords()
		return nil
	}

	userId := uuid.NewString()
	salt := uuid.NewString()[:8]
	pass, err := bcrypt.GenerateFromPassword([]byte(a.Password+salt), bcrypt.DefaultCost)

	if err != nil {
		log.Println(err)
		a.MiscErrs = append(a.MiscErrs, ErrHashPassword.Error())
		return nil
	}

	err = db.Db().InsertUser(
		ctx,
		db.InsertUserParams{
			GlobalChatColor:   "text-gray-500",
			ID:                userId,
			Email:             a.Email,
			Username:          a.Username,
			EncryptedPassword: string(pass),
			PasswordSalt:      salt,
			PasswordCreatedAt: time.Now(),
		})

	if err != nil {
		log.Println(err)
		a.MiscErrs = append(a.MiscErrs, ErrDbInsertUser.Error())
		return nil
	}

        newUser, err := db.Db().SelectUserById(ctx, userId)
	if err != nil {
		log.Println(err)
		a.MiscErrs = append(a.MiscErrs, ErrDbInsertUser.Error())
		return nil
	}

	return &User{
                ID: newUser.ID,
                Email: newUser.Email,
                Username: newUser.Email,
                GlobalChatColor: newUser.GlobalChatColor,
        }
}

type SignIn struct {
	UserOrEmail    string
	UserOrEmailErr string
	Password       string
	PassErr        string
	MiscErrs       []string
}

func (a *SignIn) RenderErrs() []string {
	errs := []string{a.UserOrEmailErr, a.PassErr}
	errs = append(errs, a.MiscErrs...)
	rendered := []string{}
	for _, v := range errs {
		if v != "" {
			rendered = append(rendered, v)
		}
	}
	return rendered
}
func (a *SignIn) FlushPassword() {
	a.Password = ""
}

func (a *SignIn) SignIn() *db.InsertUserRow {
	if a.UserOrEmail == "" || a.Password == "" {
		a.MiscErrs = append(a.MiscErrs, ErrBadLogin.Error())
		a.FlushPassword()
		return nil
	}

	var user db.User
	ctx := context.Background()

	_, err := mail.ParseAddress(a.UserOrEmail)
	if err == nil {
		user, err = db.Db().SelectUserByEmailWithPassword(ctx, a.UserOrEmail)
		if err != nil {
			a.MiscErrs = append(a.MiscErrs, ErrBadLogin.Error())
			a.FlushPassword()
			return nil
		}
	} else {
		user, err = db.Db().SelectUserByUsernameWithPassword(ctx, a.UserOrEmail)
		if err != nil {
			a.MiscErrs = append(a.MiscErrs, ErrBadLogin.Error())
			a.FlushPassword()
			return nil
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(a.Password+user.PasswordSalt))
	if err != nil {
		a.MiscErrs = append(a.MiscErrs, ErrBadLogin.Error())
		a.FlushPassword()
		return nil
	}

	user.EncryptedPassword = ""

	return &db.InsertUserRow{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
}
