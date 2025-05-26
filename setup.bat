@echo off
echo Setting up the development environment...

:: Set up GOPATH
set GOPATH=%USERPROFILE%\go
set PATH=%PATH%;%GOPATH%\bin

:: Create necessary directories
if not exist pkg\gen\go\product\v1 mkdir pkg\gen\go\product\v1

:: Generate protobuf files
echo Generating protobuf files...
protoc -I=./proto ^
  --go_out=./pkg/gen/go/product/v1 ^
  --go_opt=paths=source_relative ^
  --go-grpc_out=./pkg/gen/go/product/v1 ^
  --go-grpc_opt=paths=source_relative ^
  ./proto/product/v1/product.proto

:: Fix import paths in generated files
echo Fixing import paths...
powershell -Command "Get-ChildItem -Path .\pkg\gen -Recurse -Filter *.go | ForEach-Object { (Get-Content $_.FullName) -replace '\"stockplatform\\/', '\"github.com/leonvanderhaeghen/stockplatform/' | Set-Content $_.FullName }"

:: Install dependencies
echo Installing dependencies...
cd services\productSvc
go mod tidy

cd ..\..
echo Setup complete!
