package users

type loginInfo struct {
	Salt            string
	Password        string
	IsAuthenticated bool
}

type loginResponse struct {
	AccessToken  string
	RefreshToken string
}

////Login : Let the user login into their account
//func Login(email, password string) (*loginResponse, *pb.Status) {
//	userInfo := helpers.ConnectToDB(helpers.GetUserDatabase())
//	databaseResult, err := getLoginInfo(userInfo, email)
//	if err != nil{
//
//	}
//}
