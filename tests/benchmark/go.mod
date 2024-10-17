module github.com/0glabs/0g-chain/tests/benchmark

go 1.23.1

require (
	cosmossdk.io/errors v1.0.1
	github.com/cosmos/cosmos-sdk v0.47.10
	github.com/ethereum/go-ethereum v1.10.26
	github.com/evmos/ethermint v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/shopspring/decimal v1.4.0
	golang.org/x/crypto v0.24.0
	golang.org/x/exp v0.0.0-20230711153332-06a737ee72cb
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
)

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/cometbft/cometbft v0.37.4 // indirect
	github.com/cosmos/go-bip39 v1.0.0 // indirect
	github.com/cosmos/gogoproto v1.4.10 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.4.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/holiman/uint256 v1.3.1 // indirect
	github.com/linxGnu/grocksdb v1.9.3 // indirect
	github.com/petermattis/goid v0.0.0-20230317030725-371a4b8eda08 // indirect
	github.com/rjeczalik/notify v0.9.1 // indirect
	github.com/sasha-s/go-deadlock v0.3.1 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/tendermint/go-amino v0.16.0 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.23.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240108191215-35c7eff3a6b1 // indirect
	google.golang.org/grpc v1.60.1 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace (
	github.com/cometbft/cometbft => github.com/kava-labs/cometbft v0.37.4-kava.1
	github.com/cometbft/cometbft-db => github.com/kava-labs/cometbft-db v0.9.1-kava.1
	github.com/cosmos/cosmos-sdk => github.com/0glabs/cosmos-sdk v0.47.10-0glabs.3
	github.com/ethereum/go-ethereum => github.com/evmos/go-ethereum v1.10.26-evmos-rc2
	github.com/evmos/ethermint => github.com/0glabs/ethermint v0.21.0-0g.v3.0.3
)
