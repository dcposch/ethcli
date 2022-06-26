package action

import "github.com/ethereum/go-ethereum/common"

type TabState struct {
	EnteredAddr  string
	ContractAddr *common.Address
	ErrorText    string
}
