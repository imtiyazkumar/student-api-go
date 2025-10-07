package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOk    = "ok"
	StatusError = "error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		field := err.Field()
		value := err.Value()
		param := err.Param()

		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("%s is required", field))

		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be a valid email address", value))

		case "gte":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be greater than or equal to %s", field, param))

		case "lte":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be less than or equal to %s", field, param))

		case "gt":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be greater than %s", field, param))

		case "lt":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be less than %s", field, param))

		case "len":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be exactly %s characters long", field, param))

		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must have a minimum value of %s", field, param))

		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must have a maximum value of %s", field, param))

		case "numeric":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be a numeric value", field))

		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be a valid URL", field))

		case "uuid":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be a valid UUID", field))

		case "oneof":
			errMsgs = append(errMsgs, fmt.Sprintf("%s must be one of [%s]", field, param))

		default:
			errMsgs = append(errMsgs, fmt.Sprintf("%s is invalid (%s)", field, err.ActualTag()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
