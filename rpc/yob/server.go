package main

import (
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/quanghia24/fishystore/gen-go/yob"
)

func runServer(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, addr string) error {
	transport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		return err
	}

	handler := NewYOBHandler()
	processor := yob.NewYOBServiceProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting Year of Birth service... on ", addr)
	return server.Serve()
}
