package finitefield

import (
	"math/big"
)

type FpInt struct {
	Value *big.Int
	Order *big.Int
}

func NewFpInt(v, order *big.Int) *FpInt {
	return &FpInt{
		Value: v.Mod(v, order),
		Order: order,
	}
}
