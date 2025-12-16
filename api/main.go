package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/quanghia24/fishystore/api/utils"
	"github.com/quanghia24/fishystore/gen-go/user"
)

var (
	userClient *user.UserServiceClient
)

// initUserClient initializes the Thrift client connection to the RPC server
func initUserClient() (*user.UserServiceClient, error) {
	// Create socket transport to RPC server running on localhost:8001
	transport := thrift.NewTSocketConf("localhost:8001", nil)
	protocolFactory := thrift.NewTBinaryProtocolFactoryConf(nil)

	// Open connection
	if err := transport.Open(); err != nil {
		return nil, fmt.Errorf("failed to open transport: %w", err)
	}

	// Create Thrift client
	client := user.NewUserServiceClient(
		thrift.NewTStandardClient(
			protocolFactory.GetProtocol(transport),
			protocolFactory.GetProtocol(transport),
		),
	)

	log.Println("Successfully connected to RPC server at localhost:8001")
	return client, nil
}

// GET /user?id={id}
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		err := fmt.Errorf("missing 'id' parameter")
		utils.ReponseError(w, err, http.StatusBadRequest)
		return
	}

	// Convert to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ReponseError(w, err, http.StatusBadRequest)
		return
	}

	// Create RPC request
	req := &user.GetUserReq{
		ID: id,
	}

	// Call RPC service
	resp, err := userClient.GetUser(context.Background(), req)
	if err != nil {
		utils.ReponseError(w, err, http.StatusInternalServerError)
		return
	}

	// Marshal response to JSON
	body, err := json.Marshal(resp)
	if err != nil {
		utils.ReponseError(w, err, http.StatusInternalServerError)
		return
	}

	utils.ReponseStatus(w, body, http.StatusOK)
}

// healthHandler handles GET /health requests for health checks
func healthHandler(w http.ResponseWriter, r *http.Request) {
	msg := `{"status": "healthy"}`
	utils.ReponseStatus(w, []byte(msg), http.StatusOK)
}

func main() {
	// Initialize Thrift client connection
	var err error
	userClient, err = initUserClient()
	if err != nil {
		log.Fatalf("Failed to initialize user client: %v", err)
	}

	// Register HTTP handlers
	http.HandleFunc("/user", getUserHandler)
	http.HandleFunc("/health", healthHandler)

	// Start HTTP server
	addr := ":8080"
	log.Printf("Starting API Gateway server on %s", addr)
	log.Printf("Endpoints:")
	log.Printf("  GET /user?id={id} - Get user by ID")
	log.Printf("  GET /health       - Health check")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
