package client

import (
	"encoding/json"
	"log"
)

func Subscribe[T any](c *Client, req any, storeKey string) <-chan T {
	respChan := make(chan T)

	resp := c.streamConn.Send(req, storeKey)

	go func(rChan chan T, srChan <-chan streamResponse) {
		for {
			var respData T
			r, ok := <-srChan
			if !ok {
				close(rChan)
				return
			}

			if err := json.Unmarshal(r.Data, &respData); err != nil {
				log.Printf("unmarshal err for %T: %s", respData, err)
			}

			rChan <- respData
		}
	}(respChan, resp)

	return respChan
}

type baseStreamRequest struct {
	Command         string `json:"command"`
	StreamSessionId string `json:"streamSessionId,omitempty"`
}

type GetKeepAliveResponse struct {
	Timestamp int64 `json:"timestamp"`
}

func streamGetKeepAlive(c *Client) (stream <-chan GetKeepAliveResponse, unsubscribe func()) {
	storeKey := "keepAlive"
	return Subscribe[GetKeepAliveResponse](c,
			baseStreamRequest{
				StreamSessionId: c.streamSessionId,
				Command:         "getKeepAlive",
			},
			storeKey,
		),
		func() {
			c.streamConn.unsubscribe(baseStreamRequest{Command: "stopKeepAlive"}, storeKey)
		}
}

type GetBalanceResponse struct {
	Balance     float64 `json:"balance"`
	Credit      float64 `json:"credit"`
	Equity      float64 `json:"equity"`
	Margin      float64 `json:"margin"`
	MarginFree  float64 `json:"marginFree"`
	MarginLevel float64 `json:"marginLevel"`
}

func streamGetBalance(c *Client) (stream <-chan GetBalanceResponse, unsubscribe func()) {
	storeKey := "balance"
	return Subscribe[GetBalanceResponse](c,
			baseStreamRequest{
				StreamSessionId: c.streamSessionId,
				Command:         "getBalance",
			},
			storeKey,
		),
		func() {
			c.streamConn.unsubscribe(baseStreamRequest{Command: "stopBalance"}, storeKey)
		}
}

type GetCandlesRequest struct {
	Symbol string `json:"symbol"`
}

type GetCandlesResponse struct {
	Close     float64 `json:"close"`
	Ctm       int64   `json:"ctm"`
	CtmString string  `json:"ctmString"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Open      float64 `json:"open"`
	QuoteId   Quote   `json:"quoteId"`
	Symbol    string  `json:"symbol"`
	Vol       float64 `json:"vol"`
}

func streamGetCandles(c *Client, req GetCandlesRequest) (stream <-chan GetCandlesResponse, unsubscribe func()) {
	storeKey := "candle"
	return Subscribe[GetCandlesResponse](c,
			struct {
				baseStreamRequest
				GetCandlesRequest
			}{
				baseStreamRequest: baseStreamRequest{
					StreamSessionId: c.streamSessionId,
					Command:         "getCandles",
				},
				GetCandlesRequest: req,
			},
			storeKey,
		),
		func() {
			c.streamConn.unsubscribe(baseStreamRequest{Command: "stopCandles"}, storeKey)
		}
}

type GetNewsResponse struct {
	Body  string `json:"body"`
	Key   string `json:"key"`
	Time  int64  `json:"time"`
	Title string `json:"title"`
}

func streamGetNews(c *Client) (stream <-chan GetNewsResponse, unsubscribe func()) {
	storeKey := "news"
	return Subscribe[GetNewsResponse](c,
			baseStreamRequest{
				StreamSessionId: c.streamSessionId,
				Command:         "getNews",
			},
			storeKey,
		),
		func() {
			c.streamConn.unsubscribe(baseStreamRequest{Command: "stopNews"}, storeKey)
		}
}

type GetProfitsResponse struct {
	Order    int64   `json:"order"`
	Order2   int64   `json:"order2"`
	Position int64   `json:"position"`
	Profit   float64 `json:"profit"`
}

func streamGetProfits(c *Client) (stream <-chan GetProfitsResponse, unsubscribe func()) {
	storeKey := "profit"
	return Subscribe[GetProfitsResponse](c,
			baseStreamRequest{
				StreamSessionId: c.streamSessionId,
				Command:         "getProfits",
			},
			storeKey,
		),
		func() {
			c.streamConn.unsubscribe(baseStreamRequest{Command: "stopProfits"}, storeKey)
		}
}

type GetTickPricesRequest struct {
	Symbol         string `json:"symbol"`
	MinArrivalTime int64  `json:"minArrivalTime"`
	MaxLevel       int64  `json:"maxLevel"`
}

type GetTickPricesResponse struct {
	Ask         float64 `json:"ask"`
	AskVolume   int64   `json:"askVolume"`
	Bid         float64 `json:"bid"`
	BidVolume   int64   `json:"bidVolume"`
	High        float64 `json:"high"`
	Level       int64   `json:"level"`
	Low         float64 `json:"low"`
	QuoteId     Quote   `json:"quoteId"`
	SpreadRaw   float64 `json:"spreadRaw"`
	SpreadTable float64 `json:"spreadTable"`
	Symbol      string  `json:"symbol"`
	Timestamp   int64   `json:"timestamp"`
}

func streamGetTickPrices(c *Client, req GetTickPricesRequest) (stream <-chan GetTickPricesResponse, unsubscribe func()) {
	storeKey := "candle"
	return Subscribe[GetTickPricesResponse](c,
			struct {
				baseStreamRequest
				GetTickPricesRequest
			}{
				baseStreamRequest: baseStreamRequest{
					StreamSessionId: c.streamSessionId,
					Command:         "getTickPrices",
				},
				GetTickPricesRequest: req,
			},
			storeKey,
		),
		func() {
			c.streamConn.unsubscribe(
				struct {
					baseStreamRequest
					Symbol string `json:"symbol"`
				}{
					baseStreamRequest: baseStreamRequest{Command: "stopTickPrices"},
					Symbol:            req.Symbol,
				},
				storeKey)
		}
}

type GetTradesResponse struct {
	OpenPrice float64 `json:"open_price"`
	OpenTime  int64   `json:"open_time"`

	ClosePrice float64 `json:"close_price"`
	CloseTime  int64   `json:"close_time"`
	Closed     bool    `json:"closed"`

	Order         int64   `json:"order"`
	Order2        int64   `json:"order2"`
	Cmd           int64   `json:"cmd"`
	Comment       string  `json:"comment"`
	Commission    float64 `json:"commission"`
	CustomComment string  `json:"customComment"`
	Digits        int64   `json:"digits"`
	Expiration    int64   `json:"expiration"`
	MarginRate    float64 `json:"margin_rate"`
	Offset        int64   `json:"offset"`
	Position      int64   `json:"position"`
	Profit        float64 `json:"profit"`
	Sl            float64 `json:"sl"`
	State         string  `json:"state"`
	Storage       float64 `json:"storage"`
	Symbol        string  `json:"symbol"`
	Tp            float64 `json:"tp"`
	Type          int64   `json:"type"`
	Volume        float64 `json:"volume"`
}

func streamGetTrades(c *Client) (stream <-chan GetTradesResponse, unsubscribe func()) {
	storeKey := "trade"
	return Subscribe[GetTradesResponse](c,
			baseStreamRequest{
				StreamSessionId: c.streamSessionId,
				Command:         "getTrades",
			},
			storeKey,
		),
		func() {
			c.streamConn.unsubscribe(baseStreamRequest{Command: "stopTrades"}, storeKey)
		}
}

type GetTradeStatusResponse struct {
	CustomComment string  `json:"customComment"`
	Message       string  `json:"message"`
	Order         int64   `json:"order"`
	Price         float64 `json:"price"`
	// TODO: it's an enum
	RequestStatus int64 `json:"requestStatus"`
}

func streamGetTradeStatus(c *Client) (stream <-chan GetTradeStatusResponse, unsubscribe func()) {
	storeKey := "tradeStatus"
	return Subscribe[GetTradeStatusResponse](c,
			baseStreamRequest{
				StreamSessionId: c.streamSessionId,
				Command:         "getTradeStatus",
			},
			storeKey,
		),
		func() {
			c.streamConn.unsubscribe(baseStreamRequest{Command: "stopTradeStatus"}, storeKey)
		}
}
