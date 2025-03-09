package client

import (
	"context"

	"github.com/gorilla/websocket"
)

type connHandler struct {
	conn *websocket.Conn
	ctx  context.Context
	r    reader
	w    writter
}

type writter interface {
	Write(context.Context, *websocket.Conn)
}

type reader interface {
	Read(context.Context, *websocket.Conn)
}

func dial(addr string, opts ...option) error {
	c := connHandler{
		ctx: context.Background(),
	}

	for _, f := range opts {
		c = f(c)
	}

	conn, _, err := websocket.DefaultDialer.DialContext(c.ctx, addr, nil)
	if err != nil {
		return err
	}

	c.conn = conn

	c.handle()
	return nil
}

func (ch connHandler) handle() {
	go ch.r.Read(ch.ctx, ch.conn)
	go ch.w.Write(ch.ctx, ch.conn)
}

type option func(connHandler) connHandler

func withContext(ctx context.Context) option {
	return func(ch connHandler) connHandler {
		ch.ctx = ctx
		return ch
	}
}

func withWritter(w writter) option {
	return func(ch connHandler) connHandler {
		ch.w = w
		return ch
	}
}

func withReader(r reader) option {
	return func(ch connHandler) connHandler {
		ch.r = r
		return ch
	}
}
