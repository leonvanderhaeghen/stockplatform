@echo off
echo Generating protobuf files...

:: Create necessary directories
if not exist pkg\gen\go\product\v1 mkdir pkg\gen\go\product\v1

:: Generate protobuf files
protoc -I=./proto ^
  --go_out=./pkg/gen/go/product/v1 ^
  --go_opt=paths=source_relative ^
  --go-grpc_out=./pkg/gen/go/product/v1 ^
  --go-grpc_opt=paths=source_relative ^
  ./proto/product/v1/product.proto

echo Protobuf files generated successfully!
