package client

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Request struct {
	Command   string `json:"command"`
	Arguments any    `json:"arguments,omitempty"`
	CustomTag string `json:"customTag,omitempty"`
}

func NewRequest(cmd string, arg any) Request {
	return Request{
		Command:   cmd,
		Arguments: arg,
	}
}

func (req Request) WithRandomTag() Request {
	req.CustomTag = uuid.NewString()
	return req
}

// Do request and bind response data value pointed to by resp
func DoAndBind[T any](c *Client, req Request, v T) (T, error) {
	r, err := c.Do(req)
	if err != nil {
		return v, err
	}

	if len(r.ReturnData) > 0 {
		if err := json.Unmarshal([]byte(r.ReturnData), v); err != nil {
			return v, err
		}
	}

	return v, nil
}
