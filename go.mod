module github.com/leonvanderhaeghen/stockplatform

go 1.23.0

toolchain go1.23.3

// This is a workspace module - see go.work for more details

require (
	github.com/leonvanderhaeghen/stockplatform/services/inventorySvc v0.0.0-20250617235535-5a86d542f1f1
	github.com/leonvanderhaeghen/stockplatform/services/orderSvc v0.0.0-20250617235535-5a86d542f1f1
	github.com/leonvanderhaeghen/stockplatform/services/productSvc v0.0.0-20250617235535-5a86d542f1f1
	github.com/leonvanderhaeghen/stockplatform/services/storeSvc v0.0.0-20250617235535-5a86d542f1f1
	github.com/leonvanderhaeghen/stockplatform/services/supplierSvc v0.0.0-20250617235535-5a86d542f1f1
	github.com/leonvanderhaeghen/stockplatform/services/userSvc v0.0.0-20250617235535-5a86d542f1f1
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.73.0
)

require (
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250505200425-f936aa4a68b2 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
