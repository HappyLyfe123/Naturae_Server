package users

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"strings"
	"time"
)

type ResetCode struct {
	Email   string
	Code    string
	Expired time.Time
}

func CreateForgetPassowrdResetCode(request *pb.ForgetPasswordRequest) {
	userDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	forgetPasswordCollection := userDB.Collection("Forget_Password")
	//Generate authentication code and expired time
	resetCode, expiredTime := helpers.GenerateAuthenCode()
	newResetCode := ResetCode{Email: request.Email, Code: resetCode, Expired: expiredTime}
	//Save the reset password code to the server
	_, err := forgetPasswordCollection.InsertOne(context.Background(), newResetCode)
	if err != nil {
		log.Printf("Saving reset password code error: %v", err)
	}
	//Send the reset password code to the user email
	sendResetPasswordCode(request.GetEmail(), resetCode)
}

func VerifyForgetPasswordCode(request *pb.ForgetPasswordAuthenRequest) {
	userDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	forgetPasswordCollection := userDB.Collection("Forget_Password")
	var result *ResetCode
	err := forgetPasswordCollection.FindOne(context.Background(), bson.D{{"email", ""}}).Decode(result)
	if err != nil {
		log.Printf("Getting forget password reset code error: %v", err)
	}
	if strings.Compare(result.Code, request.AuthenCode) == 0 {

	}

}

//SendAuthenticationEmail : Send a confirmation email to the user to make sure it's the user email address
func sendResetPasswordCode(userEmail, resetCode string) {
	//The system will be send a 6 digits number to the user provided email
	//This six digits number will be use to ensure that it's the user email
	body := fmt.Sprintf("Your Naturae password reset code is %s\nThis code will expired after 30 minutes.", resetCode)
	//Send the email to the user
	newMail := helpers.Email{Receiver: userEmail, Subject: "Reset Password Code", Body: body}
	err := helpers.SendEmail(&newMail)
	//If there no error
	if err != nil {
		log.Fatalln("Send email error: ", err)
	}
	log.Println("Email", userEmail, "reset password code")

}
