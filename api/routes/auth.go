package routes

import (
	"errors"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"github.com/thesambayo/digillet-api/api/httpx"
	"github.com/thesambayo/digillet-api/internal/constants"
	"github.com/thesambayo/digillet-api/internal/validators"
)

func (routes *Routes) AuthenticateUser(resWriter http.ResponseWriter, req *http.Request) {
	// Parse the email and password from the request body.
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := routes.httpx.ReadJSON(resWriter, req, &input)
	if err != nil {
		routes.httpx.BadRequestResponse(resWriter, req, err)
		return
	}

	// Validate the email and password provided by the client.
	validator := validators.New()
	routes.models.Users.ValidateEmail(validator, input.Email)

	if !validator.Valid() {
		routes.httpx.FailedValidationResponse(resWriter, req, validator.Errors)
		return
	}

	// Lookup the user record based on the email address.
	// If no matching user was found,
	// then we call the app.invalidCredentialsResponse() helper to send a 401
	// Unauthorized response to client
	user, err := routes.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRecordNotFound):
			routes.httpx.InvalidCredentialsResponse(resWriter, req)
		default:
			routes.httpx.ServerErrorResponse(resWriter, req, err)
		}
		return
	}

	// Check if the provided password matches the actual password for the user.
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
		return
	}

	if !match {
		routes.httpx.InvalidCredentialsResponse(resWriter, req)
		return
	}

	// Create a JWT claims struct containing the user ID as the subject, with
	// time of now and validity window of the next 24 hours.
	// We also set the issuer and audience to a unique identifier for our application.
	var claims jwt.Claims
	claims.Subject = user.PublicID
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	// claims.Issuer = "operations.oryoltd.org"
	// claims.Audiences = []string{"operations.oryoltd.org"}

	// Sign the JWT claims using the HMAC-SHA256 algorithm and the secret key from the application config.
	// This returns a []byte slice containing the JWT as a base64-encoded string.
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(routes.config.Jwt.Secret))
	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
		return
	}

	// Encode the token to JSON and send
	loginData := map[string]interface{}{
		"authentication_token": string(jwtBytes),
		"user":                 *user,
	}
	err = routes.httpx.WriteJSON(
		resWriter,
		http.StatusCreated,
		httpx.Envelope{"data": loginData, "message": "logged in successfully"},
		nil,
	)
	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
	}
}
