package command

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/p-kunkel/xtb-go/client"
)

type Command struct {
	Name      string      `json:"command,omitempty"`
	Arguments interface{} `json:"arguments,omitempty"`
	CustomTag string      `json:"customTag,omitempty"`
}

func (cmd *Command) Run(c *client.Client, result interface{}) error {
	resp, err := c.Do(client.Request{
		Command:   cmd.Name,
		Arguments: cmd.Arguments,
		CustomTag: uuid.NewString(),
	})

	if err != nil {
		return err
	}

	if len(resp.ReturnData) > 0 {
		return json.Unmarshal([]byte(resp.ReturnData), result)
	}
	return nil
}

type LoginData struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

type Session struct {
	StreamSessionId string `json:"streamSessionId"`
}

func Login(c *client.Client, ld LoginData) (Session, error) {
	resp := Session{}
	cmd := Command{
		Name:      "login",
		Arguments: ld,
		CustomTag: uuid.NewString(),
	}

	if err := cmd.Run(c, &resp); err != nil {
		return Session{}, err
	}

	return resp, nil
}

type UserData struct {
	CompanyUnit        int     `json:"companyUnit"`
	Currency           string  `json:"currency"`
	Group              string  `json:"group"`
	IbAccount          bool    `json:"ibAccount"`
	Leverage           int     `json:"leverage"`
	LeverageMultiplier float64 `json:"leverageMultiplier"`
	SpreadType         string  `json:"spreadType"`
	TrailingStop       bool    `json:"trailingStop"`
}

func GetCurrentUserData(c *client.Client) (UserData, error) {
	resp := UserData{}
	cmd := Command{
		Name: "getCurrentUserData",
	}

	if err := cmd.Run(c, &resp); err != nil {
		return UserData{}, nil
	}

	return resp, nil
}
