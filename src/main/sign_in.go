package main

import (
	"crypto/sha512"
	"fmt"

	"../helper"
)

func signIn() {
	hasher := sha512.New()
	password := []byte("Hello")
	hasher.Write(password)
	helper.generateSalt()
	fmt.Println()
}

func hashPassword() {

}
