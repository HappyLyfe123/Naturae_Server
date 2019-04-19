package helpers

const (
	invalidName            = 100
	invalidEmail           = 103
	invalidPassword        = 102
	invalidLoginCredential = 103
	invalidToken           = 104
	invalidAuthenCode      = 105
	accountNotVerify       = 106
	invalidAppKey          = 107
	emailExist             = 150
	expiredAuthenCode      = 151
	expiredAccessToken     = 152
	expiredRefreshToken    = 153

	ok                  = 200
	created             = 201
	accepted            = 202
	noError             = 203
	invalidInformation  = 204
	denied              = 205
	duplicateInfo       = 206
	badRequest          = 400
	unauthorized        = 402
	forbidden           = 403
	notFound            = 404
	notAcceptable       = 406
	requestTimeout      = 408
	internalServerError = 500
	badGateway          = 502
	serviceUnavailable  = 503

	passwordSaltLength = 200
	tokenLength        = 150
	authCodeMaxNum     = 900000
	authCodeMinNum     = 100000

	gmailSMTPServer                 = "smtp.gmail.com:587"
	gmailSMTPHost                   = "smtp.gmail.com"
	gmailEmailAddress               = "naturae.outdoor@gmail.com"
	appName                         = "Naturae"
	databaseName                    = "Naturae-Server"
	userDatabase                    = "users"
	postDatabase                    = "posts"
	accountInfoCollection           = "Account_Information"
	accessTokenCollection           = "Access_Token"
	storeImageCollection            = "Images"
	storePostDescription            = "Post_Description"
	refreshTokenCollection          = "Refresh_Token"
	accountAuthenticationCollection = "Account_Authentication"
	appKey                          = "XjJzwDO6lBCQ5LbWlL1PCXIO&6jH@M&Mz8&BLq"
)

//GetInvalidNameCode : Return invalid name error code
func GetInvalidNameCode() int32 {
	return invalidName
}

//GetInvalidEmailCode : Return invalid email error code
func GetInvalidEmailCode() int32 {
	return invalidEmail
}

//GetInvalidPasswordCode : Return invalid password error code
func GetInvalidPasswordCode() int32 {
	return invalidPassword
}

//GetInvalidLoginCredential : Return invalid login credential error code
func GetInvalidLoginCredentialCode() int32 {
	return invalidLoginCredential
}

//GetInvalidTokenCode : Return invalid token error code
func GetInvalidTokenCode() int32 {
	return invalidToken
}

//GetInvalidAuthenCode : Return invalid authentication code
func GetInvalidAuthenCode() int32 {
	return invalidAuthenCode
}

//GetInvalidAppKey : Return invalid app key work
func GetInvalidAppKey() int32 {
	return invalidAppKey
}

//GetAccountNotVerifyCode : Return account not verify code
func GetAccountNotVerifyCode() int32 {
	return accountNotVerify
}

//GetEmailExistCode : Return email exist error code
func GetEmailExistCode() int32 {
	return emailExist
}

//GetExpiredAccessTokenCode : Return expired access token code
func GetExpiredAccessTokenCode() int32 {
	return expiredAccessToken
}

//GetExpiredRefreshTokenCode : Return expired refresh token code
func GetExpiredRefreshTokenCode() int32 {
	return expiredRefreshToken
}

//GetOkStatusCode : Return ok status code
func GetOkStatusCode() int32 {
	return ok
}

//GetCreatedStatusCode : Return created status code
func GetCreatedStatusCode() int32 {
	return created
}

//GetAcceptedStatusCode : Return accepted status code
func GetAcceptedStatusCode() int32 {
	return accepted
}

//GetBadRequestStatusCode : Return bad request status code
func GetBadRequestStatusCode() int32 {
	return badRequest
}

//GetUnauthorizedStatusCode : Return unauthorized status code
func GetUnauthorizedStatusCode() int32 {
	return unauthorized
}

//GetForbiddenStatusCode : Return forbidden status code
func GetForbiddenStatusCode() int32 {
	return forbidden
}

//GetNotFoundStatusCode : Return not found status code
func GetNotFoundStatusCode() int32 {
	return notFound
}

//GetNotAcceptableStatusCode : Return not acceptable status code
func GetNotAcceptableStatusCode() int32 {
	return notAcceptable
}

//GetRequestTimeoutStatusCode : Return request timeout status code
func GetRequestTimeoutStatusCode() int32 {
	return requestTimeout
}

//GetInternalServerErrorStatusCode : Return internal server error status code
func GetInternalServerErrorStatusCode() int32 {
	return internalServerError
}

//GetBadGatewayStatusCode : Return bad gateway status code
func GetBadGatewayStatusCode() int32 {
	return badGateway
}

//GetServiceUnavailableStatusCode : Return service unavailable status code
func GetServiceUnavailableStatusCode() int32 {
	return serviceUnavailable
}

//GetDeniedStatusCode : Return denied status code
func GetDeniedStatusCode() int32 {
	return denied
}

//GetDuplicateInfoCode : Return duplicate information code
func GetDuplicateInfoCode() int32 {
	return duplicateInfo
}

//GetNoErrorCode : Return no error code
func GetNoErrorCode() int32 {
	return noError
}

//GetSaltLength : Return the salt length
func GetSaltLength() int {
	return passwordSaltLength
}

//GetTokenLength : Return the length of the token
func GetTokenLength() int {
	return tokenLength
}

//GetInvalidInformation : Return invalid argument
func GetInvalidInformation() int32 {
	return invalidInformation
}

//GetAuthCodeMaxNum : Return the max number
func GetAuthCodeMaxNum() int64 {
	return authCodeMaxNum
}

//GetAuthCodeMinNum : Return the min authen number
func GetAuthCodeMinNum() int64 {
	return authCodeMinNum
}

//GetGmailSMTPServer : Return Gmail SMTP server
func GetGmailSMTPServer() string {
	return gmailSMTPServer
}

//GetGmailSMTPHost : Return Gmail SMTP Host
func GetGmailSMTPHost() string {
	return gmailSMTPHost
}

//GetAppEmailAddress : Return Gmail email address
func GetAppEmailAddress() string {
	return gmailEmailAddress
}

//GetDatabaseName : Return get the name of the main database name
func GetDatabaseName() string {
	return databaseName
}

//GetAppName : Return app name
func GetAppName() string {
	return appName
}

//GetAccountInfoCollection : Return the name of the account info collection
func GetAccountInfoCollection() string {
	return accountInfoCollection
}

//GetAccessTokenCollection : Return get the name for access token collection
func GetAccessTokenCollection() string {
	return accessTokenCollection
}

//GetRefreshTokenCollection : Return get the name for refresh token collection
func GetRefreshTokenCollection() string {
	return refreshTokenCollection
}

//GetStoreImageCollection : Return the name of store image collection
func GetStoreImageCollection() string {
	return storeImageCollection
}

//GetPostDescriptionCollection : Return the name of post description collection
func GetStorePostDescriptionCollection() string {
	return storePostDescription
}

//GetAccountVerification : Return get the name for account verification collection
func GetAccountAuthenticationCollection() string {
	return accountAuthenticationCollection
}

//GetUserDatabase : Return get the user database name
func GetUserDatabase() string {
	return userDatabase
}

//GetPostDatabase : Return the name of post database
func GetPostDatabase() string {
	return postDatabase
}

//GetAppKey : Return get the application
func GetAppKey() string {
	return appKey
}
