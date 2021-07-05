package xmr

import (
	"fmt"
	"strings"
	"time"

	"github.com/anyswap/CrossChain-Bridge/log"

	"github.com/anyswap/CrossChain-Bridge/tokens"
)

const (
	netMainnet = "mainnet"
	Stagnet    = "stagenet"
	netCustom  = "testnet"
)

// PairID unique xmr pair ID
var PairID = "xmr"

// Bridge xmr bridge
type Bridge struct {
	*tokens.CrossChainBridgeBase
	Inherit Inheritable
	// MultisigConfig *params.XMRConfig
	XMRWallet *Wallet
}

// NewCrossChainBridge new xmr bridge
func NewCrossChainBridge(isSrc bool) *Bridge {
	if !isSrc {
		log.Fatalf("xmr::NewCrossChainBridge error %v", tokens.ErrBridgeDestinationNotSupported)
	}
	instance := &Bridge{CrossChainBridgeBase: tokens.NewCrossChainBridgeBase(isSrc)}
	BridgeInstance = instance
	instance.SetInherit(instance)
	instance.XMRWallet = &Wallet{}
	// xmr.BridgeInstance = BridgeInstance
	return instance
}

// SetInherit set inherit
func (b *Bridge) SetInherit(inherit Inheritable) {
	b.Inherit = inherit
}

// SetChainAndGateway set chain and gateway config
func (b *Bridge) SetChainAndGateway(chainCfg *tokens.ChainConfig, gatewayCfg *tokens.GatewayConfig) {
	b.CrossChainBridgeBase.SetChainAndGateway(chainCfg, gatewayCfg)
	b.VerifyChainConfig()
	b.InitLatestBlockNumber()
}

// VerifyChainConfig verify chain config
func (b *Bridge) VerifyChainConfig() {
	chainCfg := b.ChainConfig
	networkID := strings.ToLower(chainCfg.NetID)
	switch networkID {
	case netMainnet, Stagnet:
	case netCustom:
		return
	default:
		log.Fatal("unsupported bitcoin network", "netID", chainCfg.NetID)
	}
}

// VerifyTokenConfig verify token config
func (b *Bridge) VerifyTokenConfig(tokenCfg *tokens.TokenConfig) error {
	return nil
}

// InitLatestBlockNumber init latest block number
func (b *Bridge) InitLatestBlockNumber() {
	var latest uint64
	var err error

	for {
		err = b.XMRWallet.Refresh()
		fmt.Printf("xmr wallet erro ", err)
		if err == nil {
			tokens.SetLatestBlockHeight(latest, b.IsSrc)
			// log.Info("get latst block number succeed.", "number", latest, "BlockChain", chainCfg.BlockChain, "NetID", chainCfg.NetID)
			break
		}
		time.Sleep(30 * time.Second)
	}
}
