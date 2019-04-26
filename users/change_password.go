package users

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
)

func ChangePassword(userEmail, password string) *pb.ChangePasswordReply {
	userInfoDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	userInfo, err := getLoginInfo(userInfoDB, userEmail)
	//Database communication error
	if err != nil {
		return &pb.ChangePasswordReply{Status: &pb.Status{Code: helpers.GetNotFoundStatusCode(), Message: "No account has been found"}}
		//User had not authenticated the account yet
	}
	return nil
}
