package main

import (
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/quanghia24/fishystore/gen-go/user"
	"github.com/quanghia24/fishystore/gen-go/yob"
)

func runServer(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, yobClient *yob.YOBServiceClient, addr string) error {
	transport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		return err
	}

	handler := NewUserHandler(yobClient)
	processor := user.NewUserServiceProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting User service... on ", addr)
	return server.Serve()
}
