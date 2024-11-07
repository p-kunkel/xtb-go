package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

var (
	HOST = "ws.xtb.com"

	MODE_DEMO = "demo"
	MODE_REAL = "real"

	MIN_REQUEST_INTERVAL = 200 * time.Millisecond
)

type Client struct {
	conn                *Connection
	restartConnectionFn func() error
}

type Config struct {
	Host            string
	Mode            string
	RequestInterval time.Duration

	RestartConnectionFn func() error

	logger *log.Logger
}

type Request struct {
	Command   string      `json:"command,omitempty"`
	Arguments interface{} `json:"arguments,omitempty"`
	CustomTag string      `json:"customTag,omitempty"`
}

type Response struct {
	Status     bool            `json:"status"`
	CustomTag  string          `json:"customTag"`
	ReturnData json.RawMessage `json:"returnData,omitempty"`
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
			receive:   make(chan Response),
			closeConn: make(chan bool),
		},
		restartConnectionFn: conf.RestartConnectionFn,
	}
}

func (c *Client) Do(req Request) (Response, error) {
	if !c.conn.isConnected {
		if c.restartConnectionFn == nil {
			return Response{}, errors.New("connection has broken")
		}

		if err := c.restartConnectionFn(); err != nil {
			return Response{}, fmt.Errorf("restart connection fail: %w", err)
		}
	}

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

func (c *Client) Dial(ctx context.Context) error {
	return c.conn.dial(ctx)
}

func (c *Client) CloseConnection() {
	c.conn.Close()
}
