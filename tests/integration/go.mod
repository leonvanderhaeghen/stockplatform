module github.com/leonvanderhaeghen/stockplatform/tests/integration

go 1.21

require (
	github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/product/v1 v0.0.0
	github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/user/v1 v0.0.0
	github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/order/v1 v0.0.0
	google.golang.org/grpc v1.58.3
)

replace github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/product/v1 => ../../pkg/gen/go/product/v1
replace github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/user/v1 => ../../pkg/gen/go/user/v1
replace github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/order/v1 => ../../pkg/gen/go/order/v1
