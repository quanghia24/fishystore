# API Implementation Plan for api/main.go

## Project Architecture Analysis

### Current Structure
This is a **Thrift-based microservices** project with the following components:

1. **IDL Definitions** (`idl/`):
   - [`base.thrift`](idl/base.thrift:1) - Base response structure
   - [`user.thrift`](idl/user.thrift:1) - User service definitions
   - [`yob.thrift`](idl/yob.thrift:1) - Year of Birth service definitions

2. **Generated Code** (`gen-go/`):
   - Auto-generated Thrift bindings for Go
   - [`user.UserServiceClient`](gen-go/user/user.go:511) - Client for calling UserService

3. **RPC Server** (`rpc/user/`):
   - [`main.go`](rpc/user/main.go:1) - Starts Thrift server on `localhost:8000`
   - [`server.go`](rpc/user/server.go:1) - Server configuration
   - [`handler.go`](rpc/user/handler.go:1) - UserService implementation

4. **API Gateway** (`api/`):
   - [`main.go`](api/main.go:1) - HTTP REST API that should connect to the Thrift RPC server

### Architecture Flow
```
HTTP Client → API Gateway (port 8080) → Thrift RPC Server (port 8000) → Business Logic
```

## Issues in Current api/main.go

### Critical Errors:
1. **Line 17**: `userClient = *user.NewUserServiceClient()` - Missing transport/protocol arguments
2. **Line 16**: Function is named `main()` but package is `api` - should be in `package main`
3. **Line 20**: `http.Handle("/user", getUserHandler)` - Handler is a function, needs `http.HandlerFunc()`
4. **Line 22-24**: Missing imports for `fmt`, `html`, and `log`
5. **Line 34**: `userClient.GetUser()` - Returns `(*GetUserResp, error)` not just response
6. **Line 35**: `json.Marshal()` - Returns `([]byte, error)` not just bytes
7. **Line 36**: `r.Response(body)` - Should be `w.Write(body)` and handle HTTP response properly

### Missing Components:
- Thrift client connection setup (transport, protocol)
- Error handling throughout
- Proper HTTP response headers
- Request parameter extraction (user ID should come from query/path params)
- Connection lifecycle management

## Implementation Plan

### Step 1: Fix Package and Imports
- Change package from `api` to `main`
- Add all required imports:
  - `fmt`, `log` for logging
  - `net/http` for HTTP server
  - `strconv` for parsing parameters
  - `github.com/apache/thrift/lib/go/thrift` for Thrift client
  - `github.com/quanghia24/fishystore/gen-go/user` for generated types

### Step 2: Initialize Thrift Client Connection
Create a proper Thrift client initialization:
```go
func initUserClient() (*user.UserServiceClient, error) {
    // Create socket transport to RPC server
    transport := thrift.NewTSocketConf("localhost:8000", nil)

    // Create protocol factory
    protocolFactory := thrift.NewTBinaryProtocolFactoryConf(nil)

    // Open connection
    if err := transport.Open(); err != nil {
        return nil, err
    }

    // Create client
    client := user.NewUserServiceClient(
        thrift.NewTStandardClient(
            protocolFactory.GetProtocol(transport),
            protocolFactory.GetProtocol(transport),
        ),
    )

    return client, nil
}
```

### Step 3: Implement getUserHandler Properly
Fix the handler to:
- Extract user ID from query parameters
- Call RPC service with proper error handling
- Marshal JSON response correctly
- Set proper HTTP headers
- Handle errors with appropriate HTTP status codes

```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Parse user ID from query params
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
        return
    }

    // 2. Convert to int64
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
        return
    }

    // 3. Create request
    req := &user.GetUserReq{ID: id}

    // 4. Call RPC service
    resp, err := userClient.GetUser(context.Background(), req)
    if err != nil {
        http.Error(w, fmt.Sprintf("RPC error: %v", err), http.StatusInternalServerError)
        return
    }

    // 5. Marshal to JSON
    body, err := json.Marshal(resp)
    if err != nil {
        http.Error(w, fmt.Sprintf("JSON marshal error: %v", err), http.StatusInternalServerError)
        return
    }

    // 6. Send response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(body)
}
```

### Step 4: Fix main() Function
- Initialize Thrift client with error handling
- Register handlers properly using `http.HandlerFunc()`
- Remove the `/bar` test endpoint (or keep if needed)
- Add graceful shutdown handling

### Step 5: Add Connection Management
- Consider connection pooling for production use
- Add defer cleanup for transport closure
- Implement health check endpoint

## API Design Recommendations

### Endpoints:
1. **GET /user?id={id}** - Get user by ID (current implementation)
2. **GET /health** - Health check endpoint (recommended addition)

### Response Format:
```json
{
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "yob": "1990"
  },
  "baseResp": {
    "code": "0",
    "msg": "Success"
  }
}
```

### Error Responses:
```json
{
  "error": "Error message here"
}
```

## Testing Plan

1. **Start RPC Server**: Run `rpc/user/main.go` on port 8000
2. **Start API Gateway**: Run `api/main.go` on port 8080
3. **Test Endpoints**:
   - `curl http://localhost:8080/user?id=1`
   - `curl http://localhost:8080/health`

## Dependencies

All dependencies are already in [`go.mod`](go.mod:1):
- `github.com/apache/thrift v0.22.0`
- Go 1.25.3

## Notes

- The RPC server handler currently returns `nil, nil` in [`handler.go`](rpc/user/handler.go:17-19), so actual user data won't be returned until that's implemented
- Consider adding middleware for logging, CORS, and authentication
- Connection pooling should be implemented for production use
- Add graceful shutdown for the HTTP server

## Next Steps

Once you approve this plan, I'll switch to **Code mode** to implement all the fixes and improvements.