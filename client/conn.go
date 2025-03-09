package client

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type connAttribute struct {
	req  Request
	resp chan Response
}

type conn struct {
	storage         *sync.Map
	pipe            chan connAttribute
	RequestInterval time.Duration
}

func newConn(minReuqestInterval time.Duration) *conn {
	return &conn{
		RequestInterval: minReuqestInterval,
		pipe:            make(chan connAttribute),
		storage:         new(sync.Map),
	}
}

func (c *conn) Read(ctx context.Context, conn *websocket.Conn) {
	var (
		resp Response
		err  error
	)

	for {
		select {
		case <-ctx.Done():
			conn.Close()
			return
		default:
			resp = Response{}

			if err = conn.ReadJSON(&resp); err != nil {
				c.storage.Clear()
				log.Printf("conn listen err: %s\n", err)
				return
			}

			if r, ok := c.storage.Load(resp.CustomTag); ok {
				respChan := r.(chan Response)
				respChan <- resp

				close(respChan)
				c.storage.Delete(resp.CustomTag)
			}
		}
	}
}

func (c *conn) Write(ctx context.Context, conn *websocket.Conn) {
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

			c.storage.Store(msg.req.CustomTag, msg.resp)

		case <-ctx.Done():
			conn.Close()
			return
		}
	}
}

func (c *conn) Send(req Request) Response {
	r := make(chan Response)
	c.pipe <- connAttribute{
		req:  req,
		resp: r,
	}
	return <-r
}
