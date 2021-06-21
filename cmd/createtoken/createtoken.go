package main

import (
	"context"
	"log"

	"github.com/sakti/o11ygo/api"
	"github.com/sakti/o11ygo/pkg/client"
	"github.com/sakti/o11ygo/pkg/config"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Println("creating token...")
	conf, err := config.LoadConfig()
	checkErr(err)
	c, err := client.NewAuthClient(conf.AuthsvcPort)
	checkErr(err)
	r, err := c.CreateToken(context.Background(), &api.CreateTokenRequest{})
	checkErr(err)
	log.Println(r.Token)
}
