package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/imtiyazkumar/students-api/internal/types"
	"github.com/imtiyazkumar/students-api/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		slog.Info("creating a new student")

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("body is empty")))
			return
		}

		if err != nil {

			// Catch type mismatch (like string → int)
			var ute *json.UnmarshalTypeError
			if errors.As(err, &ute) {
				msg := fmt.Sprintf("field %s must be of type %s", ute.Field, ute.Type)
				response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New(msg)))
				return
			}

			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request validation
		if err := validator.New().Struct(student); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		slog.Info("creating a new student")

		response.WriteJson(w, http.StatusCreated, map[string]interface{}{
			"message": "student created successfully",
			"student": student,
		})
	}
}
