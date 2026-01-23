package httpx

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"

	apperrors "github.com/vadxq/go-rest-starter/pkg/errors"
	"github.com/vadxq/go-rest-starter/pkg/logger"
)

var log logger.Logger = logger.Default()

// SetLogger sets the logger used by the response helpers.
func SetLogger(l logger.Logger) {
	if l != nil {
		log = l
	}
}

// Response is the standard HTTP response payload.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

// ErrorDetail carries structured error metadata.
type ErrorDetail struct {
	Type   string   `json:"type"`
	Fields []string `json:"fields,omitempty"`
}

// JSON writes a success response.
func JSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	payload := Response{
		Code:    status,
		Message: "OK",
		Data:    data,
		TraceID: logger.GetTraceID(r.Context()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.WithContext(r.Context()).Error("failed to encode response JSON", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Error writes an error response.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	appErr := apperrors.AsError(err)
	status := appErr.StatusCode()

	detail := ErrorDetail{Type: string(appErr.Type)}
	if fields := extractValidationFields(appErr.Err); len(fields) > 0 {
		detail.Fields = fields
	}

	payload := Response{
		Code:    status,
		Message: appErr.Message,
		Data:    detail,
		TraceID: logger.GetTraceID(r.Context()),
	}

	ctxLogger := log.WithContext(r.Context())
	if status >= http.StatusInternalServerError {
		ctxLogger.Error(appErr.Message, "error", appErr, "type", string(appErr.Type))
	} else {
		ctxLogger.Debug(appErr.Message, "error", appErr, "type", string(appErr.Type))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if encodeErr := json.NewEncoder(w).Encode(payload); encodeErr != nil {
		ctxLogger.Error("failed to encode error response JSON", "error", encodeErr)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func extractValidationFields(err error) []string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		fields := make([]string, 0, len(validationErrors))
		for _, fieldErr := range validationErrors {
			fields = append(fields, fieldErr.Field())
		}
		return fields
	}
	return nil
}
