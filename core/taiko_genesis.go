package core

import (
	"github.com/ethereum/go-ethereum/common"
	mxcGenesis "github.com/ethereum/go-ethereum/core/mxc_genesis"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

// MXCGenesisBlock returns the MXC network genesis block configs.
func MXCGenesisBlock(networkID uint64) *Genesis {
	chainConfig := params.MXCChainConfig

	var allocJSON []byte
	switch networkID {
	case params.MxcWannseeNetworkID.Uint64():
		chainConfig.ChainID = params.MxcWannseeNetworkID
		allocJSON = mxcGenesis.MxcWannseeGenesisAllocJSON
	default:
		chainConfig.ChainID = params.MXCMainnetNetworkID
		allocJSON = mxcGenesis.MXCMainNetGenesisAllocJSON
	}

	var alloc GenesisAlloc
	if err := alloc.UnmarshalJSON(allocJSON); err != nil {
		log.Crit("unmarshal alloc json error", "error", err)
	}

	return &Genesis{
		Config:     chainConfig,
		ExtraData:  []byte{},
		GasLimit:   uint64(5000000),
		Difficulty: common.Big0,
		Alloc:      alloc,
	}
}
