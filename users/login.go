package users

import (
	"Naturae_Server/helpers"
	"Naturae_Server/security"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"strings"
)

type loginInfo struct {
	Salt            string
	Password        string
	IsAuthenticated bool
}

//Login : Let the user login into their account
func Login(email, password string) {
	connectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	databaseResult, err := getLoginInfo(connectedDB, email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Getting login info in login error:", err)
			return
			//Check if the user already authenticated the account
		} else {
			log.Println("Getting login info format error: ")
			return
		}
	} else if databaseResult.IsAuthenticated == false {
		return
	}
	//Hash and salt the password the user provided
	hashPassword := security.GenerateHash(helpers.ConvertStringToByte(password),
		helpers.ConvertStringToByte(databaseResult.Salt))

	//Check if the password match
	if strings.Compare(databaseResult.Password, helpers.ConvertByteToStringBase64(hashPassword)) == 1 {
		fmt.Println("Password match")
	} else {
		fmt.Println("Password does not match")
	}
}
