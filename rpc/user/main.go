package main

import (
	"fmt"
	"log"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/quanghia24/fishystore/gen-go/yob"
)

// initYOBClient initializes the Thrift client connection to the RPC server
func initYOBClient() (*yob.YOBServiceClient, thrift.TTransport, error) {
	// Create socket transport to RPC server running on localhost:8002
	transport := thrift.NewTSocketConf("localhost:8002", nil)
	protocolFactory := thrift.NewTBinaryProtocolFactoryConf(nil)

	// Open connection
	if err := transport.Open(); err != nil {
		return nil, nil, fmt.Errorf("failed to open transport: %w", err)
	}

	// Create Thrift client
	client := yob.NewYOBServiceClient(
		thrift.NewTStandardClient(
			protocolFactory.GetProtocol(transport),
			protocolFactory.GetProtocol(transport),
		),
	)

	log.Println("Successfully connected to RPC server at localhost:8002")
	return client, transport, nil
}

func main() {
	yobClient, yobTransport, err := initYOBClient()
	if err != nil {
		log.Fatal("failed to initialize YOB client:", err)
	}
	defer yobTransport.Close()

	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTBinaryProtocolFactoryConf(nil)

	if err := runServer(transportFactory, protocolFactory, yobClient, "localhost:8001"); err != nil {
		fmt.Println("error running server:", err)
	}
}
