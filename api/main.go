package main

import (
	"context"
	"log"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/quanghia24/fishystore/kitex_gen/fishystore/item"
	"github.com/quanghia24/fishystore/kitex_gen/fishystore/item/itemservice"
)

var (
	itemCli itemservice.Client
)

func main() {
	itemClient, err := itemservice.NewClient("fishystore", client.WithHostPorts("0.0.0.0:8889"))
	if err != nil {
		log.Fatal(err)
	}
	itemCli = itemClient

	hz := server.New(server.WithHostPorts("localhost:8888"))

	hz.GET("/api/item", Handler)

	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

func Handler(ctx context.Context, c *app.RequestContext) {
	req := item.NewGetItemReq()
	req.Id = 1024
	resp, err := itemCli.GetItem(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, resp)
}
