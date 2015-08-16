# go-pb-grpc
Go + Protocol Buffers + gRPC client/server example

To compile and run the server, simply:

    $ go run server/main.go

Likewise, to run the client:

    $ go run client/main.go

If you wish to alter the .proto file, install [proto3](https://github.com/google/protobuf) then do:

    $ protoc --go_out=plugins=grpc:message *.proto
