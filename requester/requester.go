package requester

import (
	"bot/request"
	"bot/response"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ExecRequestByName(name string, requests []request.Request) *response.Response {
	for _, req := range requests {
		if *req.Name == name {
			return call(req)
		}
	}
	return nil
}

func concatHeaderValues(values []string) string {
	// Concatenate multiple header values
	return strings.Join(values, ", ")
}

func call(req request.Request) *response.Response {
	protocol := req.Data.GetProtocol()
	host := req.Data.GetHost()
	path := req.Data.GetPath()
	port := req.Data.GetPort()
	method := req.Data.GetMethod()
	body := req.Data.GetBody()

	var resp *http.Response
	var err error
	switch method {
	case "get":
		resp, err = http.Get(fmt.Sprintf("%s://%s:%s%s", protocol, host, port, path))
	case "post":
		resp, err = http.Post(fmt.Sprintf("%s://%s:%s%s", protocol, host, port, path),
			"application/json",
			strings.NewReader(fmt.Sprintf("%s", body)))
	}
	if err != nil {
		fmt.Println("Error al realizar la solicitud:", err)
		return nil
	}
	defer resp.Body.Close()

	body2, err2 := io.ReadAll(io.Reader(resp.Body))
	if err2 != nil {
		fmt.Println("Error al leer la respuesta:", err)
		return nil
	}
	responseObj := response.Response{
		Body: string(body2),
	}

	for key, values := range resp.Header {
		for _, value := range values {
			k := fmt.Sprintf("%s", key)
			v := fmt.Sprintf("%s", value)
			header123 := response.Header{
				Key:   &k,
				Value: &v,
			}
			responseObj.Headers = append(responseObj.Headers, header123)
		}
	}
	responseObj.StatusCode = &resp.StatusCode
	return &responseObj
}
