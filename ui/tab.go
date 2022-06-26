package ui

import "github.com/ethereum/go-ethereum/common"

type Tab struct {
	enteredAddr  string
	contractAddr *common.Address
	errorText    string
}
