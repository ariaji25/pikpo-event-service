# PIKPO Events Service

A GRPC Service for creating events with support recurring event

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
```go test ./app/events/```

Run with native code without docker
1. Run the database first with docker like the first step above
2. Create an environment variable to define if it runs natively
``` export mode=native```
4. Run the service
``` go run main.go```
5. Run the Unit test
``` go test ./app/events```


Notes : To run the unit test make sure the service is running on the port 8000,