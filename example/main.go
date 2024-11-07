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

	if err := c.Dial(context.Background()); err != nil {
		log.Fatalln(err)
	}

	if _, err := command.Login(c, command.LoginData{
		UserId:   os.Getenv("XTB_USER_ID"),
		Password: os.Getenv("XTB_PASSWORD"),
	}); err != nil {
		log.Fatalln(err)
	}

	resp, err := command.GetCurrentUserData(c)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v\n", resp)

	c.CloseConnection()
}
