module github.com/leonvanderhaeghen/stockplatform/services/gatewaySvc

go 1.23.0

toolchain go1.23.3

require (
	github.com/gin-contrib/cors v1.5.0
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/leonvanderhaeghen/stockplatform v0.1.0
	github.com/leonvanderhaeghen/stockplatform/services/inventorySvc v0.0.0-20250617235535-5a86d542f1f1
	github.com/leonvanderhaeghen/stockplatform/services/orderSvc v0.0.0-20250617235535-5a86d542f1f1
	github.com/leonvanderhaeghen/stockplatform/services/productSvc v0.0.0-20250617235535-5a86d542f1f1
	github.com/leonvanderhaeghen/stockplatform/services/storeSvc v0.0.0-20250617235535-5a86d542f1f1
	github.com/leonvanderhaeghen/stockplatform/services/supplierSvc v0.0.0-20250617235535-5a86d542f1f1
	github.com/leonvanderhaeghen/stockplatform/services/userSvc v0.0.0-20250617235535-5a86d542f1f1
	github.com/spf13/viper v1.20.1
	github.com/swaggo/files v1.0.1
	github.com/swaggo/gin-swagger v1.6.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.73.0
)

replace github.com/leonvanderhaeghen/stockplatform => ../..

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/bytedance/sonic v1.10.2 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.1 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/spec v0.20.9 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.16.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/swaggo/swag v1.16.2 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.6.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250505200425-f936aa4a68b2 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
