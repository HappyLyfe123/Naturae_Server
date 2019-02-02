package main

import (
	"crypto/sha512"
	"strings"
)

type User struct {
	Username string
	Email    string
	Salt     string
	Password string
}

//CreateAccount : User want to create an account
/*
 * emai: user email address
 * username: user username
 * password: user password
 *
 */
func CreateAccount(email, username, password string) {
	//Check if email, username, or password is empty
	if len(strings.TrimSpace(email)) != 0 || len(strings.TrimSpace(username)) != 0 || len(strings.TrimSpace(password)) != 0 {

	} else {

	}
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
