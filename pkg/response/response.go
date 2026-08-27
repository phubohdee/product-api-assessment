package response

type Response struct {
	Successful bool   `json:"successful"`
	ErrorCode  string `json:"error_code"`
	Data       any    `json:"data"`
}

func Success(data any) Response {
	return Response{
		Successful: true,
		ErrorCode:  "",
		Data:       data,
	}
}

func Error(errorCode string) Response {
	return Response{
		Successful: false,
		ErrorCode:  errorCode,
		Data:       nil,
	}
}
