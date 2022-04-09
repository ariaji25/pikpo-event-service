# PIKPO Events Service

A GRPC Service for create events with support recurring event

# Generate rpc from proto

Enter the proto file directory
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative {file.proto} 

```

# How To Run
1. Run the database
``` docker-compose up -d db```
2. Run the migrations
``` docker-compose up migration```
3. Run the service
``` docker-compose up -d event```
4. Run the test
```go test```