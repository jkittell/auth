package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	r := struct {
		Id int `json:"id"`
	}{}
	err = json.Unmarshal(res, &r)

	// get
	res, err = client.Get(fmt.Sprintf("http://localhost/users/%d", r.Id), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(res))

	// update
	err = json.Unmarshal(res, user)
	if err != nil {
		t.Fatal(err)
	}

	user.FirstName = "updated"

	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	res, err = client.Put(fmt.Sprintf("http://localhost/users/%d", user.ID), nil, bytes.NewReader(userJSON))

	all, err := client.Get("http://localhost/users", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(all))
}

func TestUsersHandler_Authenticate(t *testing.T) {
	user := newUser()
	password := user.Password

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

	r := struct {
		Id int `json:"id"`
	}{}
	err = json.Unmarshal(res, &r)

	// get
	res, err = client.Get(fmt.Sprintf("http://localhost/users/%d", r.Id), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(res, user)
	if err != nil {
		t.Fatal(err)
	}

	auth := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    user.Email,
		Password: password,
	}

	data, err = json.Marshal(auth)
	if err != nil {
		t.Fatal(err)
	}

	res, err = client.Post("http://localhost/authenticate", nil, bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(res))

}
