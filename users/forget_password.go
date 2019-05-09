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

//The structure for reset password code
type ResetCode struct {
	Email   string
	Code    string
	Expired time.Time
}

//ForgetPasswordCreateResetCode: create an reset password code for the user
func ForgetPasswordCreateResetCode(request *pb.ForgetPasswordRequest) *pb.ForgetPasswordReply {
	userDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	forgetPasswordCollection := userDB.Collection("Forget_Password")
	//Generate authentication code and expired time
	resetCode, expiredTime := helpers.GenerateAuthenCode()
	newResetCode := ResetCode{Email: request.Email, Code: resetCode, Expired: expiredTime}
	//Save the reset password code to the server
	_, err := forgetPasswordCollection.InsertOne(context.Background(), newResetCode)
	if err != nil {
		log.Printf("Saving reset password code error: %v", err)
		return &pb.ForgetPasswordReply{Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "server error"}}
	}
	//Send the reset password code to the user email
	sendResetPasswordCode(request.GetEmail(), resetCode)
	return &pb.ForgetPasswordReply{Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "code created"}}
}

//VerifyForgetPasswordCode : verify if the code that the user entered to reset their password
func ForgetPasswordVerifyCode(request *pb.ForgetPasswordVerifyCodeRequest) *pb.ForgetPasswordVerifyCodeReply {
	userDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	forgetPasswordCollection := userDB.Collection("Forget_Password")
	var result ResetCode
	//Get the verification code from the server
	err := forgetPasswordCollection.FindOne(context.Background(), bson.D{{"email", request.Email}}).Decode(&result)
	if err != nil {
		log.Printf("Getting forget password reset code error: %v", err)
		return &pb.ForgetPasswordVerifyCodeReply{Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "server error"}}
	}
	//Compare the verification code the user provided with the one that stored in the database
	if strings.Compare(result.Code, request.GetVerificationCode()) == 0 {
		_, err := forgetPasswordCollection.DeleteOne(context.Background(), bson.D{{"email", request.Email}})
		if err != nil {
			return &pb.ForgetPasswordVerifyCodeReply{Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "server error"}}
		}
		return &pb.ForgetPasswordVerifyCodeReply{Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "code match"}}
	}
	return &pb.ForgetPasswordVerifyCodeReply{Status: &pb.Status{Code: helpers.GetInvalidCode(), Message: "code is invalid"}}

}

//ForgetPasswordNewPassword : update the user change the user password and salt
func ForgetPasswordNewPassword(request *pb.ForgetPasswordNewPasswordRequest) *pb.ForgetPasswordNewPasswordReply {
	//Check if the password is in a valid format
	if helpers.IsPasswordValid(request.Password) {
		userDB := helpers.ConnectToDB(helpers.GetUserDatabase())
		accountInfo := userDB.Collection(helpers.GetAccountInfoCollection())
		//Generate random bytes of data to be use as salt for the password
		salt := helpers.GenerateRandomBytes(helpers.GetSaltLength())
		//Generate a hash for the user password
		hashPassword := helpers.GenerateHash(helpers.ConvertStringToByte(request.GetPassword()), salt)
		filter := bson.D{{"email", request.Email}}
		update := bson.D{{"$set", bson.D{{"password", helpers.ConvertByteToStringBase64(hashPassword)},
			{"salt", helpers.ConvertByteToStringBase64(salt)}}}}
		_, err := accountInfo.UpdateOne(context.Background(), filter, update)
		if err != nil {
			log.Printf("Error while saving user hash password and salt from forget password: %v", err)
			return &pb.ForgetPasswordNewPasswordReply{Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(),
				Message: "server error"}}
		}
		//Send the user an email that their account password had been reset
		sendResetPasswordSuccessMessage(request.Email)
		return &pb.ForgetPasswordNewPasswordReply{Status: &pb.Status{Code: helpers.GetOkStatusCode(),
			Message: "password is had been reset"}}
	}
	return &pb.ForgetPasswordNewPasswordReply{Status: &pb.Status{Code: helpers.GetInvalidPasswordCode(),
		Message: "password is invalid"}}
}

func sendResetPasswordSuccessMessage(userEmail string) {
	body := fmt.Sprint("Your have successfully reset your password. If your didn't request for this change, please" +
		"contact us at naturae.outdoor@gmail.com immediately otherwise ignore this message.")
	newMail := helpers.Email{Receiver: userEmail, Subject: "Password reset", Body: body}
	err := helpers.SendEmail(&newMail)
	//If there no error
	if err != nil {
		log.Fatalln("Send email error: ", err)
	}
	log.Println("Email", userEmail, "successfully reset password")
}

//SendAuthenticationEmail : Send a confirmation email to the user to make sure it's the user email address
func sendResetPasswordCode(userEmail, resetCode string) {
	//The system will be send a 6 digits number to the user provided email
	//This six digits number will be use to ensure that it's the user email
	body := fmt.Sprintf("Your Naturae password reset code is: %s\nThis code will expired after 30 minutes.", resetCode)
	//Send the email to the user
	newMail := helpers.Email{Receiver: userEmail, Subject: "Reset Password Code", Body: body}
	err := helpers.SendEmail(&newMail)
	//If there no error
	if err != nil {
		log.Fatalln("Send email error: ", err)
	}
	log.Println("Email", userEmail, "reset password code")

}
