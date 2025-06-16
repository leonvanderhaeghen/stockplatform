module github.com/leonvanderhaeghen/stockplatform/tests/integration

go 1.21

require (
	github.com/leonvanderhaeghen/stockplatform v0.0.0
	google.golang.org/grpc v1.58.3
)

replace github.com/leonvanderhaeghen/stockplatform => ../..
