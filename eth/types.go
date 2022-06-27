package eth

import (
	"encoding/binary"
	"encoding/json"

	"dcposch.eth/cli/util"
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

// const (
// 	TypeInAmountJson = `{"components":[{"name":"label","type":"string"},{"name":"decimals","type":"uint64"}],"name":"data","type":"tuple"}`
// )

type ElemTs []abi.ArgumentMarshaling

var (
	ElemsInAmount = ElemTs{{Name: "label", Type: "string"}, {Name: "decimals", Type: "uint64"}}
	ElemsButton   = ElemTs{{Name: "text", Type: "string"}}
)

func ParseTuple(bytes []byte, elems ElemTs, ret interface{}) {
	typ, err := abi.NewType("tuple", "", elems)
	util.Must(err)

	// TODO: replace the geth abi-parsing.
	// It is both hard to use and incredibly ugly. And reflective, likely slow.
	args := abi.Arguments{{Type: typ}}
	res, err := args.UnpackValues(bytes)
	util.Must(err)

	// Ineffective: util.Must(args.Copy(&wrap, res))
	// As a workaround, round-trip through JSON instead.
	wrap := []interface{}{ret}
	js, err := json.Marshal(res)
	util.Must(err)
	util.Must(json.Unmarshal(js, &wrap))
}

type VdomElem struct {
	TypeHash uint64
	// Text for a text field, options for a dropdown, etc.
	Data []byte
	// Data parsed into a struct. See DataDropdown, etc.
	DataStruct interface{}
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

type DataBtnAction struct {
	// 0 = first button, etc.
	ButtonId uint64
	// Value of each input.
	Inputs [][]byte
}
