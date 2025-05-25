module testclient

go 1.21

require (
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
	stockplatform v0.0.0
)

replace stockplatform => .
