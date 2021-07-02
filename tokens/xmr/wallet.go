package xmr

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"time"

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

func (w *Wallet) Balance() error {
	cmd := `
		monero-wallet-cli --stagenet --trusted-daemon --wallet-file wallet_stagnet1  --password '' >> 
		balances
	`
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)

	defer cancel()

	return exec.CommandContext(ctx, cmd, "5").Run()
}

func (w *Wallet) Addresses() error {
	cmd := `
		monero-wallet-cli --stagenet --trusted-daemon --wallet-file wallet_stagnet1  --password '' >> 
		addresses
	`
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)

	defer cancel()

	return exec.CommandContext(ctx, cmd, "5").Run()
}

func (w *Wallet) ExportMultisigInfo() (string, error) {
	resp, err := client.ExportMultisigInfo()

	return resp.Info, err
}

func (w *Wallet) ImportMultisigInfo() error {
	return errors.New("Not implemented yet")
}

func (w *Wallet) Transfer() error {
	return errors.New("Not implemented yet")
}

func (w *Wallet) SignTransaction() error {
	return errors.New("Not implemented yet")
}
