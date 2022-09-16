package models

import (
	"context"
	"github.com/kamva/mgm/v3"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string `json:"email" bson:"email"`
	FirstName        string `json:"first_name,omitempty" bson:"first_name"`
	LastName         string `json:"last_name,omitempty" bson:"last_name"`
	Password         string `json:"password,omitempty" bson:"password"`
	Active           int    `json:"active" bson:"active"`
}

func (u *User) Creating(ctx context.Context) error {
	log.Println("at creating", u.Password)

	hashedByte, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}
	u.Password = string(hashedByte)

	u.Active = 1

	return nil
}

func (u *User) Created(ctx context.Context) error {
	u.Password = ""

	return nil
}

func (u *User) PasswordMatches(hashedPassword, password string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		log.Println("incorrect password")
		return err
	}

	return nil
}
