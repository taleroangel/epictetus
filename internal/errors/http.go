package errors

import (
	"fmt"
	"net/http"
)

// Return HTTP status error
type HttpStatusError struct {
	Code int
	Err  string
}

func (e HttpStatusError) Error() string {
	return fmt.Sprintf("(%d) %s", e.Code, e.Err)
}

func NotFoundHttpError() HttpStatusError {
	return HttpStatusError{
		Code: http.StatusNotFound,
		Err:  "The given resource was not found",
	}
}

func BadRequestHttpError() HttpStatusError {
	return HttpStatusError{
		Code: http.StatusBadRequest,
		Err:  "Invalid request: Check verb, body or header for missing or invalid information",
	}
}
