package eth

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	TypeText       = binary.BigEndian.Uint64(crypto.Keccak256([]byte("text"))[24:])
	TypeInAmount   = binary.BigEndian.Uint64(crypto.Keccak256([]byte("amount"))[24:])
	TypeInDropdown = binary.BigEndian.Uint64(crypto.Keccak256([]byte("dropdown"))[24:])
	TypeInTextbox  = binary.BigEndian.Uint64(crypto.Keccak256([]byte("textbox"))[24:])
	TypeButton     = binary.BigEndian.Uint64(crypto.Keccak256([]byte("button"))[24:])
)

type VdomElem struct {
	TypeHash uint64
	// Text for a text field, options for a dropdown, etc.
	Data []byte
}

type DataAmount struct {
	Label string
	// Amount input will return fixed-point uint256 to n decimals.
	Decimals uint64
}

type DataDropdown struct {
	Label string
	// Options. User must pick one.
	Options []DataDropOption
}

type DataDropOption struct {
	// Dropdown option ID
	Id uint64
	// Dropdown option display string
	Display string
}

type Action struct {
	// 0 = first button, etc.
	ButtonId uint64
	// Value of each input.
	Inputs [][]byte
}
