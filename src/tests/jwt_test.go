package test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nathan-hello/personal-site/src/auth"
	"github.com/nathan-hello/personal-site/src/db"
	"github.com/nathan-hello/personal-site/src/utils"
)

func TestNewPairAndParse(t *testing.T) {
	access, refresh, err := auth.NewTokenPair(
		&auth.JwtParams{
			UserId:   uuid.New().String(),
			Username: "black-bear",
			Family:   uuid.New(),
		})

	if err != nil {
		t.Error(err)
	}

	// t.Logf("TestNewPairAndParse/access: %v\n", access)
	// t.Logf("TestNewPairAndParse/refresh: %v\n", refresh)

	_, err = auth.ParseToken(access)
	if err != nil {
		t.Error(err)
	}

	_, err = auth.ParseToken(refresh)
	if err != nil {
		t.Error(err)
	}

	// t.Logf("TestNewPairAndParse/ap: %v\n", ap)
	// t.Logf("TestNewPairAndParse/rp: %v\n", rp)

}
func TestJwtExpiry(t *testing.T) {
	c := utils.Env()
	c.ACCESS_EXPIRY_TIME = time.Second * 1
	c.REFRESH_EXPIRY_TIME = time.Second * 2

	access, refresh, err := auth.NewTokenPair(
		&auth.JwtParams{
			UserId:   uuid.New().String(),
			Username: "black-bear",
			Family:   uuid.New(),
		})

	if err != nil {
		t.Error(err)
	}

	// t.Logf("access: %#v\n", access)
	// t.Logf("refresh: %#v\n", refresh)

	_, err = auth.ParseToken(access)

	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Millisecond * 1100)

	_, err = auth.ParseToken(access)

	if err == nil {
		t.Error("access is still valid even after waiting past expiration time")
	}

	// t.Logf("token successfully invalidated: %#v\n", err)

	_, err = auth.ParseToken(refresh)

	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 1)

	_, err = auth.ParseToken(refresh)

	if err == nil {
		t.Error("refresh is still valid even after waiting past expiration time")
	}

	// t.Logf("refresh successfully invalidated: %#v\n", err)
}

func TestDbJwt(t *testing.T) {
	ctx := context.Background()
	db.InitDb()

	f := db.Conn()

	fullUser, err := f.InsertUser(ctx, db.InsertUserParams{
		Username:          "black-bear-test-1",
		EncryptedPassword: "honey",
	})

	defer func() {
		err = f.DeleteUser(ctx, fullUser.ID)

		if err != nil {
			t.Error(err)
		}
		// fmt.Printf("deleted user: %#v\n", fullUser.ID.String())
	}()

	if err != nil {
		t.Error(err)
	}

	access, refresh, err := auth.NewTokenPair(
		&auth.JwtParams{
			Username: fullUser.Username,
			UserId:   fullUser.ID,
			Family:   uuid.New(),
		})
	if err != nil {
		t.Error(err)
	}

	_, err = auth.ParseToken(access)
	if err != nil {
		t.Error(err)
	}
	_, err = auth.ParseToken(refresh)
	if err != nil {
		t.Error(err)
	}

	err = auth.InsertNewToken(access, "access_token")
	defer func() {
		err := f.DeleteTokensByUserId(ctx, fullUser.ID)
		if err != nil {
			t.Error(err)
		}
	}()
	if err != nil {
		t.Error(err)
	}
	err = auth.InsertNewToken(refresh, "refresh_token")
	if err != nil {
		t.Error(err)
	}

	tokens, err := f.SelectUsersTokens(ctx, fullUser.ID)
	if err != nil {
		t.Error(err)
	}

	if len(tokens) == 0 {
		t.Error("Token length 0\n")
	}

	err = f.UpdateUserTokensToInvalid(ctx, fullUser.ID)
	if err != nil {
		t.Error(err)
	}

	_, err = f.SelectUsersTokens(ctx, fullUser.ID)
	if err != nil {
		t.Error(err)
	}

}
