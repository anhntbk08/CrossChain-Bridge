package xmr

import (
	"errors"

	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/xmr/electrs"
)

// GetTransaction impl
func (b *Bridge) GetTransaction(txHash string) (interface{}, error) {
	return b.GetTransactionByHash(txHash)
}

// GetTransactionStatus impl
func (b *Bridge) GetTransactionStatus(txHash string) *tokens.TxStatus {
	txStatus := &tokens.TxStatus{}
	electStatus, err := b.GetElectTransactionStatus(txHash)
	if err != nil {
		log.Trace(b.ChainConfig.BlockChain+" Bridge::GetElectTransactionStatus fail", "tx", txHash, "err", err)
		return txStatus
	}
	if !*electStatus.Confirmed {
		return txStatus
	}
	if electStatus.BlockHash != nil {
		txStatus.BlockHash = *electStatus.BlockHash
	}
	if electStatus.BlockTime != nil {
		txStatus.BlockTime = *electStatus.BlockTime
	}
	if electStatus.BlockHeight != nil {
		txStatus.BlockHeight = *electStatus.BlockHeight
		latest, err := b.GetLatestBlockNumber()
		if err != nil {
			log.Debug(b.ChainConfig.BlockChain+" Bridge::GetLatestBlockNumber fail", "err", err)
			return txStatus
		}
		if latest > txStatus.BlockHeight {
			txStatus.Confirmations = latest - txStatus.BlockHeight
		}
	}
	return txStatus
}

// VerifyMsgHash verify msg hash
func (b *Bridge) VerifyMsgHash(rawTx interface{}, msgHash []string) (err error) {
	return errors.New("Not implemented yet")
}

// VerifyTransaction impl
func (b *Bridge) VerifyTransaction(pairID, txHash string, allowUnstable bool) (*tokens.TxSwapInfo, error) {
	if !b.IsSrc {
		return nil, tokens.ErrBridgeDestinationNotSupported
	}
	return b.verifySwapinTx(pairID, txHash, allowUnstable)
}

func (b *Bridge) verifySwapinTx(pairID, txHash string, allowUnstable bool) (*tokens.TxSwapInfo, error) {
	return nil, errors.New("not implemented yet")
}

func (b *Bridge) checkSwapinInfo(swapInfo *tokens.TxSwapInfo) error {
	if swapInfo.From == swapInfo.To {
		return tokens.ErrTxWithWrongSender
	}
	// if !tokens.CheckSwapValue(swapInfo.PairID, swapInfo.Value, b.IsSrc) {
	// 	return tokens.ErrTxWithWrongValue
	// }
	if !tokens.DstBridge.IsValidAddress(swapInfo.Bind) {
		log.Debug("wrong bind address in swapin", "bind", swapInfo.Bind)
		return tokens.ErrTxWithWrongMemo
	}
	return nil
}

func (b *Bridge) checkStable(txHash string) bool {
	txStatus := b.GetTransactionStatus(txHash)
	confirmations := *b.GetChainConfig().Confirmations
	return txStatus.BlockHeight > 0 && txStatus.Confirmations >= confirmations
}

// GetReceivedValue get received value
func (b *Bridge) GetReceivedValue(vout []*electrs.ElectTxOut, receiver, pubkeyType string) (value uint64, memoScript string, rightReceiver bool) {
	return 0, "", false
}

// return priorityAddress if has it in Vin
// return the first address in Vin if has no priorityAddress
func getTxFrom(vin []*electrs.ElectTxin, priorityAddress string) string {
	from := ""
	for _, input := range vin {
		if input != nil &&
			input.Prevout != nil &&
			input.Prevout.ScriptpubkeyAddress != nil {
			if *input.Prevout.ScriptpubkeyAddress == priorityAddress {
				return priorityAddress
			}
			if from == "" {
				from = *input.Prevout.ScriptpubkeyAddress
			}
		}
	}
	return from
}
