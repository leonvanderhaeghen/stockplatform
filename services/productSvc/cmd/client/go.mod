module github.com/leonvanderhaeghen/stockplatform/services/productSvc/cmd/client

go 1.23.0

toolchain go1.23.3

require (
	github.com/leonvanderhaeghen/stockplatform v0.0.0
	google.golang.org/grpc v1.72.2
	google.golang.org/protobuf v1.36.6
)

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250505200425-f936aa4a68b2 // indirect
)

replace github.com/leonvanderhaeghen/stockplatform => ../../../../
