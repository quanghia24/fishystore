package main

import (
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
)

func main() {
	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTBinaryProtocolFactoryConf(nil)

	if err := runServer(transportFactory, protocolFactory, "localhost:8002"); err != nil {
		fmt.Println("error running server:", err)
	}
}
