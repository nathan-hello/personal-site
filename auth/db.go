package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/nathan-hello/personal-site/db"
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

func dbValidateJwt(t string, user string) (*db.SelectUserByIdRow, error) {
	ctx := context.Background()

	token, err := db.Db().SelectTokenFromJwtString(ctx, t)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrJwtNotInDb
		}
		return nil, ErrDbSelectJwt
	}
	if !token.Valid {
		return nil, ErrJwtInvalidInDb
	}

	userObj, err := db.Db().SelectUserById(context.Background(), user)

	if err != nil {
		return nil, err
	}

	return &userObj, nil
}

func dbInvalidateJwtFamily(t string) error {

	ctx := context.Background()

	token, err := db.Db().SelectTokenFromJwtString(ctx, t)
	if err != nil {
		return ErrDbSelectJwt
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
