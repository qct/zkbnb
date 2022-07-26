module github.com/bnb-chain/zkbas

go 1.16

require (
	github.com/zeromicro/go-zero v1.3.4
	gorm.io/gorm v1.23.4
)

require (
	github.com/bnb-chain/zkbas-crypto v0.0.3-stable
	github.com/bnb-chain/zkbas-eth-rpc v0.0.2-stable
	github.com/consensys/gnark v0.7.0
	github.com/consensys/gnark-crypto v0.7.0
	github.com/eko/gocache/v2 v2.3.1
	github.com/ethereum/go-ethereum v1.10.17
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/robfig/cron/v3 v3.0.1
	github.com/stretchr/testify v1.7.1
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v3 v3.0.0-20220512140231-539c8e751b99 // indirect
	gorm.io/driver/postgres v1.3.6
	k8s.io/client-go v0.24.1 // indirect
)

replace (
	github.com/bnb-chain/zkbas-crypto => ../zkbas-crypto
	github.com/bnb-chain/zkbas-eth-rpc => ../zkbas-eth-rpc
)
