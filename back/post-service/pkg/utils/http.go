package utils

type HttpResponse struct {
	Status  int
	Headers map[string]string
	Body    interface{}
}

func NewHttpResponse(status int, body interface{}) *HttpResponse {
	var b interface{}
	if body != nil {
		b = body
	}
	return &HttpResponse{
		Status: status,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
			"Access-Control-Allow-Methods": "*",
		},
		Body: b,
	}
}

func Ok(body interface{}) *HttpResponse {
	return NewHttpResponse(200, body)
}

func NoContent() *HttpResponse {
	return NewHttpResponse(204, nil)
}

func BadRequest(body interface{}) *HttpResponse {
	return NewHttpResponse(400, body)
}

func MethodNotAllowed() *HttpResponse {
	return NewHttpResponse(405, nil)
}

func InternalServerError(body interface{}) *HttpResponse {
	return NewHttpResponse(500, body)
}
