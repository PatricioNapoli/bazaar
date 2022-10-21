package utils

import (
	"math"
	"math/big"
)

// ExtendBigInt extends a big int by decimals amount and returns it
func ExtendBigInt(n *big.Int, precision int) *big.Int {
	mul := math.Pow(10.0, float64(precision))
	mulBig := new(big.Int).SetUint64(uint64(mul))

	nBig := new(big.Int).Mul(n, mulBig)

	return nBig
}

// ReduceBigInt reduces a big int by decimals amount and returns it
func ReduceBigInt(n *big.Int, precision int) *big.Int {
	mul := math.Pow(10.0, float64(precision))
	mulBig := new(big.Int).SetUint64(uint64(mul))

	nBig := new(big.Int).Div(n, mulBig)

	return nBig
}

// ExtractFees fee like (30% -> use 0.3), applies a fee reduction to n
func ExtractFees(n *big.Int, fee float64, precision int) *big.Int {
	f := 1 - fee
	fPrec := math.Pow(10.0, float64(precision))
	fBig := new(big.Int).SetInt64(int64(f * fPrec))

	out := new(big.Int).Mul(n, fBig)
	out = ReduceBigInt(out, precision)
	return out
}
