package main

import (
	"net"

	"github.com/cloudwego/kitex/server"
	item "github.com/quanghia24/fishystore/kitex_gen/fishystore/item/itemservice"

	"log"
)

func main() {
	itemServiceImpl := new(ItemServiceImpl)
	stockCli, err := NewStockClient("0.0.0.0:8890")
	if err != nil {
		log.Fatal(err)
	}
	itemServiceImpl.stockCli = stockCli

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8889")
	svr := item.NewServer(itemServiceImpl, server.WithServiceAddr(addr))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
