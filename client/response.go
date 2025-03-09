package client

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Status          bool            `json:"status"`
	CustomTag       string          `json:"customTag"`
	ReturnData      json.RawMessage `json:"returnData,omitempty"`
	StreamSessionId string          `json:"streamSessionId,omitempty"`
	ErrorResp
}

type ErrorResp struct {
	ErrorCode  string `json:"errorCode"`
	ErrorDescr string `json:"errorDescr"`
}

func (e ErrorResp) Error() string {
	return fmt.Sprintf("error code: %s, desc: %s", e.ErrorCode, e.ErrorDescr)
}
