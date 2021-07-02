package xmr

import (
	"crypto/ecdsa"
	"errors"

	"github.com/anyswap/CrossChain-Bridge/tokens"
)

// DcrmSignTransaction dcrm sign raw tx
func (b *Bridge) DcrmSignTransaction(rawTx interface{}, args *tokens.BuildTxArgs) (signedTx interface{}, txHash string, err error) {
	return nil, "", errors.New("Not implemented yet")
}

// MakeSignedTransaction make signed tx
func (b *Bridge) MakeSignedTransaction(authoredTx interface{}, msgHash, rsv []string, sigScripts [][]byte, cPkData []byte) (signedTx interface{}, txHash string, err error) {
	return nil, "", errors.New("Not implemented yet")
}

// GetCompressedPublicKey get compressed public key
func (b *Bridge) GetCompressedPublicKey(fromPublicKey string, needVerify bool) (cPkData []byte, err error) {
	return []byte{}, errors.New("Not implemented yet")
}

// DcrmSignMsgHash dcrm sign msg hash
func (b *Bridge) DcrmSignMsgHash(msgHash []string, args *tokens.BuildTxArgs) (rsv []string, err error) {
	return []string{}, errors.New("Not implemented yet")
}

// SignTransaction sign tx with pairID
func (b *Bridge) SignTransaction(rawTx interface{}, pairID string) (signedTx interface{}, txHash string, err error) {
	privKey := b.GetTokenConfig(pairID).GetDcrmAddressPrivateKey()
	return b.SignTransactionWithPrivateKey(rawTx, privKey)
}

// SignTransactionWithPrivateKey sign tx with ECDSA private key
func (b *Bridge) SignTransactionWithPrivateKey(rawTx interface{}, privKey *ecdsa.PrivateKey) (signTx interface{}, txHash string, err error) {
	return nil, "", errors.New("Not implemented yet")
}
