package xmr

import (
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/tokens"
)

// Init init xmr extra
func Init(btcExtra *tokens.BtcExtraConfig) {
	// if xmr.BridgeInstance == nil {
	// 	return
	// }

	// if btcExtra == nil {
	// 	log.Fatal("xmr bridge must config 'BtcExtra'")
	// }

	initFromPublicKey()
}

func initFromPublicKey() {
	if len(tokens.GetTokenPairsConfig()) != 1 {
		log.Fatalf("xmr bridge does not support multiple tokens")
	}

	log.Println("initFromPublicKey PairID ", PairID)
	_, exist := tokens.GetTokenPairsConfig()[PairID]
	if !exist {
		log.Fatalf("xmr bridge must have pairID %v", PairID)
	}

}
