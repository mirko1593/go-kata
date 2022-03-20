# Reflection

This example shows how reflection can be registered on a gRPC server.

See
https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md
for a tutorial.


# Try it

```go
go run server/main.go
```

```bash
grpcurl --plaintext localhost:50051 list
grpcurl -cert "" -key "" localhost:50551 list
```

There are multiple existing reflection clients.

To use `gRPC CLI`, follow
https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md#grpc-cli.

To use `grpcurl`, see https://github.com/fullstorydev/grpcurl.
