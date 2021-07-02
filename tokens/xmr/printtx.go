package xmr

import (
	"encoding/json"
	"math/big"

	"github.com/anyswap/CrossChain-Bridge/common/hexutil"
)

// MarshalToJSON marshal to json
func MarshalToJSON(obj interface{}, pretty bool) string {
	var jsdata []byte
	if pretty {
		jsdata, _ = json.MarshalIndent(obj, "", "  ")
	} else {
		jsdata, _ = json.Marshal(obj)
	}
	return string(jsdata)
}

// AuthoredTxToString AuthoredTx to string
func AuthoredTxToString(authtx interface{}, pretty bool) string {
	return ""
}

// EncAuthoredTx stuct
type EncAuthoredTx struct {
	Tx          *EncMsgTx
	TotalInput  *big.Int
	ChangeIndex int
}

// EncMsgTx struct
type EncMsgTx struct {
	Txid     string
	Version  int32
	TxIn     []*EncTxIn
	TxOut    []*EncTxOut
	LockTime uint32
}

// EncTxOut struct
type EncTxOut struct {
	PkScript string
	Value    int64
}

// EncOutPoint struct
type EncOutPoint struct {
	Hash  string
	Index uint32
}

// EncTxIn struct
type EncTxIn struct {
	PreviousOutPoint EncOutPoint
	SignatureScript  string
	Witness          []hexutil.Bytes
	Sequence         uint32
	Value            *big.Int
}
