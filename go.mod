module github.com/sisu-network/sisu

go 1.15

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/cosmos/cosmos-sdk v0.42.1
	github.com/ethereum/go-ethereum v1.10.12
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gogo/protobuf v1.3.3
	github.com/golang/mock v1.4.4
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/hashicorp/go-retryablehttp v0.7.0
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/joho/godotenv v1.3.0
	github.com/libp2p/go-libp2p-core v0.8.0
	github.com/sisu-network/dcore v0.1.11
	github.com/sisu-network/deyes v0.1.2
	github.com/sisu-network/dheart v0.1.5-alpha1
	github.com/sisu-network/lib v0.0.1-alpha9
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.13
	github.com/tendermint/tm-db v0.6.4
	github.com/tyler-smith/go-bip39 v1.0.2
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	google.golang.org/genproto v0.0.0-20210303154014-9728d6b83eeb // indirect
	google.golang.org/grpc v1.37.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
