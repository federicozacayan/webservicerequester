package request

import "strings"

/*
{
    "url": "test",
    "method": "post",
    "headers": [
        {
            "key": "Content-Type",
            "value": "application/json"
        }
    ],
    "body": {
        "keyExample": "valueExample"
    }
}
*/
type Requests struct {
	Requests []Request `json:"requests"`
	Actions  []Action  `json:"actions"`
	Mappings []Mapping `json:"mappings"`
}

type Mapping struct {
	Name   *string       `json:"name"`
	Params MappingParams `json:"params"`
}

type MappingParams struct {
	Objects     MappingObject       `json:"objects"`
	Expressions []MappingExpression `json:"expressions"`
}

type MappingObject struct {
	Source *string `json:"source"`
	Target *string `json:"target"`
}

type MappingExpression struct {
	Source *string `json:"source"`
	Target *string `json:"target"`
}

// Response *string `json:"response"`

/*
{
	"response": "hello",
	"jsonPath": "$.store.book[1].title"
}
*/
type Action struct {
	Type *string `json:"type"`
	Name *string `json:"name"`
	// Response *string `json:"response"`
	// JsonPath *string `json:"jsonPath"`
}

type Request struct {
	Name *string `json:"name"`
	Data *Data   `json:"data"`
}
type Data struct {
	URL      *string     `json:"url"`
	Host     *string     `json:"host"`
	Path     *string     `json:"path"`
	Protocol *string     `json:"protocol"`
	Port     *string     `json:"port"`
	Method   *string     `json:"method"`
	Headers  []Header    `json:"headers"`
	Body     interface{} `json:"body"`
}

type Header struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

func (r *Data) GetMethod() string {
	if r.Method != nil {
		return strings.ToLower(*r.Method)
	}
	return "get"
}

func (r *Data) GetProtocol() string {
	if r.Protocol != nil {
		return *r.Protocol
	}
	return "http"
}

func (r *Data) GetPort() string {
	if r.Port != nil {
		return *r.Port
	}
	return "80"
}

func (r *Data) GetHost() string {
	if r.Host != nil {
		return *r.Host
	}
	return "localhost"
}

func (r *Data) GetPath() string {
	if r.Path != nil {
		return *r.Path
	}
	return "/"
}

func (r *Data) GetURL() (string, error) {
	if r.URL != nil {
		return *r.URL, nil
	}
	return "", &CustomError{Message: "URL is nil"}
}

func (r *Data) GetBody() interface{} {
	if r.Body != nil {
		return r.Body
	}
	return ""
}
