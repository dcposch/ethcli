package util

import (
	"math/big"
	"strings"
)

func ToFixedPrecision(val *big.Int, dec int) string {
	strV := strings.Repeat("0", dec) + val.String()
	decIx := len(strV) - dec
	ret := strings.TrimLeft(strV[:decIx]+"."+strV[decIx:], "0")
	if strings.HasPrefix(ret, ".") {
		ret = "0" + ret
	}
	return ret
}
