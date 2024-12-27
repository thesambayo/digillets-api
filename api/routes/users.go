package routes

import (
	"errors"
	"net/http"

	"github.com/thesambayo/digillet-api/api/contexts"
	"github.com/thesambayo/digillet-api/api/httpx"
	"github.com/thesambayo/digillet-api/internal/constants"
	"github.com/thesambayo/digillet-api/internal/data/users"
	"github.com/thesambayo/digillet-api/internal/publicid"
	"github.com/thesambayo/digillet-api/internal/validators"
)

func (routes *Routes) CreateUser(resWriter http.ResponseWriter, req *http.Request) {
	// an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the request body into the anonymous struct.
	err := routes.httpx.ReadJSON(resWriter, req, &input)
	if err != nil {
		routes.httpx.BadRequestResponse(resWriter, req, err)
		return
	}

	publicID, _ := publicid.New(constants.PrefixUserID)
	user := &users.User{
		Name:     input.Name,
		PublicID: publicID,
		Email:    input.Email,
	}

	// Use the Password.Set() method to generate and store the hashed and plaintext passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
		return
	}

	validator := validators.New()
	routes.models.Users.ValidateName(validator, user.Name)
	routes.models.Users.ValidateEmail(validator, user.Email)
	routes.models.Users.ValidatePasswordPlaintext(validator, input.Password)
	if !validator.Valid() {
		routes.httpx.FailedValidationResponse(resWriter, req, validator.Errors)
		return
	}

	user, err = routes.models.Users.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our failedValidationResponse() helper.
		case errors.Is(err, constants.ErrDuplicateEmail):
			validator.AddError("email", "a user with this email address already exists")
			routes.httpx.FailedValidationResponse(resWriter, req, validator.Errors)
		default:
			routes.httpx.ServerErrorResponse(resWriter, req, err)
		}
		return
	}

	err = routes.httpx.WriteJSON(
		resWriter,
		http.StatusOK,
		httpx.Envelope{"message": "user registered successfully", "data": user},
		nil,
	)

	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
	}
}

func (routes *Routes) GetUserProfile(resWriter http.ResponseWriter, req *http.Request) {
	user := contexts.ContextGetUser(req)

	err := routes.httpx.WriteJSON(
		resWriter,
		http.StatusOK,
		httpx.Envelope{"message": "user profile fetched successfully", "data": user},
		nil,
	)

	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
	}
}
