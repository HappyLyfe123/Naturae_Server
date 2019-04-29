package users

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"bytes"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"log"
)

func ChangePassword(userEmail, currPassword, newPassword string) *pb.ChangePasswordReply {
	userInfoDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	accountInfoCollection := userInfoDB.Collection(helpers.GetAccountInfoCollection())
	accountInfo, err := getLoginInfo(userEmail)
	//Database communication error
	if err != nil {
		return &pb.ChangePasswordReply{Status: &pb.Status{Code: helpers.GetNotFoundStatusCode(), Message: "No account has been found"}}
		//User had not authenticated the account yet
	}
	//Hash the user password
	checkHashPassword := helpers.GenerateHash(helpers.ConvertStringToByte(currPassword),
		helpers.ConvertStringToByte(accountInfo.Salt))

	//Compare the hash stored in the database and the curr hash password
	//If result is 1 then the password doesn't match
	if bytes.Compare(helpers.ConvertStringToByte(accountInfo.Password), checkHashPassword) == 1 {
		return &pb.ChangePasswordReply{Status: &pb.Status{Code: helpers.GetInvalidPasswordCode(), Message: "Password doesn't match"}}
	}
	//Generate random bytes of data to be use as salt for the password
	salt := helpers.GenerateRandomBytes(helpers.GetSaltLength())
	//Generate a hash for the user password
	hashPassword := helpers.GenerateHash(helpers.ConvertStringToByte(newPassword), salt)
	filter := bson.D{{"email", userEmail}}
	update := bson.D{{"$set", bson.D{{"password", hashPassword}, {"salt", salt}}}}
	_, err = accountInfoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error while saving user hash password and salt from forget password: %v", err)
		return &pb.ChangePasswordReply{Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(),
			Message: "server error"}}
	}
	log.Printf("%s change the account password", userEmail)
	//Get the user access and refresh token id
	return &pb.ChangePasswordReply{Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "Password had been change"}}
}
