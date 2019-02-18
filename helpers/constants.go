package helpers

const (
	ok                     = 200
	created                = 201
	accepted               = 202
	invalidArguments       = 204
	denied                 = 205
	found                  = 302
	badRequest             = 400
	unauthorized           = 402
	forbidden              = 403
	notFound               = 404
	notAcceptable          = 406
	requestTimeout         = 408
	internalServerError    = 500
	badGateway             = 502
	serviceUnavailable     = 503
	passwordSaltLength     = 200
	authCodeMaxNum         = 900000
	authCodeMinNum         = 100000
	gmailSMTPServer        = "smtp.gmail.com:587"
	gmailSMTPHost          = "smtp.gmail.com"
	gmailEmailAddress      = "naturae.outdoor@gmail.com"
	appName                = "Naturae"
	databaseName           = "Naturae-Server"
	accountInfoCollection  = "Account_Information"
	accessTokenCollection  = "Access_Token"
	refreshTokenCollection = "Refresh_Token"
)

//GetOkStatusCode : Return ok status code
func GetOkStatusCode() int16 {
	return ok
}

//GetCreatedStatusCode : Return created status code
func GetCreatedStatusCode() int16 {
	return created
}

//GetAcceptedStatusCode : Return accepted status code
func GetAcceptedStatusCode() int16 {
	return accepted
}

//GetFoundStatusCode : Return found status code
func GetFoundStatusCode() int16 {
	return found
}

//GetBadRequestStatusCode : Return bad request status code
func GetBadRequestStatusCode() int16 {
	return badRequest
}

//GetUnauthorizedStatusCode : Return unauthorized status code
func GetUnauthorizedStatusCode() int16 {
	return unauthorized
}

//GetForbiddenStatusCode : Return forbidden status code
func GetForbiddenStatusCode() int16 {
	return forbidden
}

//GetNotFoundStatusCode : Return not found status code
func GetNotFoundStatusCode() int16 {
	return notFound
}

//GetNotAcceptableStatusCode : Return not acceptable status code
func GetNotAcceptableStatusCode() int16 {
	return notAcceptable
}

//GetRequestTimeoutStatusCode : Return request timeout status code
func GetRequestTimeoutStatusCode() int16 {
	return requestTimeout
}

//GetInternalServerErrorStatusCode : Return internal server error status code
func GetInternalServerErrorStatusCode() int16 {
	return internalServerError
}

//GetBadGatewayStatusCode : Return bad gateway status code
func GetBadGatewayStatusCode() int16 {
	return badGateway
}

//GetServiceUnavailableStatusCode : Return service unavailable status code
func GetServiceUnavailableStatusCode() int16 {
	return serviceUnavailable
}

//GetDeniedStatusCode : Return denied status code
func GetDeniedStatusCode() int {
	return denied
}

//GetSaltLength : return the salt length
func GetSaltLength() int {
	return passwordSaltLength
}

//GetInvalidArgument : return invalid argument
func GetInvalidArgument() int16 {
	return invalidArguments
}

//GetAuthCodeMaxNum : the max number
func GetAuthCodeMaxNum() int64 {
	return authCodeMaxNum
}

//GetAuthCodeMinNum : the min
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

//GetGmailEmailAddress : Return Gmail email address
func GetGmailEmailAddress() string {
	return gmailEmailAddress
}

//GetDatabaseName : Get the name of the main database name
func GetDatabaseName() string {
	return databaseName
}

//GetAppName : return app name
func GetAppName() string {
	return appName
}

//GetAccountInfoColl : return the name of the account info collection
func GetAccountInfoCollection() string {
	return accountInfoCollection
}

//GetAccessTokenColl : get the name for access token collection
func GetAccessTokenCollection() string {
	return accessTokenCollection
}

//GetRefreshTokenColl : get the name for refresh token collection
func GetRefreshTokenCollection() string {
	return refreshTokenCollection
}
