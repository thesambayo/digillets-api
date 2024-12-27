package users

import (
	"regexp"
	"strings"

	"github.com/thesambayo/digillets-api/internal/validators"
)

func (userModel UserModel) ValidateName(validator *validators.Validator, name string) {
	validator.Check(name != "", "name", "must be provided")

	parts := strings.Fields(name)
	// Check if there are at least two parts
	validator.Check(len(parts) >= 2, "name", "must contain both first and last name")

	// Check if the first and second parts are not empty
	if len(parts) > 1 {
		validator.Check(parts[0] != "" && parts[1] != "", "name", "must not have empty first or last name")
	}
}

func (userModel UserModel) ValidateEmail(validator *validators.Validator, email string) {
	validator.Check(email != "", "email", "must be provided")
	validator.Check(validators.Matches(email, validators.EmailREGEX), "email", "must be a valid email address")
}

func (userModel UserModel) ValidatePasswordPlaintext(validator *validators.Validator, password string) {
	// Check if the password is provided
	validator.Check(password != "", "password", "must be provided")

	// Check password length
	validator.Check(len(password) >= 8, "password", "must be at least 8 characters long")
	validator.Check(len(password) <= 72, "password", "must not be more than 72 characters long")

	// Regular expressions for validations
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	numberRegex := regexp.MustCompile(`[0-9]`)
	symbolRegex := regexp.MustCompile(`[\W_]`) // Matches any non-word character, including symbols

	// Validate for at least one uppercase letter
	validator.Check(validators.Matches(password, uppercaseRegex), "password", "must contain at least one uppercase letter")

	// Validate for at least one lowercase letter
	validator.Check(validators.Matches(password, lowercaseRegex), "password", "must contain at least one lowercase letter")

	// Validate for at least one number
	validator.Check(validators.Matches(password, numberRegex), "password", "must contain at least one number")

	// Validate for at least one symbol
	validator.Check(validators.Matches(password, symbolRegex), "password", "must contain at least one symbol character")
}
