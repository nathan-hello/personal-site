package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/nathan-hello/personal-site/src/db"
	"github.com/nathan-hello/personal-site/src/utils"
)

func TestDbSelectUsers(t *testing.T) {
	ctx := context.Background()
	err := utils.InitEnv()

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(utils.Env())

	err = db.InitDb()
	if err != nil {
		t.Fatal(err)
	}
	f := db.Conn()

	users, err := f.SelectAllUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Print(users)
}

// func TestDbInsertUser(t *testing.T) {
// 	ctx := context.Background()
// 	err  := db.InitDb()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	f := db.Conn()
//
// 	user, err := f.InsertUser(ctx, db.InsertUserParams{
// 		Username:          "black-bear-test-22121231321",
// 		EncryptedPassword: "honey",
// 	})
//
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	defer func() {
// 		err = f.DeleteUser(ctx, user.ID)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()
//
//
// 	t.Logf("user: %#v\n", user)
// 	fullUser, err := f.SelectUserByUsername(ctx, user.Username)
//
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	t.Logf("fullUser user: %#v\n", fullUser)
//
//
//
//
// }
