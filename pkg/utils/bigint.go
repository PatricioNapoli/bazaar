package utils

import (
	"math"
	"math/big"
)

func ReduceBigInt(n *big.Int, decimals int) float64 {
	mul := math.Pow(10.0, -float64(decimals))

	flt := new(big.Float).SetPrec(128).SetInt(n)
	flt = new(big.Float).Mul(flt, new(big.Float).SetFloat64(mul))

	final, _ := flt.Float64()
	return final
}
