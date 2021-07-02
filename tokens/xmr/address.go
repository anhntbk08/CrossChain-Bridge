package xmr

// DecodeAddress decode address
func (b *Bridge) DecodeAddress(addr string) (address []byte, err error) {
	return
}

// IsValidAddress check address
func (b *Bridge) IsValidAddress(addr string) bool {
	// _, err := b.DecodeAddress(addr)
	return true
	// return err == nil
}
