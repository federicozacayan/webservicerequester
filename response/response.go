package response

type GeneralResponse struct {
	Responses interface{} `json:"responses"`
	Payload   interface{} `json:"payload"`
}

type Response struct {
	Headers    []Header    `json:"headers"`
	Body       interface{} `json:"body"`
	StatusCode *int        `json:"statusCode"`
}

type Header struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

func (r *Response) GetHeaders() []Header {
	return r.Headers
}

func (r *Response) GetBody() interface{} {
	return r.Body
}
