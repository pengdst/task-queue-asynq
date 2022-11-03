package exceptions

type ErrorBadRequest struct {
	ErrMsg string
}

func NewErrorBadRequest(errMsg string) ErrorBadRequest {
	return ErrorBadRequest{ErrMsg: errMsg}
}

func (e *ErrorBadRequest) Error() string {
	return e.ErrMsg
}
