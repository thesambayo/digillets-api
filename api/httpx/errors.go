package httpx

import (
	"fmt"
	"net/http"
)

// The logError() method is a generic helper for logging an error message.
// Later, we'll upgrade this to use structured logging, and record additional information
// about the request including the HTTP method and URL.
func (utils *Utils) logError(req *http.Request, err error) {
	utils.logger.PrintError(err, map[string]string{
		"request_method": req.Method,
		"request_url":    req.URL.String(),
	})
}

// The ErrorResponse() method is a generic helper for sending JSON-formatted error
// messages to the client with a given status code. Note that we're using an interface{}
// type for the message parameter, rather than just a string type, as this gives us
// more flexibility over the values that we can include in the response.
func (utils *Utils) ErrorResponse(resWriter http.ResponseWriter, req *http.Request, status int, message interface{}) {
	env := Envelope{"error": message}
	// Write the response using the writeJSON() helper. If this happens to return an
	// error then log it, and fall back to sending the client an empty response with a
	// 500 Internal Server Error status code.
	err := utils.WriteJSON(resWriter, status, env, nil)
	if err != nil {
		utils.logError(req, err)
		resWriter.WriteHeader(http.StatusInternalServerError)
	}
}

// The ServerErrorResponse() method will be used when our application encounters an
// unexpected problem at runtime. It logs the detailed error message, then uses the
// errorResponse() helper to send a 500 Internal Server Error status and JSON
// response (containing a generic error message) to the client.”
func (utils *Utils) ServerErrorResponse(resWriter http.ResponseWriter, req *http.Request, err error) {
	utils.logError(req, err)
	message := "the server encountered a problem and could not process your request"
	utils.ErrorResponse(resWriter, req, http.StatusInternalServerError, message)
}

// The NotFoundResponse() method will be used to send a 404 Not Found status code and
// JSON response to the client.
func (utils *Utils) NotFoundResponse(resWriter http.ResponseWriter, req *http.Request) {
	message := "the requested resource could not be found"
	utils.ErrorResponse(resWriter, req, http.StatusNotFound, message)
}

// The MethodNotAllowedResponse() method will be used to send a 405 Method Not Allowed
// status code and JSON response to the client.
func (utils *Utils) MethodNotAllowedResponse(resWriter http.ResponseWriter, req *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", req.Method)
	utils.ErrorResponse(resWriter, req, http.StatusMethodNotAllowed, message)
}

func (utils *Utils) BadRequestResponse(resWriter http.ResponseWriter, req *http.Request, err error) {
	utils.ErrorResponse(resWriter, req, http.StatusBadRequest, err.Error())
}

// Note that the errors parameter here has the type map[string]string, which is exactly
// the same as the errors map contained in our Validator type.
func (utils *Utils) FailedValidationResponse(resWriter http.ResponseWriter, req *http.Request, errors map[string]string) {
	utils.ErrorResponse(resWriter, req, http.StatusUnprocessableEntity, errors)
}

func (utils *Utils) EditConflictResponse(resWriter http.ResponseWriter, req *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	utils.ErrorResponse(resWriter, req, http.StatusConflict, message)
}

func (utils *Utils) RateLimitExceededResponse(resWriter http.ResponseWriter, req *http.Request) {
	message := "rate limit exceeded"
	utils.ErrorResponse(resWriter, req, http.StatusTooManyRequests, message)
}

func (utils *Utils) InvalidCredentialsResponse(resWriter http.ResponseWriter, req *http.Request) {
	message := "invalid authentication credentials"
	utils.ErrorResponse(resWriter, req, http.StatusUnauthorized, message)
}

func (utils *Utils) InvalidAuthenticationTokenResponse(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	utils.ErrorResponse(resWriter, req, http.StatusUnauthorized, message)
}

func (utils *Utils) AuthenticationRequiredResponse(resWriter http.ResponseWriter, req *http.Request) {
	message := "you must be authenticated to access this resource"
	utils.ErrorResponse(resWriter, req, http.StatusUnauthorized, message)
}

// func InactiveAccountResponse(resWriter http.ResponseWriter, req *http.Request) {
// 	message := "your staff account must be activated to access this resource"
// 	ErrorResponse(resWriter, req, http.StatusForbidden, message)
// }

// func NotPermittedResponse(resWriter http.ResponseWriter, req *http.Request) {
// 	message := "your staff account doesn't have the necessary permissions to access this resource"
// 	ErrorResponse(resWriter, req, http.StatusForbidden, message)
// }
