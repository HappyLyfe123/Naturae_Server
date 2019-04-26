package main

import (
	"Naturae_Server/helpers"
	. "Naturae_Server/naturaeproto"
	"Naturae_Server/post"
	"Naturae_Server/users"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct{}

func main() {
	//Close the connection to the database when the server is turn off
	defer cleanUpServer()
	//post.GetPost()
	//Lat := 14.55
	//Lon := 25.00
	//fmt.Println(math.Acos(math.Sin(1.3963)*math.Sin(Lat)+math.Cos(1.3963)*math.Cos(Lat)*math.Cos(Lon-(-0.6981))) * 6371)
}

//Initialize all of the variable to be uses
func init() {
	//Initialize global variable in the helper package
	helpers.ConnectToGmailAccount()
	helpers.ConnectToDBAccount()
	//Create listener for server
	createServer()

}

//Close all of the connection to everything that the server is connected to
func cleanUpServer() {
	err := helpers.CloseConnectionToDatabaseAccount()
	if err != nil {
		log.Println("Closed connection to DB error: ", err)
	}
}

//Initialize and start the server
func createServer() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("unable to listen on 8080 port: %v", err)
	}
	log.Println("listening on port 8080")
	srv := grpc.NewServer()
	RegisterServerRequestsServer(srv, &server{})
	reflection.Register(srv)
	err = srv.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (s *server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{
		Message: "Hello " + in.Name,
	}, nil
}

//Create user account
func (s *server) CreateAccount(ctx context.Context, request *CreateAccountRequest) (*CreateAccountReply, error) {
	var result *CreateAccountReply
	//Check if the app key is valid
	if helpers.CheckAppKey(request.GetAppKey()) {
		result = users.CreateAccount(request)
	} else {
		result = &CreateAccountReply{Status: &Status{
			Code: helpers.GetInvalidAppKey(), Message: "Invalid app key"}}
	}

	return result, nil
}

//Login user
func (s *server) Login(ctx context.Context, request *LoginRequest) (*LoginReply, error) {
	var result *LoginReply
	//Check if the app key is valid
	if helpers.CheckAppKey(request.GetAppKey()) {
		result = users.Login(request)
	} else {
		result = &LoginReply{AccessToken: "", RefreshToken: "", FirstName: "", LastName: "", Email: "", Status: &Status{
			Code: helpers.GetInvalidAppKey(), Message: "Invalid app key"}}
	}
	return result, nil
}

//Account authentication
func (s *server) AccountAuthentication(ctx context.Context, request *AccountAuthenRequest) (*AccountAuthenReply, error) {
	var result *AccountAuthenReply
	//Check if app key is valid
	if helpers.CheckAppKey(request.GetAppKey()) {
		result = users.AuthenticateAccount(request)
	} else {
		result = &AccountAuthenReply{Status: &Status{
			Code: helpers.GetInvalidAppKey(), Message: "Invalid app key"}}
	}
	return result, nil
}

//Generate a new access token
func (s *server) GetNewAccessToken(ctx context.Context, request *GetAccessTokenRequest) (*GetAccessTokenReply, error) {
	var result *GetAccessTokenReply
	//Check if app key is valid
	if helpers.CheckAppKey(request.GetAppKey()) {
		result = users.RefreshAccessToken(request)
	} else {
		result = &GetAccessTokenReply{AccessToken: "", Status: &Status{
			Code: helpers.GetInvalidAppKey(), Message: "Invalid app key"}}
	}

	return result, nil
}

//func (s *server) ChangePassword(ctx context.Context, request *ChangePasswordRequest) (*ChangePasswordReply, error) {
//	var result *ChangePasswordReply
//
//	if helpers.CheckAppKey(request.GetAppKey()) {
//		connectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())
//		accessToken, err := helpers.GetAccessToken(connectedDB, request.GetAccessToken())
//		//Check if there an error then the access token provided is not in the database
//		if err != nil {
//			result = &ChangePasswordReply{Status: &Status{Code: helpers.GetInvalidTokenCode(), Message: "token is not valid"}}
//		} else {
//			//Check if the access token is expired
//			if helpers.IsTokenExpired(accessToken.ExpiredTime) {
//				result = &ChangePasswordReply{Status: &Status{Code: helpers.GetExpiredAccessTokenCode(), Message: "token is " +
//					"had expired"}}
//			} else {
//
//			}
//
//		}
//	}
//
//	return result, nil
//}

func (s *server) CreatePost(ctx context.Context, request *CreatePostRequest) (*CreatePostReply, error) {
	var result *CreatePostReply
	if helpers.CheckAppKey(request.GetAppKey()) {
		connectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())
		accessToken, err := helpers.GetAccessToken(connectedDB, request.GetAccessToken())
		//Check if there an error then the access token provided is not in the database
		if err != nil {
			result = &CreatePostReply{Status: &Status{Code: helpers.GetInvalidTokenCode(), Message: "token is not valid"}}
		} else {
			//Check if the access token is expired
			if helpers.IsTokenExpired(accessToken.ExpiredTime) {
				result = &CreatePostReply{Status: &Status{Code: helpers.GetExpiredAccessTokenCode(), Message: "token had expired"}}
			} else {
				fmt.Println("Post create by:", accessToken.Email)
				result = post.SavePost(request, accessToken.Email)
			}
		}
	}

	return result, nil
}

func (s *server) GetPosts(ctx context.Context, request *GetPostRequest) (*GetPostReply, error) {
	if helpers.CheckAppKey(request.AppKey) {

	}
	panic("implement me")
}

func (s *server) ForgetPassword(ctx context.Context, request *ForgetPasswordRequest) (*ForgetPasswordReply, error) {
	var result *ForgetPasswordReply
	if helpers.CheckAppKey(request.GetAppKey()) {
		result = users.ForgetPasswordCreateResetCode(request)
	}
	return result, nil
}

func (s *server) ForgetPasswordVerifyCode(ctx context.Context, request *ForgetPasswordVerifyCodeRequest) (*ForgetPasswordVerifyCodeReply, error) {
	var result *ForgetPasswordVerifyCodeReply
	if helpers.CheckAppKey(request.GetAppKey()) {
		result = users.ForgetPasswordVerifyCode(request)
	}
	return result, nil
}

func (s *server) ForgetPasswordResetPassword(ctx context.Context, request *ForgetPasswordNewPasswordRequest) (*ForgetPasswordNewPasswordReply, error) {
	var result *ForgetPasswordNewPasswordReply
	if helpers.CheckAppKey(request.GetAppKey()) {
		result = users.ForgetPasswordNewPassword(request)
	}
	return result, nil
}

func (s *server) ChangePassword(ctx context.Context, request *ChangePasswordRequest) (*ChangePasswordReply, error) {
	var result *ChangePasswordReply
	if helpers.CheckAppKey(request.GetAppKey()) {
		connectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())
		accessToken, err := helpers.GetAccessToken(connectedDB, request.GetAccessToken())
		//Check if there an error then the access token provided is not in the database
		if err != nil {
			result = &ChangePasswordReply{Status: &Status{Code: helpers.GetInvalidTokenCode(), Message: "token is not valid"}}
		} else {
			//Check if the access token is expired
			if helpers.IsTokenExpired(accessToken.ExpiredTime) {
				result = &ChangePasswordReply{Status: &Status{Code: helpers.GetExpiredAccessTokenCode(), Message: "token had expired"}}
			} else {
				result = users.ChangePassword(accessToken.Email, request.GetCurrentPassword(), request.GetNewPassword())
			}
		}
	}

	return result, nil
}
