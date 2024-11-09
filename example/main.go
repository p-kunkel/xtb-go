package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/p-kunkel/xtb-go/client"
	"github.com/p-kunkel/xtb-go/command"
)

func main() {
	c := client.New(client.Config{
		Host:            client.HOST,
		Mode:            client.MODE_DEMO,
		RequestInterval: client.MIN_REQUEST_INTERVAL,
	})

	lResp, err := c.Login(context.Background(), command.LoginRequest{
		UserId:   os.Getenv("XTB_USER_ID"),
		Password: os.Getenv("XTB_PASSWORD"),
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v\n", lResp)

	resp, err := c.GetCurrentUserData()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v\n", resp)

	if err := c.Logout(); err != nil {
		log.Fatalln(err)
	}
}
