package mxc_genesis

import (
	_ "embed"
)

//go:embed mainnet.json
var MainnetGenesisAllocJSON []byte

//go:embed wannsee.json
var MxcWannseeGenesisAllocJSON []byte

//go:embed geneva.json
var MxcGenevaGenesisAllocJSON []byte
