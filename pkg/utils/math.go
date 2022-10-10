package utils

import (
	"math"
	"math/big"
)

// ReduceBigInt reduces a big int by decimals amount and returns it in a float64 representation.
func ReduceBigInt(n *big.Int, decimals int) float64 {
	mul := math.Pow(10.0, -float64(decimals))

	flt := new(big.Float).SetPrec(128).SetInt(n)
	flt = new(big.Float).Mul(flt, new(big.Float).SetFloat64(mul))

	final, _ := flt.Float64()
	return final
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
