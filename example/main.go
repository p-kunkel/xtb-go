package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/p-kunkel/xtb-go/client"
)

func main() {
	c, err := client.New(client.Config{
		Host:            client.Host,
		Mode:            client.ModeDemo,
		RequestInterval: client.MinReuqestInterval,
	})
	if err != nil {
		log.Fatalln(err)
	}

	lResp, err := c.Login(context.Background(), client.LoginRequest{
		UserId:   os.Getenv("XTB_USER_ID"),
		Password: os.Getenv("XTB_PASSWORD"),
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", lResp)

	resp, err := c.GetCurrentUserData()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%+v\n", resp)

	streamResp, unsubscribe := c.StreamGetKeepAlive()

	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println("unsubscribe")
		unsubscribe()
	}()

	for {
		r, ok := <-streamResp
		if !ok {
			break
		}

		fmt.Printf("stream:  %+v\n", r)
	}

	if err := c.Logout(); err != nil {
		log.Fatalln(err)
	}
}
