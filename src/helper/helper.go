package helper

import (
	"crypto/rand"
)

//GenerateRandomBytes : random data
func GenerateRandomBytes(len int) ([]byte, error) {
	newData := make([]byte, len)
	_, err := rand.Read(newData)
	if err != nil {
		return nil, err
	}

	return newData, nil
}

//GetUsername : Get the user useranme
func GetUsername(email string) {

}

//GetEmail : Get the user email address
func GetEmail(username string) {

}
