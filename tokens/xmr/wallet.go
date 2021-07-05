package xmr

import (
	"encoding/json"

	"github.com/anyswap/CrossChain-Bridge/tokens"
	monero "github.com/monero-ecosystem/go-monero-rpc-client/wallet"
)

type WalletInterface interface {
	Refresh() error
	ExportMultisigInfo() (string, error)
}

// XMR Wallet
type Wallet struct {
	tokens.WalletInterface
	// MultisigConfig *params.XMRConfig
}

// just for poc, local wallet rpc node
var client monero.Client = monero.New(monero.Config{
	Address: "http://127.0.0.1:18082/json_rpc",
})

func (w *Wallet) Refresh() error {
	// check wallet balance
	resp, err := client.GetBalance(&monero.RequestGetBalance{AccountIndex: 0})

	if err != nil {
		return err
	}
	_, err = json.MarshalIndent(resp, "", "\t")

	return err
}

func (w *Wallet) Balance() (*monero.ResponseGetBalance, error) {
	return client.GetBalance(&monero.RequestGetBalance{AccountIndex: 0})
}
func (w *Wallet) ExportMultisigInfo() (string, error) {
	resp, err := client.ExportMultisigInfo()

	return resp.Info, err
}

func (w *Wallet) ImportMultisigInfo(info string) (*monero.ResponseImportMultisigInfo, error) {
	return client.ImportMultisigInfo(
		&monero.RequestImportMultisigInfo{
			Info: []string{info},
		},
	)
}

func (w *Wallet) Transfer(target string, amount uint64) (*monero.ResponseTransfer, error) {
	return client.Transfer(
		&monero.RequestTransfer{
			Destinations: []*monero.Destination{
				{
					Amount:  amount,
					Address: target,
				},
			},
		},
	)
}

func (w *Wallet) SignTransaction(multisigTxSet string) (*monero.ResponseSignMultisig, error) {
	return client.SignMultisig(
		&monero.RequestSignMultisig{
			TxDataHex: multisigTxSet,
		},
	)
}
