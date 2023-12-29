package model

import (
	"context"
	"errors"
	"github.com/jkittell/data/database"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

var users *database.PosgresDB[*User]

// User is the structure which holds one user from the database.
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Primary() (string, any) {
	return "id", u.ID
}

func (u *User) Scan(fields []string, scan database.ScanFunc) error {
	return database.Scan(map[string]any{
		"id":         &u.ID,
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

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	if err = users.Update(ctx, u); err != nil {
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
