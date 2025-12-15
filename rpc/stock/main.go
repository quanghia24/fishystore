package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	stock "github.com/quanghia24/fishystore/kitex_gen/fishystore/stock/stockservice"
)

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8890")
	svr := stock.NewServer(new(StockServiceImpl), server.WithServiceAddr(addr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
