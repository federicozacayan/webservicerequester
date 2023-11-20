package actions

import (
	"bot/mapping"
	"bot/request"
	"bot/requester"
	"bot/response"
	"encoding/json"
	"fmt"
)

func Exec(body []byte) []byte {
	request := &request.Requests{}
	err := json.Unmarshal(body, request)

	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	responses := map[string]response.Response{}

	str := []byte{}
	for _, step := range request.Actions {
		if step.Type == nil {
			fmt.Println("Error: type is nil")
			return nil
		}
		if *step.Type == "request" {
			fmt.Println("step.Type == request")
			requestName := *step.Name
			response := requester.ExecRequestByName(requestName, request.Requests)
			responses[requestName] = *response
		}
		if *step.Type == "mapping" {
			fmt.Println("step.Type == mapping")
			str = mapping.ExecMappingByName(*step.Name, request, responses)

		}
	}
	generalResponse := response.GeneralResponse{
		Responses: responses,
		Payload:   string(str),
	}
	//convert generalResponse to []byte
	responsesByte, err4 := json.Marshal(generalResponse)
	if err4 != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return nil
	}
	return responsesByte
}
