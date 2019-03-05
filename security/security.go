package security

import (
	"Naturae_Server/helpers"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// CheckAppKey : Check if the app key is valid
// @appKey : provided app key to be check
func CheckAppKey(appKey string) bool {
	//Compare the provided app key with current database app key
	if strings.Compare(appKey, "") == 0 {
		return true
	} else {
		return false
	}
}

//CheckAccessToken : Check is access token is valid
func CheckAccessToken(username, tokenID string) {
	fmt.Println(helpers.GetDeniedStatusCode())
}

//CheckRefreshToken : Check if the token is valid
func CheckRefreshToken(username, tokenID string) {

}

//GenerateTokenID : Generate a token id for the user
func GenerateTokenID() string {
	return helpers.ConvertByteToStringBase64(GenerateHash(GenerateRandomBytes(helpers.GetTokenLength()), nil))
}

//GenerateAuthenCode : generate an authentication to verify the user
func GenerateAuthenCode() (string, time.Time) {
	//The time that the authentication code is created
	expiredTime := time.Now().Add(time.Minute * 30)
	authenCode := GenerateRandomNumber(helpers.GetAuthCodeMinNum(), helpers.GetAuthCodeMaxNum()).String()

	return authenCode, expiredTime
}

//GenerateHash : Hash password with salt using sha512
func GenerateHash(originalData, salt []byte) []byte {
	//Add the salt to the end of the user password
	//After adding salt to user password convert the password with
	//salt into byte array
	//Hash the byte array password using sha515 then convert the hash to
	//base64 encoding and return it
	hashPassword := sha512.New().Sum(append(originalData, salt...))
	return hashPassword
}

//GenerateRandomBytes : random data
func GenerateRandomBytes(len int) []byte {
	//Generate an array of the specific byte in length
	newData := make([]byte, len)
	//If there an error it going to loop until there are no error
	for {
		//Generate a pseudo random number
		_, err := rand.Read(newData)
		if err == nil {
			return newData
		} else {
			log.Fatalf("Generate random bytes error: %v", err)
		}
	}

}

//GenerateRandomNumber : generate a random between 100000 and 999999
func GenerateRandomNumber(minNum, maxNum int64) *big.Int {
	//It going to loop until there are no error
	for {
		//Generate a number between 0 and the given number
		randomNum, err := rand.Int(rand.Reader, big.NewInt(maxNum))
		if err != nil {
			log.Fatalf("Random number generator error: %v", err)
		}
		//Add minNum to the generated number
		randomNum.Add(randomNum, big.NewInt(minNum))

		return randomNum

	}

}

func IsTokenValid(currDatabase *mongo.Database, collectionName, tokenID string) {
	switch collectionName {
	case helpers.GetAccessTokenCollection():

	}
}
