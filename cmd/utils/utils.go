package utils

import (
	"github.com/tomochain/go-tomochain/dashboard"
	"github.com/tomochain/go-tomochain/eth"
	"github.com/tomochain/go-tomochain/eth/downloader"
	"github.com/tomochain/go-tomochain/ethstats"
	"github.com/tomochain/go-tomochain/les"
	"github.com/tomochain/go-tomochain/node"
	"github.com/tomochain/go-tomochain/tomox"
	whisper "github.com/tomochain/go-tomochain/whisper/whisperv6"
)

// RegisterEthService adds an Ethereum client to the stack.
func RegisterEthService(stack *node.Node, cfg *eth.Config) {
	var err error
	if cfg.SyncMode == downloader.LightSync {
		err = stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
			return les.New(ctx, cfg)
		})
	} else {
		err = stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
			var tomoXServ *tomox.TomoX
			ctx.Service(&tomoXServ)

			fullNode, err := eth.New(ctx, cfg, tomoXServ)
			if fullNode != nil && cfg.LightServ > 0 {
				ls, _ := les.NewLesServer(fullNode, cfg)
				fullNode.AddLesServer(ls)
			}
			return fullNode, err
		})
	}
	if err != nil {
		Fatalf("Failed to register the Ethereum service: %v", err)
	}
}

// RegisterDashboardService adds a dashboard to the stack.
func RegisterDashboardService(stack *node.Node, cfg *dashboard.Config, commit string) {
	stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		return dashboard.New(cfg, commit)
	})
}

// RegisterShhService configures Whisper and adds it to the given node.
func RegisterShhService(stack *node.Node, cfg *whisper.Config) {
	if err := stack.Register(func(n *node.ServiceContext) (node.Service, error) {
		return whisper.New(cfg), nil
	}); err != nil {
		Fatalf("Failed to register the Whisper service: %v", err)
	}
}

// RegisterEthStatsService configures the Ethereum Stats daemon and adds it to
// th egiven node.
func RegisterEthStatsService(stack *node.Node, url string) {
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		// Retrieve both eth and les services
		var ethServ *eth.Ethereum
		ctx.Service(&ethServ)

		var lesServ *les.LightEthereum
		ctx.Service(&lesServ)

		return ethstats.New(url, ethServ, lesServ)
	}); err != nil {
		Fatalf("Failed to register the Ethereum Stats service: %v", err)
	}
}

func RegisterTomoXService(stack *node.Node, cfg *tomox.Config) {
	if err := stack.Register(func(n *node.ServiceContext) (node.Service, error) {
		return tomox.New(cfg), nil
	}); err != nil {
		Fatalf("Failed to register the TomoX service: %v", err)
	}
}
