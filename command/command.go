package command

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
