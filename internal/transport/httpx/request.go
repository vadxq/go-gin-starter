package httpx

import (
	"encoding/json"
	"net/http"

	apperrors "github.com/vadxq/go-rest-starter/pkg/errors"
)

// DecodeJSON decodes JSON payload from request body.
func DecodeJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return apperrors.BadRequestError("invalid JSON payload", err)
	}
	return nil
}

// BindJSON decodes JSON and validates it if validate is provided.
func BindJSON(r *http.Request, v interface{}, validate func(interface{}) error) error {
	if err := DecodeJSON(r, v); err != nil {
		return err
	}

	if validate != nil {
		if err := validate(v); err != nil {
			return apperrors.ValidationError("validation failed", err)
		}
	}

	return nil
}
