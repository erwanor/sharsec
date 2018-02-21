package curvewrapper

import (
	"crypto/elliptic"
	"math/big"
)

type Point struct {
	X     *big.Int
	Y     *big.Int
	Curve elliptic.Curve
}

func NewPoint(x, y *big.Int, ec elliptic.Curve) Point {
	return Point{
		X:     x,
		Y:     y,
		Curve: ec,
	}
}

func (p Point) Add(a, b Point) Point {
	var result Point
	result.Curve = p.Curve
	result.X, result.Y = result.Curve.Add(a.X, a.Y, b.X, b.Y)
	return result
}

func (p Point) ScalarMul(k *big.Int) Point {
	var result Point
	result.X, result.Y = p.Curve.ScalarMult(p.X, p.Y, k.Bytes())
	result.Curve = p.Curve
	return result
}

func (p Point) ScalarDiv(k *big.Int) Point {
	scalarInverse := k.ModInverse(k, p.Curve.Params().N)
	return p.ScalarMul(scalarInverse)
}
