// utilities are related to HTTP handling.
package httpx

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ReadIDParam gets id from request url
func (utils *Utils) ReadIDParam(req *http.Request) string {
	return req.PathValue("id")
}

// ReadIntIDParam gets id from request url and converts to an int64
func (utils *Utils) ReadIntIDParam(req *http.Request) (int64, error) {
	id, err := strconv.Atoi(req.PathValue("id"))
	// id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return int64(id), nil
}

// The ReadString() helper returns a string value from the query string, or the provided
// default value if no matching key could be found.
func (utils *Utils) ReadString(queryString url.Values, key string, defaultValue string) string {
	// Extract the value for a given key from the query string. If no key exists this
	// will return the empty string "".
	queryStringValue := queryString.Get(key)

	// If no key exists (or the value is empty) then return the default value.
	if queryStringValue == "" {
		return defaultValue
	}

	// Otherwise return the string.
	return queryStringValue
}

// The ReadCSV() helper reads a string value from the query string and then splits it
// into a slice on the comma character. If no matching key could be found, it returns
// the provided default value.
//
// for queryURL values like url?status=HIGH,URGENT to return []string{"HIGH", "URGENT"}
func (utils *Utils) ReadCSV(queryString url.Values, key string, defaultValue []string) []string {
	// Extract the value from the query string.
	csv := queryString.Get(key)

	// If no key exists (or the value is empty) then return the default value.
	if csv == "" {
		return defaultValue
	}

	// Otherwise parse the value into a []string slice and return it.
	return strings.Split(csv, ",")
}

// The ReadInt() helper reads a string value from the query string and converts it to an
// integer before returning. If no matching key could be found it returns the provided
// default value. If the value couldn't be converted to an integer, then we record an
// error message in the provided Validator instance.
// func ReadInt(queryString url.Values, key string, defaultValue int, validator *validators.Validator) int {
// 	// Extract the value from the query string.
// 	queryStringValue := queryString.Get(key)

// 	// If no key exists (or the value is empty) then return the default value.
// 	if queryStringValue == "" {
// 		return defaultValue
// 	}

// 	// Try to convert the value to an int. If this fails, add an error message to the
// 	// validator instance and return the default value.
// 	queryIntValue, err := strconv.Atoi(queryStringValue)
// 	if err != nil {
// 		validator.AddError(key, "must be an integer value")
// 		return defaultValue
// 	}

// 	// Otherwise, return the converted integer value.
// 	return queryIntValue
// }
