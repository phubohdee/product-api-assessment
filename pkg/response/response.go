package response

// Response is the standard API response format.
type Response struct {
	Successful bool   `json:"successful"`
	ErrorCode  string `json:"error_code"`
	Data       any    `json:"data"`
}

// Success returns a successful response with data.
func Success(data any) Response {
	return Response{
		Successful: true,
		ErrorCode:  "",
		Data:       data,
	}
}

// Error returns an error response with the given error code.
func Error(errorCode string) Response {
	return Response{
		Successful: false,
		ErrorCode:  errorCode,
		Data:       nil,
	}
}
