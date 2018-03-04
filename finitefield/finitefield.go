package finitefield

import (
	"math/big"
)

type FpInt struct {
	Value *big.Int
	Order *big.Int
}

type Field struct {
	Order *big.Int
}

func NewFpInt(v, order *big.Int) *FpInt {
	return &FpInt{
		Value: v.Mod(v, order),
		Order: order,
	}
}

func (k *FpInt) Add(a, b *FpInt) *FpInt {
	k.Value.Add(a.Value, b.Value)
	k.Value.Mod(k.Value, a.Order)
	return k
}

func (k *FpInt) Sub(a, b *FpInt) *FpInt {
	k.Value.Sub(a.Value, b.Value)
	k.Value.Mod(k.Value, a.Order)
	return k
}

func (k *FpInt) Mul(a, b *FpInt) *FpInt {
	k.Value.Mul(a.Value, b.Value)
	k.Value.Mod(k.Value, a.Order)
	return k
}

func (k *FpInt) ModInv(a *FpInt) *FpInt {
	k.Value.ModInverse(a.Value, a.Order)
	return k
}

func NewField(order *big.Int) Field {
	return Field{
		Order: order,
	}
}

func (f Field) NewInt(v *big.Int) *FpInt {
	return NewFpInt(v, f.Order)
}

func (f Field) Zero() *FpInt {
	return NewFpInt(big.NewInt(int64(0)), f.Order)
}

func (f Field) One() *FpInt {
	return NewFpInt(big.NewInt(int64(1)), f.Order)
}
