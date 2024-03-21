package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func u64(val uint64) *uint64 { return &val }

// Network IDs
var (
	MxcMainnetNetworkID = big.NewInt(18686)
	MxcWannseeNetworkID = big.NewInt(5167003)
	MxcGenevaNetworkID  = big.NewInt(5167004)
)

var MxcChainConfig = &ChainConfig{
	ChainID:                       MxcMainnetNetworkID, // Use mainnet network ID by default.
	HomesteadBlock:                common.Big0,
	EIP150Block:                   common.Big0,
	EIP155Block:                   common.Big0,
	EIP158Block:                   common.Big0,
	ByzantiumBlock:                common.Big0,
	ConstantinopleBlock:           common.Big0,
	PetersburgBlock:               common.Big0,
	IstanbulBlock:                 common.Big0,
	BerlinBlock:                   common.Big0,
	LondonBlock:                   common.Big0,
	ShanghaiTime:                  u64(0),
	MergeNetsplitBlock:            nil,
	TerminalTotalDifficulty:       common.Big0,
	TerminalTotalDifficultyPassed: true,
	Mxc:                           true,
	Treasury:                      common.HexToAddress("0x2000777700000000000000000000000000000001"),
}
