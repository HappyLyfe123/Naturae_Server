package main

import (
	"Naturae_Server/helpers"
	. "Naturae_Server/naturaeproto"
	"Naturae_Server/post"
	"Naturae_Server/users"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct{}

func main() {
	//Close the connection to the database when the server is turn off
	defer cleanUpServer()
}

//Initialize all of the variable to be uses
func init() {
	//Initialize global variable in the helper package
	helpers.ConnectToGmailAccount()
	helpers.ConnectToDBAccount()
	helpers.ConnectToS3()
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
	log.Println("Hello world")
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
		log.Println(request.GetEmail(), "fail to login")
		result = &LoginReply{AccessToken: "", RefreshToken: "", FirstName: "", LastName: "", Email: "", Status: &Status{
			Code: helpers.GetInvalidAppKey(), Message: "Invalid app key"}}
	}
	log.Printf("%s login", request.GetEmail())
	return result, nil
}

//Account authentication
func (s *server) AccountAuthentication(ctx context.Context, request *AccountAuthenRequest) (*AccountAuthenReply, error) {
	var result *AccountAuthenReply
	//Check if app key is valid
	if helpers.CheckAppKey(request.GetAppKey()) {
		result = users.AuthenticateAccount(request)
	} else {
		log.Println(request.GetEmail(), "failed to authenticated account")
		result = &AccountAuthenReply{Status: &Status{
			Code: helpers.GetInvalidAppKey(), Message: "Invalid app key"}}
	}
	log.Println(request.GetEmail(), "has authenticated the account")
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
				result = post.SavePost(request, accessToken.Email)
				log.Println("Post create by:", accessToken.Email)
			}
		}
	}

	return result, nil
}

//GetPosts : get all of the post within the give location
func (s *server) GetPosts(ctx context.Context, request *GetPostRequest) (*GetPostReply, error) {
	var result *GetPostReply
	if helpers.CheckAppKey(request.AppKey) {
		result = post.GetPost(request.GetLat(), request.GetLng())
	} else {
		result = &GetPostReply{Status: &Status{Code: helpers.GetInvalidAppKey(), Message: "invalid app key"},
			Reply: nil}
	}
	return result, nil
}

//SearchPost : Search for a specific post
func (s *server) SearchPost(context.Context, *SearchPostRequest) (*SearchPostReply, error) {
	panic("implement me")
}

//GetPostPreview : get all of the post that is within the radius
func (s *server) GetPostPreview(ctx context.Context, request *GetPostPreviewRequest) (*GetPostPreviewReply, error) {
	var result *GetPostPreviewReply
	if helpers.CheckAppKey(request.AppKey) {
		log.Println("Getting post preview")
		result = post.GetPostPreview(float64(request.GetRadius()), helpers.ConvertDegreeToRadian(float64(request.GetLat())),
			helpers.ConvertDegreeToRadian(float64(request.GetLng())))
	} else {
		result = &GetPostPreviewReply{Status: &Status{Code: helpers.GetInvalidAppKey(), Message: "invalid app key"},
			Reply: nil}
	}

	return result, nil
}

//ForgetPassword : Start the process of resetting user password by checking if the user email valid
// then it will send an email with a verification code to the user
func (s *server) ForgetPassword(ctx context.Context, request *ForgetPasswordRequest) (*ForgetPasswordReply, error) {
	var result *ForgetPasswordReply
	if helpers.CheckAppKey(request.GetAppKey()) {
		result = users.ForgetPasswordCreateResetCode(request)
	}
	return result, nil
}

//ForgetPasswordVerifyCode : User to verifying forget password verification code
func (s *server) ForgetPasswordVerifyCode(ctx context.Context, request *ForgetPasswordVerifyCodeRequest) (*ForgetPasswordVerifyCodeReply, error) {
	var result *ForgetPasswordVerifyCodeReply
	if helpers.CheckAppKey(request.GetAppKey()) {
		result = users.ForgetPasswordVerifyCode(request)
	}
	return result, nil
}

//ForgetPasswordResetPassword : reset the password for the user
func (s *server) ForgetPasswordResetPassword(ctx context.Context, request *ForgetPasswordNewPasswordRequest) (*ForgetPasswordNewPasswordReply, error) {
	var result *ForgetPasswordNewPasswordReply
	if helpers.CheckAppKey(request.GetAppKey()) {
		result = users.ForgetPasswordNewPassword(request)
	}
	return result, nil
}

//ChangePassword : change the user password
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

func (s *server) GetProfileImage(ctx context.Context, request *GetProfileImageRequest) (*GetProfileImageReply, error) {
	var result *GetProfileImageReply
	//Check if the app key if valid
	if helpers.CheckAppKey(request.GetAppKey()) {
		connectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())
		accessToken, err := helpers.GetAccessToken(connectedDB, request.GetAccessToken())
		if err != nil {
			result = &GetProfileImageReply{EncodedImage: "", Status: &Status{Code: helpers.GetInvalidTokenCode(), Message: "token is not valid"}}
		} else {
			//Check if the token is expired
			if helpers.IsTokenExpired(accessToken.ExpiredTime) {
				result = &GetProfileImageReply{EncodedImage: "", Status: &Status{Code: helpers.GetExpiredAccessTokenCode(), Message: "token had expired"}}
			} else {
				//Get the user profile image
				result = users.GetProfileImage(accessToken.Email)
			}
		}
	} else {
		result = &GetProfileImageReply{EncodedImage: "", Status: &Status{Code: helpers.GetInvalidAppKey(), Message: "app key is invalid"}}
	}
	return result, nil
}

func (s *server) SetProfileImage(ctx context.Context, request *SetProfileImageRequest) (*SetProfileImageReply, error) {
	var result *SetProfileImageReply
	//Check if the user app key is valid
	if helpers.CheckAppKey(request.GetAppKey()) {
		connectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())
		accessToken, err := helpers.GetAccessToken(connectedDB, request.GetAccessToken())
		if err != nil {
			result = &SetProfileImageReply{Status: &Status{Code: helpers.GetInvalidTokenCode(), Message: "token is not valid"}}
		} else {
			//Check the if the token is still valid or if it expired
			if helpers.IsTokenExpired(accessToken.ExpiredTime) {
				result = &SetProfileImageReply{Status: &Status{Code: helpers.GetExpiredAccessTokenCode(), Message: "token had expired"}}
			} else {
				//Save the image to AWS S3
				imageURL, saveImageResult := users.SaveProfileImage(request)
				if saveImageResult {
					//User the user information of the new profile image
					result = users.UpdateProfileImage(accessToken.Email, imageURL)
				} else {
					result = &SetProfileImageReply{Status: &Status{Code: helpers.GetExpiredAccessTokenCode(), Message: "token had expired"}}
				}
			}
		}
	} else {
		result = &SetProfileImageReply{Status: &Status{Code: helpers.GetInvalidAppKey(), Message: "app key is invalid"}}
	}
	return result, nil
}

//User/Friend Search
func (s *server) SearchUsers(ctx context.Context, request *UserSearchRequest) (*UserListReply, error) {
	return users.SearchUsers(request), nil
}

//Friend Adding
func (s *server) AddFriend(ctx context.Context, request *FriendRequest) (*FriendReply, error) {
	return users.AddFriend(request), nil
}

//Friend Removal
func (s *server) RemoveFriend(ctx context.Context, request *FriendRequest) (*FriendReply, error) {
	return users.RemoveFriend(request), nil
}

//Room Retrieval
func (s *server) GetRoomName(ctx context.Context, request *RoomRequest) (*RoomReply, error) {
	return users.GetRoomName(request), nil
}
