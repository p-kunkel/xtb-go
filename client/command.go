package client

type LoginRequest struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
	// (optional, deprecated)
	AppId string `json:"appId,omitempty"`
	// (optional) application name
	AppName string `json:"appName,omitempty"`
}

type LoginResponse struct {
	StreamSessionId string `json:"streamSessionId"`
}

func login(c *Client, req LoginRequest) (*LoginResponse, error) {
	resp, err := c.Do(NewRequest("login", req).WithRandomTag())
	if err != nil {
		return nil, err
	}
	return &LoginResponse{StreamSessionId: resp.StreamSessionId}, nil
}

type GetCurrentUserDataResponse struct {
	CompanyUnit        int     `json:"companyUnit"`
	Currency           string  `json:"currency"`
	Group              string  `json:"group"`
	IbAccount          bool    `json:"ibAccount"`
	Leverage           int     `json:"leverage"`
	LeverageMultiplier float64 `json:"leverageMultiplier"`
	SpreadType         string  `json:"spreadType"`
	TrailingStop       bool    `json:"trailingStop"`
}

func getCurrentUserData(c *Client) (*GetCurrentUserDataResponse, error) {
	return DoAndBind(c,
		NewRequest("getCurrentUserData", nil).WithRandomTag(),
		&GetCurrentUserDataResponse{},
	)
}
