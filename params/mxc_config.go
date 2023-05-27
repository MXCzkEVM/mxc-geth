package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func u64(val uint64) *uint64 { return &val }

// Network IDs
var (
	MxcMainnetNetworkID = big.NewInt(516)
	MxcWannseeNetworkID = big.NewInt(5167003)
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
	Treasury:                      common.HexToAddress("0xdf09A0afD09a63fb04ab3573922437e1e637dE8b"),
}
