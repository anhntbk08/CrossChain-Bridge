package dcrm

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os/exec"

	"github.com/anyswap/CrossChain-Bridge/tools/keystore"
)

type XMRWallet struct {
	CLI_PATH        string
	MultisigAddress string
	balance         *big.Int
	key             *keystore.Key
	PublicViewKey   ecdsa.PublicKey
	PublicSpentKey  ecdsa.PublicKey
	PrivateViewKey  ecdsa.PrivateKey
	PrivateSpentKey ecdsa.PrivateKey
}

func NewXMRWallet(cliPath string, key *keystore.Key) *XMRWallet {
	xmr := &XMRWallet{
		CLI_PATH: cliPath,
		key:      key,
	}

	xmr.init()

	return xmr
}

func (wallet *XMRWallet) init() error {
	err := exec.Command(wallet.CLI_PATH).Run()
	if err != nil {
		fmt.Println(err.Error())
	}

	return err
}

func (wallet *XMRWallet) GetMultisigAddress() error {
	err := exec.Command(wallet.CLI_PATH).Run()
	if err != nil {
		fmt.Println(err.Error())
	}

	return err
}

func (wallet *XMRWallet) GetPublicKeys() error {
	return errors.New("Not implemented yet")
}

func (wallet *XMRWallet) GetPrivateKeys() error {
	return errors.New("Not implemented yet")
}

func (wallet *XMRWallet) GetBalance() error {
	return errors.New("Not implemented yet")
}

func (wallet *XMRWallet) ExportKeyImages() error {
	return errors.New("Not implemented yet")
}
