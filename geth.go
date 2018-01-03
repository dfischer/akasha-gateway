package main

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/params"
)

// makeGeth assembles a go-ethereum in-process node to interact with the Rinkeby
// blockchain and the Akasha smart contracts.
func makeGeth(datadir string) (*node.Node, error) {
	// Define the basic configurations for the Ethereum node
	config := &node.Config{
		Name:    "geth",
		Version: params.Version,
		DataDir: datadir,
		P2P: p2p.Config{
			MaxPeers: 25,
		},
		NoUSB: true,
	}
	// Inject the Rinkeby bootnodes into the configs
	config.P2P.BootstrapNodes = make([]*discover.Node, 0, len(params.RinkebyBootnodes))
	for _, url := range params.RinkebyBootnodes {
		node, err := discover.ParseNode(url)
		if err != nil {
			return nil, err
		}
		config.P2P.BootstrapNodes = append(config.P2P.BootstrapNodes, node)
	}
	// Start the node and configure a full Ethereum node on it
	stack, err := node.New(config)
	if err != nil {
		return nil, err
	}
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		return eth.New(ctx, &eth.Config{
			Genesis:         core.DefaultRinkebyGenesisBlock(),
			NetworkId:       params.RinkebyChainConfig.ChainId.Uint64(),
			DatabaseHandles: 2048,
			DatabaseCache:   2048,
			TxPool:          core.DefaultTxPoolConfig,
			GPO:             eth.DefaultConfig.GPO,
		})
	}); err != nil {
		return nil, err
	}
	return stack, nil
}