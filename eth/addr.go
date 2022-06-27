package eth

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
)

// An address plus ENS name
type NamedAddr struct {
	Addr common.Address
	Name string
	Err  string
}

// Displays an address, eg "0x123456..." or "vitalik.eth" or "⚠️ invalid.eth"
func (a *NamedAddr) Disp() string {
	ret := ""
	if a.Err != "" {
		ret += "⚠️ "
	}
	if a.Name != "" {
		ret += a.Name
	} else {
		hex := a.Addr.Hex()
		ret += hex[0:8] + "…" + hex[36:]
	}
	return ret
}

var ZeroAddr = common.Address{}

func IsZeroAddr(a common.Address) bool {
	return bytes.Equal(a[:], ZeroAddr[:])
}
