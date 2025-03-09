package client

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type streamResponse struct {
	Command string          `json:"command"`
	Data    json.RawMessage `json:"data,omitempty"`
	ErrorResp
}

type streamAttribute struct {
	storeKey string
	req      any
	resp     chan streamResponse
}

type connStream struct {
	storage         *sync.Map
	pipe            chan streamAttribute
	RequestInterval time.Duration
}

func newConnStream(minReuqestInterval time.Duration) *connStream {
	return &connStream{
		RequestInterval: minReuqestInterval,
		pipe:            make(chan streamAttribute),
		storage:         new(sync.Map),
	}
}

func (c *connStream) Read(ctx context.Context, conn *websocket.Conn) {
	var (
		resp streamResponse
		err  error
	)

	for {
		select {
		case <-ctx.Done():
			conn.Close()
			return
		default:
			resp = streamResponse{}

			if err = conn.ReadJSON(&resp); err != nil {
				c.storage.Clear()
				log.Printf("stream listen err: %s\n", err)
				return
			}

			if r, ok := c.storage.Load(resp.Command); ok {
				respChan := r.(chan streamResponse)
				respChan <- resp
			}
		}
	}
}

func (c *connStream) Write(ctx context.Context, conn *websocket.Conn) {
	requestIntervalTicker := time.NewTicker(c.RequestInterval)
	defer requestIntervalTicker.Stop()

	for {
		select {
		case msg := <-c.pipe:
			<-requestIntervalTicker.C

			if err := conn.WriteJSON(msg.req); err != nil {
				log.Fatalf("send err: %s", err)
				return
			}

			if msg.storeKey != "" {
				c.storage.Store(msg.storeKey, msg.resp)
			}

		case <-ctx.Done():
			conn.Close()
			return
		}
	}
}

func (c *connStream) Send(req any, storeKey string) chan streamResponse {
	r := make(chan streamResponse)
	c.pipe <- streamAttribute{
		storeKey: storeKey,
		req:      req,
		resp:     r,
	}
	return r
}

func (c *connStream) unsubscribe(req any, storeKey string) {
	ch := c.Send(req, "")
	close(ch)

	if r, ok := c.storage.LoadAndDelete(storeKey); ok {
		respChan := r.(chan streamResponse)
		close(respChan)
	}
}
