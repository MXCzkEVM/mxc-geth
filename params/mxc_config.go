package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Network IDs
var (
	MXCMainnetNetworkID = big.NewInt(5167)
	MxcWannseeNetworkID = big.NewInt(5167003)
)

var MXCChainConfig = &ChainConfig{
	ChainID:                       MXCMainnetNetworkID, // Use mainnet network ID by default.
	HomesteadBlock:                common.Big0,
	EIP150Block:                   common.Big0,
	EIP155Block:                   common.Big0,
	EIP158Block:                   common.Big0,
	ByzantiumBlock:                common.Big0,
	ConstantinopleBlock:           common.Big0,
	PetersburgBlock:               common.Big0,
	IstanbulBlock:                 common.Big0,
	BerlinBlock:                   common.Big0,
	LondonBlock:                   nil,
	MergeNetsplitBlock:            nil,
	TerminalTotalDifficulty:       common.Big0,
	TerminalTotalDifficultyPassed: true,
	MXC:                           true,
}
