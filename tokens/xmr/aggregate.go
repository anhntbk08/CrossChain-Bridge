package xmr

import (
	"errors"

	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/xmr/electrs"
)

const (
	redeemAggregateP2SHInputSize = 198
)

// ShouldAggregate should aggregate
func (b *Bridge) ShouldAggregate(aggUtxoCount int, aggSumVal uint64) bool {
	return false
}

// AggregateUtxos aggregate uxtos
func (b *Bridge) AggregateUtxos(addrs []string, utxos []*electrs.ElectUtxo) (string, error) {
	return "", nil
}

// VerifyAggregateMsgHash verify aggregate msgHash
func (b *Bridge) VerifyAggregateMsgHash(msgHash []string, args *tokens.BuildTxArgs) error {
	return errors.New("Not implemented yet")
}
