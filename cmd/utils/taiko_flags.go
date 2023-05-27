package utils

import (
	"os"

	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/urfave/cli/v2"
)

var (
	MxcFlag = cli.BoolFlag{
		Name:  "mxc",
		Usage: "mxc network",
	}
)

// RegisterMxcAPIs initializes and registers the mxc RPC APIs.
func RegisterMxcAPIs(stack *node.Node, cfg *ethconfig.Config, backend *eth.Ethereum) {
	if os.Getenv("MXC_TEST") != "" {
		return
	}
	// Add methods under "mxc_" RPC namespace to the available APIs list
	stack.RegisterAPIs([]rpc.API{
		{
			Namespace: "mxc",
			Version:   params.VersionWithMeta,
			Service:   eth.NewMxcAPIBackend(backend),
			Public:    true,
		},
	})
}
