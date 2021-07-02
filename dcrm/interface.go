package dcrm

import "github.com/anyswap/CrossChain-Bridge/tools/keystore"

type DRCM interface {
	DoSignOne(signPubkey, msgHash, msgContext string) (keyID string, rsvs []string, err error)
	DoSign(signPubkey string, msgHash, msgContext []string) (keyID string, rsvs []string, err error)
	doSignImpl(dcrmNode *NodeInfo, signGroupIndex int64, signPubkey string, msgHash, msgContext []string) (keyID string, rsvs []string, err error)
	getSignResult(keyID, rpcAddr string) (rsvs []string, err error)
	BuildDcrmRawTx(nonce uint64, payload []byte, keyWrapper *keystore.Key) (string, error)
}
