package helper

const (
	ok                  = 200
	created             = 201
	accepted            = 202
	cancelled           = 203
	invalidArugment     = 204
	found               = 302
	badRequest          = 400
	unauthorized        = 402
	forbidden           = 403
	notFound            = 404
	notAcceptable       = 406
	requestTimeout      = 408
	internalServerError = 500
	badGateway          = 502
	serviceUnavilable   = 503
)

//GetOkStatusCode : Return ok status code
func GetOkStatusCode() int {
	return ok
}

//GetCreatedStatusCode : Return created status code
func GetCreatedStatusCode() int {
	return created
}

//GetAcceptedStatusCode : Return accepted status code
func GetAcceptedStatusCode() int {
	return accepted
}

//GetFoundStatusCode : Return found status code
func GetFoundStatusCode() int {
	return found
}

//GetBadRequestStatusCode : Return bad request status code
func GetBadRequestStatusCode() int {
	return badRequest
}

//GetUnauthorizedStatusCode : Return unauthorized status code
func GetUnauthorizedStatusCode() int {
	return unauthorized
}

//GetForbiddenStatusCode : Return forbidden status code
func GetForbiddenStatusCode() int {
	return forbidden
}

//GetNotFoundStatusCode : Return not found status code
func GetNotFoundStatusCode() int {
	return notFound
}

//GetNotAcceptableStatusCode : Return not acceptable status code
func GetNotAcceptableStatusCode() int {
	return notAcceptable
}

//GetRequestTimeoutStatusCode : Return request timeout status code
func GetRequestTimeoutStatusCode() int {
	return requestTimeout
}

//GetInternalServerErrorStatusCode : Return internal server error status code
func GetInternalServerErrorStatusCode() int {
	return internalServerError
}

//GetBadGatewayStatusCode : Return bad gateway status code
func GetBadGatewayStatusCode() int {
	return badGateway
}

//GetServiceUnavailableStatusCode : Return service unavailable status code
func GetServiceUnavailableStatusCode() int {
	return serviceUnavilable
}
