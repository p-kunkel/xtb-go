package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/p-kunkel/xtb-go/command"
)

var (
	HOST = "ws.xtb.com"

	MODE_DEMO = "demo"
	MODE_REAL = "real"

	MIN_REQUEST_INTERVAL = 200 * time.Millisecond
)

type Client struct {
	conn            *Connection
	streamSessionId string
}

type Config struct {
	Host            string
	Mode            string
	RequestInterval time.Duration
}

type Request struct {
	Command   string      `json:"command"`
	Arguments interface{} `json:"arguments,omitempty"`
	CustomTag string      `json:"customTag,omitempty"`
}

func NewRequestWithRandomTag(cmd string, arg interface{}) Request {
	return Request{
		Command:   cmd,
		Arguments: arg,
		CustomTag: uuid.NewString(),
	}
}

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

func New(conf Config) *Client {
	return &Client{
		conn: &Connection{
			conf:      conf,
			send:      make(chan requestResponse),
			closeConn: make(chan bool),
		},
	}
}

func (c *Client) Do(req Request) (Response, error) {
	if req.CustomTag == "" {
		req.CustomTag = uuid.NewString()
	}

	r := requestResponse{
		req:  req,
		resp: make(chan Response),
	}

	c.conn.send <- r
	resp := <-r.resp

	if !resp.Status {
		return Response{}, resp.ErrorResp
	}

	return resp, nil
}

// Do request and bind response.ReturnData to result
func (c *Client) DoAndBind(req Request, result interface{}) error {
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if len(resp.ReturnData) > 0 {
		if err := json.Unmarshal([]byte(resp.ReturnData), result); err != nil {
			return nil
		}
	}

	return nil
}

func (c *Client) Login(ctx context.Context, lr command.LoginRequest) (command.LoginResponse, error) {
	if err := c.conn.dial(ctx); err != nil {
		return command.LoginResponse{}, err
	}

	resp, err := c.Do(NewRequestWithRandomTag("login", lr))
	if err != nil {
		return command.LoginResponse{}, err
	}

	c.streamSessionId = resp.StreamSessionId
	result := command.LoginResponse{
		StreamSessionId: resp.StreamSessionId,
	}

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			<-ticker.C
			if !c.conn.isConnected {
				break
			}
			c.Ping()
		}
	}()

	return result, nil
}

func (c *Client) Logout() error {
	if _, err := c.Do(NewRequestWithRandomTag("logout", nil)); err != nil {
		return err
	}

	c.conn.Close()
	return nil
}

func (c *Client) Ping() error {
	_, err := c.Do(NewRequestWithRandomTag("ping", nil))
	return err
}

func (c *Client) GetCurrentUserData() (command.GetCurrentUserDataResponse, error) {
	result := command.GetCurrentUserDataResponse{}
	if err := c.DoAndBind(NewRequestWithRandomTag("getCurrentUserData", nil), &result); err != nil {
		return command.GetCurrentUserDataResponse{}, err
	}
	return result, nil
}
