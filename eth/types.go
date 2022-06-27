package eth

import (
	"encoding/binary"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	TypeText       = binary.BigEndian.Uint64(crypto.Keccak256([]byte("text"))[24:])
	TypeInAmount   = binary.BigEndian.Uint64(crypto.Keccak256([]byte("amount"))[24:])
	TypeInDropdown = binary.BigEndian.Uint64(crypto.Keccak256([]byte("dropdown"))[24:])
	TypeInTextbox  = binary.BigEndian.Uint64(crypto.Keccak256([]byte("textbox"))[24:])
	TypeButton     = binary.BigEndian.Uint64(crypto.Keccak256([]byte("button"))[24:])
)

type ElemTs []abi.ArgumentMarshaling

var (
	PropsText     = ElemTs{{Name: "key", Type: "uint256"}, {Name: "text", Type: "string"}}
	PropsAmount   = ElemTs{{Name: "key", Type: "uint256"}, {Name: "label", Type: "string"}, {Name: "decimals", Type: "uint64"}}
	DropOpt       = ElemTs{{Name: "val", Type: "uint256"}, {Name: "text", Type: "string"}}
	PropsDropdown = ElemTs{{Name: "key", Type: "uint256"}, {Name: "label", Type: "string"}, {Name: "options", Type: "tuple[]", Components: DropOpt}}
	PropsButton   = ElemTs{{Name: "key", Type: "uint256"}, {Name: "text", Type: "string"}}
)

func ParseTuple(bytes []byte, elems ElemTs, ret interface{}) error {
	typ, err := abi.NewType("tuple", "", elems)
	if err != nil {
		return err
	}

	// TODO: replace the geth abi-parsing.
	// It is both hard to use and incredibly ugly. And reflective, likely slow.
	args := abi.Arguments{{Type: typ}}
	res, err := args.UnpackValues(bytes)
	if err != nil {
		return err
	}

	// Ineffective: util.Must(args.Copy(&wrap, res))
	// As a workaround, round-trip through JSON instead.
	wrap := []interface{}{ret}
	js, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return json.Unmarshal(js, &wrap)
}

// Virtual DOM element.
// Loosely inspired by React, but radically simplified to fit EVM constraints.
// The VDOM is a flat list of VElems, not a tree. Styling options are tightly
// constrained. Focus is on functionality.
type VElem struct {
	TypeHash uint64
	// Raw ABI-encoded data. Text for a text field, options for a dropdown, etc.
	Data []byte
	// Parsed data. See ElemText, etc.
	DataElem interface{}
}

type Elem struct {
	Key big.Int
}

type ElemText struct {
	Elem
	Text string
}

type ElemAmount struct {
	Elem
	Label string
	// Amount input will return fixed-point uint256 to n decimals.
	Decimals uint64
}

type ElemDropdown struct {
	Elem
	Label string
	// Options. User must pick one.
	Options []ElemDropOption
}

type ElemDropOption struct {
	// Dropdown option value
	Val big.Int
	// Dropdown option display string
	Display string
}

type ElemButton struct {
	Elem
	// Button label
	Text string
}

type ButtonAction struct {
	// Which button was pressed.
	ButtonKey uint64
	// ABI serialization of each input.
	Inputs [][]byte
}
