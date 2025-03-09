package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Config struct {
	Host            string
	Mode            string
	RequestInterval time.Duration
}

type Client struct {
	addr            string
	conn            *conn
	streamAddr      string
	streamConn      *connStream
	streamSessionId string
	cancelCtxFn     context.CancelFunc
}

func New(conf Config) (*Client, error) {
	streamMode := ""
	switch conf.Mode {
	case ModeDemo:
		streamMode = modeDemoStream
	case ModeReal:
		streamMode = modeRealStream
	default:
		return nil, errors.New("invalid mode")
	}

	return &Client{
		addr:       fmt.Sprintf("wss://%s/%s", conf.Host, conf.Mode),
		streamAddr: fmt.Sprintf("wss://%s/%s", conf.Host, streamMode),
		conn:       newConn(conf.RequestInterval),
		streamConn: newConnStream(conf.RequestInterval),
	}, nil
}

func (c *Client) Do(req Request) (Response, error) {
	if req.CustomTag == "" {
		req.CustomTag = uuid.NewString()
	}

	resp := c.conn.Send(req)

	if !resp.Status {
		return Response{}, resp.ErrorResp
	}

	return resp, nil
}

func (c *Client) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	ctx, c.cancelCtxFn = context.WithCancel(ctx)

	if err := dial(c.addr, withContext(ctx), withReader(c.conn), withWritter(c.conn)); err != nil {
		return nil, err
	}

	if err := dial(c.streamAddr, withContext(ctx), withReader(c.streamConn), withWritter(c.streamConn)); err != nil {
		return nil, err
	}

	resp, err := login(c, req)
	if err != nil {
		return nil, err
	}

	c.streamSessionId = resp.StreamSessionId

	go func(ctx context.Context) {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				c.Ping()
				<-ticker.C
			}
		}
	}(ctx)

	return resp, nil
}

func (c *Client) Logout() error {
	if _, err := c.Do(NewRequest("logout", nil).WithRandomTag()); err != nil {
		return err
	}

	c.cancelCtxFn()
	return nil
}

func (c *Client) Ping() error {
	_, err := c.Do(NewRequest("ping", nil).WithRandomTag())
	ch := c.streamConn.Send(baseStreamRequest{Command: "ping", StreamSessionId: c.streamSessionId}, "")
	close(ch)
	return err
}

func (c *Client) GetCurrentUserData() (*GetCurrentUserDataResponse, error) {
	return getCurrentUserData(c)
}

func (c *Client) StreamGetKeepAlive() (streamChan <-chan GetKeepAliveResponse, unsubscribe func()) {
	return streamGetKeepAlive(c)
}

func (c *Client) StreamGetBalance() (streamChan <-chan GetBalanceResponse, unsubscribe func()) {
	return streamGetBalance(c)
}

func (c *Client) StreamGetCandles(req GetCandlesRequest) (streamChan <-chan GetCandlesResponse, unsubscribe func()) {
	return streamGetCandles(c, req)
}

func (c *Client) StreamGetNews() (streamChan <-chan GetNewsResponse, unsubscribe func()) {
	return streamGetNews(c)
}

func (c *Client) StreamGetProfits() (streamChan <-chan GetProfitsResponse, unsubscribe func()) {
	return streamGetProfits(c)
}

func (c *Client) StreamGetTickPrices(req GetTickPricesRequest) (streamChan <-chan GetTickPricesResponse, unsubscribe func()) {
	return streamGetTickPrices(c, req)
}

func (c *Client) StreamGetTrades() (streamChan <-chan GetTradesResponse, unsubscribe func()) {
	return streamGetTrades(c)
}

func (c *Client) StreamGetTradeStatus() (streamChan <-chan GetTradeStatusResponse, unsubscribe func()) {
	return streamGetTradeStatus(c)
}
