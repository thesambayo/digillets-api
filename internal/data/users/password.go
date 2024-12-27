package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// password, a custom type which is a struct containing the plaintext and hashed
// versions of the password for a staff. The plaintext field is a *pointer* to a string,
// so that we're able to distinguish between a plaintext password not being present in
// the struct at all (which will be nil), versus a plaintext password which is the empty string "".
type password struct {
	plaintext *string
	hash      []byte
}

// The Set() method calculates the bcrypt hash of a plaintext password, and stores both
// the hash and the plaintext versions in the struct.
func (password *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	password.plaintext = &plaintextPassword
	password.hash = hash

	return nil
}

// The Matches() method checks whether the provided plaintext password matches the
// hashed password stored in the struct, returning true if it matches and false otherwise.
func (password *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(password.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
