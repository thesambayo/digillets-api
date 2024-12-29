package publicid

import (
	"fmt"

	nanoid "github.com/matoous/go-nanoid/v2"
)

const (
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	length   = 13
)

// New generates a public id using nanoid
// it requires a prefix:
// returns such as usr_{randomFromNanoid}
func New(prefix string) (string, error) {
	id, err := nanoid.Generate(alphabet, length)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v%v", prefix, id), nil
}
