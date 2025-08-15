module github.com/leonvanderhaeghen/stockplatform/services/userSvc

go 1.23.0

toolchain go1.23.3

require github.com/leonvanderhaeghen/stockplatform v0.0.0-00010101000000-000000000000

replace github.com/leonvanderhaeghen/stockplatform => ../..

require (
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/google/uuid v1.6.0
	go.mongodb.org/mongo-driver v1.17.4
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.38.0
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/leonvanderhaeghen/stockplatform/services/orderSvc v0.0.0-20250617235535-5a86d542f1f1 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250505200425-f936aa4a68b2 // indirect
)
