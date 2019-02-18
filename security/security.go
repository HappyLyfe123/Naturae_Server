package security

import (
	"Naturae_Server/helpers"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type Token struct {
	TokenID    string
	Email      string
	StartTime  time.Time
	ExpireTime time.Time
}

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

//GenerateToken : Generate an access token and refresh token for user
//userEmail : the user that the token will be created for
//lifeSpan : the amount of time that the token going to last
func GenerateToken(userEmail *string, lifeSpan []int) (newToken *Token) {
	//Generate a start time for the token
	startTime := time.Now()
	//Generate an end time for the token
	//Index 0 = year Index 1 = month Index 2 = day
	expireTime := startTime.AddDate(lifeSpan[0], lifeSpan[1], lifeSpan[2])

	newToken = &Token{*userEmail, "", startTime, expireTime}

	return
}

//GenerateHash : Hash password with salt using sha512
func GenerateHash(userPassword, salt *[]byte) *[]byte {
	//Add the salt to the end of the user password
	//After adding salt to user password convert the password with
	//salt into byte array
	//Hash the byte array password using sha515 then convert the hash to
	//base64 encoding and return it
	hashPassword := sha512.New().Sum(append(*userPassword, *salt...))
	return &hashPassword
}

//GenerateRandomBytes : random data
func GenerateRandomBytes(len int) (newData *[]byte, err error) {
	//Generate an array of the specific byte in length
	*newData = make([]byte, len)
	//Generate a pseudorandom number
	_, err = rand.Read(*newData)
	if err != nil {
		return
	}
	//Convert the new byte of data into base64 string
	return
}

//GenerateRandomNumber : generate a random between 100000 and 999999
func GenerateRandomNumber(minNum, maxNum int64) (randomNum *big.Int) {
	//Generate a number between 0 and the given number
	randomNum, err := rand.Int(rand.Reader, big.NewInt(maxNum))
	if err != nil {
		fmt.Println(err)
	}
	//Add minNum to the generated number
	randomNum.Add(randomNum, big.NewInt(minNum))

	return
}
