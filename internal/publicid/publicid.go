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
// it requires a prefix such as can be found in the gist below:
// https://gist.github.com/fnky/76f533366f75cf75802c8052b577e2a5
// returns such as usr_{randomFromNanoid}
func New(prefix string) (string, error) {
	id, err := nanoid.Generate(alphabet, length)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v%v", prefix, id), nil
}
