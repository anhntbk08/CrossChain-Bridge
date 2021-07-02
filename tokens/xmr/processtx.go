package xmr

import (
	"errors"

	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/tools"
	"github.com/anyswap/CrossChain-Bridge/tokens/xmr/electrs"
)

func (b *Bridge) processTransaction(txid string) {
	var tx *electrs.ElectTx
	var err error
	for i := 0; i < 2; i++ {
		tx, err = b.GetTransactionByHash(txid)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Debug("[processTransaction] "+b.ChainConfig.BlockChain+" Bridge::GetTransaction fail", "tx", txid, "err", err)
		return
	}
	b.processTransactionImpl(tx)
}

func (b *Bridge) processTransactionImpl(tx *electrs.ElectTx) {
}

func (b *Bridge) processSwapin(txid string) {
	if tools.IsSwapExist(txid, PairID, "", true) {
		return
	}
	swapInfo, err := b.verifySwapinTx(PairID, txid, true)
	tools.RegisterSwapin(txid, []*tokens.TxSwapInfo{swapInfo}, []error{err})
}

// CheckSwapinTxType check swapin type
func (b *Bridge) CheckSwapinTxType(tx *electrs.ElectTx) (p2shBindAddrs []string, err error) {
	return nil, errors.New("Not impelmented yet")
}
