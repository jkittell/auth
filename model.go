package main

import (
	"github.com/jkittell/data/database"
	"time"
)

// User is the structure which holds one user from the database.
type User struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Primary() (string, any) {
	return "id", u.Id
}

func (u *User) Scan(fields []string, scan database.ScanFunc) error {
	return database.Scan(map[string]any{
		"id":         &u.Id,
		"email":      &u.Email,
		"first_name": &u.FirstName,
		"last_name":  &u.LastName,
		"password":   &u.Password,
		"active":     &u.Active,
		"created_at": &u.CreatedAt,
		"updated_at": &u.UpdatedAt,
	}, fields, scan)
}

func (u *User) Params() map[string]any {
	return map[string]any{
		"email":      &u.Email,
		"first_name": &u.FirstName,
		"last_name":  &u.LastName,
		"password":   &u.Password,
		"active":     &u.Active,
		"created_at": &u.CreatedAt,
		"updated_at": &u.UpdatedAt,
	}
}
