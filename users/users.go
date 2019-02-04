package users

import (
	sha512 "crypto/sha512"
	"strings"
	"time"

	"../helpers"
	"../security"
)

type user struct {
	Username string
	Email    string
	Salt     string
	Password string
}

type token struct {
	Username  string
	TokenID   string
	StartTime time.Time
	EndTime   time.Time
}

var saltLength int16 = 200


//CreateAccount : User want to create an account
/*
 * emai: user email address
 * username: user username
 * password: user password
 *
 */
func CreateAccount(email, username, password string) error {
	//Check if email, username, or password is  not empty
	if len(strings.TrimSpace(email)) != 0 || len(strings.TrimSpace(username)) != 0 || len(strings.TrimSpace(password)) != 0 {
		salt, err := helpers.GenerateRandomBytes(saltLength)
		if err != nil {
			return err
		}
		//Hash the user password
		hashPassword := security.HashPassword(password, salt)
		//Save the user into a struct
		newUser := user{username, email, salt, hashPassword}

	} else {

	}

	return nil
}

// SignIn : User sign in
/*
 * email: email that want to login
 *
 */
func SignIn(email string) {
	hasher := sha512.New()
	password := []byte("Hello")
	hasher.Write(password)
}

// ForgotPassword : User forget password and want to reset it
/*
 *
 * email: user emal
 *
 */
func ForgotPassword(email string) {

}

//ChangePassword : User want to change password
/*
 * username: User email address
 * oldPassword: User old password
 * newPassword: The password the user want to change to
 *
 */
func ChangePassword(email, oldPassword, newPassword string) {

}
