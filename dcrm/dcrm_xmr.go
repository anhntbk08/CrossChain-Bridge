package dcrm

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/anyswap/CrossChain-Bridge/common"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/params"
	cryptoEd25519 "github.com/anyswap/CrossChain-Bridge/tools/crypto_ed25519"
	"github.com/anyswap/CrossChain-Bridge/tools/keystore"
	"github.com/anyswap/CrossChain-Bridge/tools/rlp"
	"github.com/anyswap/CrossChain-Bridge/types"
)

type XMRDcrm struct {
	DRCM
}

// DoSignOne dcrm sign single msgHash with context msgContext
func (xmr *XMRDcrm) DoSignOne(signPubkey, msgHash, msgContext string) (keyID string, rsvs []string, err error) {
	return xmr.DoSign(signPubkey, []string{msgHash}, []string{msgContext})
}

// DoSign dcrm sign msgHash with context msgContext
func (xmr *XMRDcrm) DoSign(signPubkey string, msgHash, msgContext []string) (keyID string, rsvs []string, err error) {
	if !params.IsDcrmEnabled() {
		return "", nil, errSignIsDisabled
	}
	log.Debug("dcrm DoSign", "msgHash", msgHash, "msgContext", msgContext)
	if signPubkey == "" {
		return "", nil, errSignWithoutPublickey
	}
	for i := 0; i < retrySignLoop; i++ {
		for _, dcrmNode := range allInitiatorNodes {
			if err = pingDcrmNode(dcrmNode); err != nil {
				continue
			}
			signGroupsCount := int64(len(dcrmNode.signGroups))
			// randomly pick first subgroup to sign
			randIndex, _ := rand.Int(rand.Reader, big.NewInt(signGroupsCount))
			startIndex := randIndex.Int64()
			i := startIndex
			for {
				keyID, rsvs, err = doSignImpl(dcrmNode, i, signPubkey, msgHash, msgContext)
				if err == nil {
					return keyID, rsvs, nil
				}
				i = (i + 1) % signGroupsCount
				if i == startIndex {
					break
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
	log.Warn("dcrm DoSign failed", "msgHash", msgHash, "msgContext", msgContext, "err", err)
	return "", nil, errDoSignFailed
}

func (xmr *XMRDcrm) doSignImpl(dcrmNode *NodeInfo, signGroupIndex int64, signPubkey string, msgHash, msgContext []string) (keyID string, rsvs []string, err error) {
	nonce, err := GetSignNonce(dcrmNode.dcrmUser.String(), dcrmNode.dcrmRPCAddress)
	if err != nil {
		return "", nil, err
	}
	txdata := SignData{
		TxType:     "SIGN",
		PubKey:     signPubkey,
		MsgHash:    msgHash,
		MsgContext: msgContext,
		Keytype:    "ECDSA",
		GroupID:    dcrmNode.signGroups[signGroupIndex],
		ThresHold:  dcrmThreshold,
		Mode:       dcrmMode,
		TimeStamp:  common.NowMilliStr(),
	}
	payload, _ := json.Marshal(txdata)
	rawTX, err := BuildDcrmRawTx(nonce, payload, dcrmNode.keyWrapper)
	if err != nil {
		return "", nil, err
	}

	rpcAddr := dcrmNode.dcrmRPCAddress
	keyID, err = Sign(rawTX, rpcAddr)
	if err != nil {
		return "", nil, err
	}

	rsvs, err = getSignResult(keyID, rpcAddr)
	if err != nil {
		return "", nil, err
	}
	return keyID, rsvs, nil
}

func (xmr *XMRDcrm) getSignResult(keyID, rpcAddr string) (rsvs []string, err error) {
	log.Info("start get sign status", "keyID", keyID)
	var signStatus *SignStatus
	i := 0
	signTimer := time.NewTimer(signTimeout)
	defer signTimer.Stop()
LOOP_GET_SIGN_STATUS:
	for {
		i++
		select {
		case <-signTimer.C:
			if err == nil {
				err = errSignTimerTimeout
			}
			break LOOP_GET_SIGN_STATUS
		default:
			signStatus, err = GetSignStatus(keyID, rpcAddr)
			if err == nil {
				rsvs = signStatus.Rsv
				break LOOP_GET_SIGN_STATUS
			}
			switch {
			case errors.Is(err, ErrGetSignStatusFailed),
				errors.Is(err, ErrGetSignStatusTimeout):
				break LOOP_GET_SIGN_STATUS
			}
		}
		time.Sleep(1 * time.Second)
	}
	if len(rsvs) == 0 || err != nil {
		log.Info("get sign status failed", "keyID", keyID, "retryCount", i, "err", err)
		return nil, errGetSignResultFailed
	}
	log.Info("get sign status success", "keyID", keyID, "retryCount", i)
	return rsvs, nil
}

// BuildDcrmRawTx build dcrm raw tx
func (xmr *XMRDcrm) BuildDcrmRawTx(nonce uint64, payload []byte, keyWrapper *keystore.Key) (string, error) {
	tx := types.NewTransaction(
		nonce,             // nonce
		dcrmToAddr,        // to address
		big.NewInt(0),     // value
		100000,            // gasLimit
		big.NewInt(80000), // gasPrice
		payload,           // data
	)
	signature, err := cryptoEd25519.Sign(dcrmSigner.Hash(tx).Bytes(), keyWrapper.PrivateKey)
	if err != nil {
		return "", err
	}
	sigTx, err := tx.WithSignature(dcrmSigner, signature)
	if err != nil {
		return "", err
	}
	txdata, err := rlp.EncodeToBytes(sigTx)
	if err != nil {
		return "", err
	}
	rawTX := common.ToHex(txdata)
	return rawTX, nil
}
