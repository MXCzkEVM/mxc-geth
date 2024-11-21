package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Network IDs
var (
	MoonchainMainnetNetworkID = big.NewInt(18686)
	MoonchainGenevaNetworkID  = big.NewInt(5167004)
)

var networkIDToChainConfigByMoonchain = map[*big.Int]*ChainConfig{
	MoonchainMainnetNetworkID:  MoonchainChainConfig,
	MoonchainGenevaNetworkID:   MoonchainChainConfig,
	MainnetChainConfig.ChainID: MainnetChainConfig,
	SepoliaChainConfig.ChainID: SepoliaChainConfig,
	TestChainConfig.ChainID:    TestChainConfig,
	NonActivatedConfig.ChainID: NonActivatedConfig,
}

func NetworkIDToChainConfigOrDefaultByMoonchain(networkID *big.Int) *ChainConfig {
	if config, ok := networkIDToChainConfigByMoonchain[networkID]; ok {
		return config
	}

	return AllEthashProtocolChanges
}

var MoonchainChainConfig = &ChainConfig{
	ChainID:                       MoonchainGenevaNetworkID, // Use Mxc Geneva network ID by default.
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
