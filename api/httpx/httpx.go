// httpx is a simple contraction of `http extensions
// a collection of methods and useful tools for HTTP handling.
package httpx

import "github.com/thesambayo/digillets-api/internal/jsonlog"

type Utils struct {
	logger *jsonlog.Logger
}

func New(logger *jsonlog.Logger) *Utils {
	return &Utils{logger: logger}
}
