package client

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

type requestResponse struct {
	req  Request
	resp chan Response
}

type Connection struct {
	conn        *websocket.Conn
	conf        Config
	send        chan requestResponse
	receive     chan Response
	closeConn   chan bool
	isConnected bool
}

func (c *Connection) dial(ctx context.Context) error {
	var (
		err    error
		addr   = fmt.Sprintf("wss://%s/%s", c.conf.Host, c.conf.Mode)
		origin = fmt.Sprintf("https://%s/", c.conf.Host)
	)

	cfg, err := websocket.NewConfig(addr, origin)
	if err != nil {
		return err
	}
	c.conn, err = cfg.DialContext(ctx)
	if err != nil {
		return err
	}

	c.isConnected = true
	go c.handle(ctx)
	return nil
}

func (c *Connection) handle(ctx context.Context) {
	go c.read(ctx)

	requestIntervalTicker := time.NewTicker(c.conf.RequestInterval)
	defer requestIntervalTicker.Stop()

	waitingForResponse := map[string]chan Response{}

	c.ping()
	for {
		select {
		case msg := <-c.send:
			<-requestIntervalTicker.C

			if err := c.write(msg.req); err != nil {
				c.receive <- Response{
					Status:    false,
					CustomTag: msg.req.CustomTag,
					ErrorResp: ErrorResp{
						ErrorCode:  "write_loop",
						ErrorDescr: err.Error(),
					},
				}
			}

			waitingForResponse[msg.req.CustomTag] = msg.resp

		case resp := <-c.receive:
			if respChan, ok := waitingForResponse[resp.CustomTag]; ok {
				respChan <- resp

				close(respChan)
				delete(waitingForResponse, resp.CustomTag)
			}

		case <-time.After(time.Second * 10):
			c.ping()

		case <-ctx.Done():
			c.Close()

		case <-c.closeConn:
			close(c.send)
			close(c.receive)
			c.conn.Close()
			return
		}
	}
}

func (c *Connection) ping() {
	c.write(Request{
		Command: "ping",
	})
}

func (c *Connection) read(ctx context.Context) {
	b := Response{}
	for c.isConnected {
		b = Response{}
		if err := websocket.JSON.Receive(c.conn, &b); err != nil {
			c.Close()
			break
		}

		c.receive <- b
	}
	c.closeConn <- true
}

func (c *Connection) write(r Request) error {
	return websocket.JSON.Send(c.conn, r)
}

func (c *Connection) Close() {
	c.isConnected = false
}
