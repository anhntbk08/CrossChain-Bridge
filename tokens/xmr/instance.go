package xmr

import (
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/xmr/electrs"
)

// BridgeInstance xmr bridge instance
var BridgeInstance BridgeInterface

// BridgeInterface xmr bridge interface
type BridgeInterface interface {
	tokens.CrossChainBridge

	GetCompressedPublicKey(fromPublicKey string, needVerify bool) (cPkData []byte, err error)
	VerifyAggregateMsgHash(msgHash []string, args *tokens.BuildTxArgs) error
	AggregateUtxos(addrs []string, utxos []*electrs.ElectUtxo) (string, error)
	FindUtxos(addr string) ([]*electrs.ElectUtxo, error)
	StartSwapHistoryScanJob()
	ShouldAggregate(aggUtxoCount int, aggSumVal uint64) bool
}
