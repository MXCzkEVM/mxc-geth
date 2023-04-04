package mxc_genesis

import (
	_ "embed"
)

//go:embed mainnet.json
var MXCMainNetGenesisAllocJSON []byte

//go:embed mxc-wannsee.json
var MxcWannseeGenesisAllocJSON []byte
