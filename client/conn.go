package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type requestResponse struct {
	req  Request
	resp chan Response
}

type Connection struct {
	conn               *websocket.Conn
	conf               Config
	send               chan requestResponse
	waitingForResponse sync.Map
	closeConn          chan bool
	isConnected        bool
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

	go c.readLoop(ctx)
	go c.writeLoop(ctx)
	return nil
}

func (c *Connection) writeLoop(ctx context.Context) {
	requestIntervalTicker := time.NewTicker(c.conf.RequestInterval)
	defer requestIntervalTicker.Stop()

	for {
		select {
		case msg := <-c.send:
			<-requestIntervalTicker.C

			if err := websocket.JSON.Send(c.conn, msg.req); err != nil {
				msg.resp <- Response{
					Status:    false,
					CustomTag: msg.req.CustomTag,
					ErrorResp: ErrorResp{
						ErrorCode:  "write_loop",
						ErrorDescr: err.Error(),
					},
				}
				continue
			}

			c.waitingForResponse.Store(msg.req.CustomTag, msg.resp)

		case <-ctx.Done():
			c.Close()

		case <-c.closeConn:
			close(c.send)
			c.waitingForResponse.Clear()
			c.conn.Close()
			return
		}
	}
}

func (c *Connection) readLoop(ctx context.Context) {
	var (
		resp Response
		err  error
	)

	for c.isConnected {
		resp = Response{}

		if err = websocket.JSON.Receive(c.conn, &resp); err != nil {
			c.waitingForResponse.Range(func(key, value any) bool {
				resp := Response{
					Status:    false,
					CustomTag: key.(string),
					ErrorResp: ErrorResp{
						ErrorCode:  "read_loop",
						ErrorDescr: err.Error(),
					},
				}

				respChan := value.(chan Response)
				respChan <- resp
				close(respChan)
				return true
			})

			c.waitingForResponse.Clear()
			c.Close()
			break
		}

		if r, ok := c.waitingForResponse.Load(resp.CustomTag); ok {
			respChan := r.(chan Response)
			respChan <- resp

			close(respChan)
			c.waitingForResponse.Delete(resp.CustomTag)
		}
	}
	c.closeConn <- true
}

func (c *Connection) Close() {
	c.isConnected = false
}
