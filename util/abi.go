package util

import (
	"math/big"

	gethmath "github.com/ethereum/go-ethereum/common/math"
)

func EncodeUint(v *big.Int) []byte {
	return gethmath.U256Bytes(v)
}

func DecodeUint(bytes []byte) *big.Int {
	return big.NewInt(0).SetBytes(bytes)
}
