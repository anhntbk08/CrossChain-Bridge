package xmr

import (
	"errors"
	"math/big"

	"github.com/anyswap/CrossChain-Bridge/tokens"
)

// BuildRawTransaction build raw tx
func (b *Bridge) BuildRawTransaction(args *tokens.BuildTxArgs) (rawTx interface{}, err error) {
	return nil, errors.New("not implemented yet")
}

// BuildTransaction build tx
func (b *Bridge) BuildTransaction(from string, receivers []string, amounts []int64) (rawTx interface{}, err error) {
	return nil, errors.New("not implemented yet")
}

func (b *Bridge) getTxOutputs(to string, amount *big.Int, memo string) (err error) {
	return errors.New("not implemented yet")
}

type insufficientFundsError struct{}

func (insufficientFundsError) InputSourceError() {}
func (insufficientFundsError) Error() string {
	return "insufficient funds available to construct transaction"
}
