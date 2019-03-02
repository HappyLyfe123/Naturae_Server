package security

import (
	"Naturae_Server/helpers"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"math/big"
	"strings"
	"time"
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

//GenerateToken : Generate an access token and refresh token for user
//userEmail : the user that the token will be created for
//lifeSpan : the amount of time that the token going to last
func GenerateToken(userEmail string, tokenType int8, tokenIDChan chan string, startTimeChan chan time.Time) {
	//Initialize tokeID to an empty string
	tokenID := ""
	//Generate a start time for the token
	startTime := time.Now()
	//Generate an end time for the token
	randomData := GenerateRandomBytes(300)
	//Generate a token id
	tempTokenID := GenerateHash(randomData, helpers.ConvertStringToByte(userEmail))

	if tokenType == 1 {
		tokenID = helpers.ConvertByteToString(helpers.RandomizeData(tempTokenID))
	} else {
		tokenID = helpers.ConvertByteToString(tempTokenID)
	}

	tokenIDChan <- tokenID
	startTimeChan <- startTime
}

//GenerateAuthenCode : generate an authentication to verify the user
func GenerateAuthenCode() (string, time.Time) {
	//The time that the authentication code is created
	startTime := time.Now()
	authenCode := GenerateRandomNumber(helpers.GetAuthCodeMinNum(), helpers.GetAuthCodeMaxNum()).String()

	return authenCode, startTime
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
		//Generate a pseudorandom number
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

func IsTokenValid(currDatabase *mongo.Database, tokenID, collectionName string){

	//Check if it's a access token
	if strings.Compare(helpers.GetAccessTokenCollection(), collectionName) == 1{

	}else{

	}

}
