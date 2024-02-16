package errors

import "net/http"

func UnauthorizedHttpError() HttpStatusError {
	return HttpStatusError{
		Code: http.StatusUnauthorized,
		Err:  "You are not authenticated or your credentials are invalid",
	}
}

func ForbiddenHttpError() HttpStatusError {
	return HttpStatusError{
		Code: http.StatusForbidden,
		Err:  "You have no persmission to access this resource",
	}
}
