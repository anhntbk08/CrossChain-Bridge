package xmr

import (
	"fmt"
	"strings"
)

// Inheritable interface
type Inheritable interface {
	GetChainParams() string
}

// GetChainParams get chain config (net params)
func (b *Bridge) GetChainParams() string {
	networkID := strings.ToLower(b.ChainConfig.NetID)
	switch networkID {
	case netMainnet:
		return netMainnet
	default:
		return netCustom
	}
}

// GetPayToAddrScript get pay to address script
func (b *Bridge) GetPayToAddrScript(address string) ([]byte, error) {
	_, err := b.DecodeAddress(address)
	return nil, fmt.Errorf("decode xmr address '%v' failed. %w", address, err)
}
