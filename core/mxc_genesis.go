package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	mxcGenesis "github.com/ethereum/go-ethereum/core/mxc_genesis"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

// MxcGenesisBlock returns the Mxc network genesis block configs.
func MxcGenesisBlock(networkID uint64) *Genesis {
	chainConfig := params.MxcChainConfig

	var allocJSON []byte
	switch networkID {
	case params.MxcWannseeNetworkID.Uint64():
		chainConfig.ChainID = params.MxcWannseeNetworkID
		allocJSON = mxcGenesis.MxcWannseeGenesisAllocJSON
	case params.MxcGenevaNetworkID.Uint64():
		chainConfig.ChainID = params.MxcGenevaNetworkID
		allocJSON = mxcGenesis.MxcGenevaGenesisAllocJSON
	default:
		chainConfig.ChainID = params.MxcMainnetNetworkID
		allocJSON = mxcGenesis.MainnetGenesisAllocJSON
	}

	var alloc GenesisAlloc
	if err := alloc.UnmarshalJSON(allocJSON); err != nil {
		log.Crit("unmarshal alloc json error", "error", err)
	}

	return &Genesis{
		Config:     chainConfig,
		ExtraData:  []byte{},
		GasLimit:   uint64(6000000),
		Difficulty: common.Big0,
		Alloc:      alloc,
		GasUsed:    0,
		BaseFee:    new(big.Int).SetUint64(10000000),
	}
}
