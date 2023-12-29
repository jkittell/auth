package main

import (
	"bytes"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jkittell/auth/model"
	"github.com/jkittell/data/api/client"
	"testing"
	"time"
)

func newUser() *model.User {
	return &model.User{
		ID:        0,
		Email:     gofakeit.Email(),
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Password:  gofakeit.Password(true, true, true, true, true, 32),
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestUsersHandler(t *testing.T) {
	user := newUser()

	data, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
	}
	reader := bytes.NewReader(data)

	res, err := client.Post("http://localhost/users", nil, reader)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(res))

	// get
	res, err = client.Get("http://localhost/users", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(res))

}

func TestUsersHandler_ListUsers(t *testing.T) {

}

func TestUsersHandler_Authenticate(t *testing.T) {

}

func TestUsersHandler_DeleteUser(t *testing.T) {

}
