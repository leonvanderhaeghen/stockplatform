module github.com/leonvanderhaeghen/stockplatform/cmd/test_grpc_client

go 1.21

require (
	github.com/leonvanderhaeghen/stockplatform/services/productSvc v0.1.0
	google.golang.org/grpc v1.62.0
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
)

replace github.com/leonvanderhaeghen/stockplatform/services/productSvc => ../../services/productSvc
