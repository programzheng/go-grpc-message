# go-grpc-message
## generate *.pb.go
```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative message.proto
```