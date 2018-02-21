package sharsec

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
	"sharsec/curvewrapper"
	"sharsec/finitefield"
)

type Shamir struct {
	ec elliptic.Curve
}

// For future iterations.
type EncryptedShare struct {
	SID   *finitefield.FpInt
	Value curvewrapper.Point
}

type ClearShare = EncryptedShare

func NewShamirSystem(curve elliptic.Curve) *Shamir {
	return &Shamir{
		ec: curve,
	}
}

func (s *Shamir) NewKey() Key {
	var k Key
	k.priv, k.pub.X, k.pub.Y, _ = elliptic.GenerateKey(s.ec, rand.Reader)
	k.pub.Curve = s.ec
	return k
}

type ShamirPoly []*big.Int

func (s *Shamir) GeneratePolynomial(deg int) ShamirPoly {
	var p []*big.Int
	// We want deg+1 elements because the secret is hidden in the constant term
	// of the polynomial. That is, consider F(x)_n = SEC + A_1*x^1 + ... + A_n * x^n (A_0 = SEC)
	for i := 0; i < deg; i++ {
		k := s.NewKey()
		var coeff big.Int
		// Note that SetBytes assume Big-Endian uint
		coeff.SetBytes(k.priv)
		coeff.Mod(&coeff, s.ec.Params().N)
		p = append(p, &coeff)
	}
	return p
}

func (p ShamirPoly) Eval(x *big.Int, mod *big.Int) *big.Int {
	var total big.Int
	for i := 0; i < len(p); i++ {
		var current big.Int

		// Note that modular exponentiation is not a cryptographically constant-time operation
		// IDEA: https://www.cse.buffalo.edu/srds2009/escs2009_submission_Gopal.pdf
		// There's a proposal to provide built-in constant-time arithmetic: https://github.com/golang/go/issues/20654
		current.Exp(x, big.NewInt(int64(i)), nil)
		current.Mul(p[i], &current)
		current.Mod(&current, mod)
		total.Add(&total, &current)
		total.Mod(&total, mod)
	}
	return &total
}

func (p ShamirPoly) String() {
	fmt.Printf("%d", p[0])
	for i := 1; i < len(p); i++ {
		fmt.Printf(" + %d * x^%d", p[i], i)
	}
	fmt.Printf("\n")
}

func (s *Shamir) Split(sec *big.Int, threshold int, pubkeys []Key) []ClearShare {
	polynomial := s.GeneratePolynomial(threshold)
	polynomial[0] = sec

	polynomial.String()

	var shares []ClearShare
	for i := 0; i < len(pubkeys); i++ {
		order := s.ec.Params().N
		// We don't want to evaluate at 0 since that would reveal the secret
		image := polynomial.Eval(big.NewInt(int64(i+1)), order)
		currShare := ClearShare{
			SID: &finitefield.FpInt{
				Value: big.NewInt(int64(i + 1)),
				Order: s.ec.Params().N,
			},
			Value: curvewrapper.NewPoint(big.NewInt(int64(i+1)), image, s.ec),
		}
		shares = append(shares, currShare)
	}

	return shares
}

func (e EncryptedShare) Decrypt(priv []byte) ClearShare {
	var privKey big.Int
	v := e.Value.ScalarDiv(privKey.SetBytes(priv))
	return ClearShare{
		SID:   e.SID,
		Value: v,
	}
}

func (s *Shamir) Combine(c []ClearShare) *big.Int {
	fieldOrder := s.ec.Params().N
	pooledSecret := finitefield.NewFpInt(big.NewInt(0), fieldOrder)

	for i := 0; i < len(c); i++ {
		// Lagrange interpolation
		// We will proceed as follow:
		// f(x) = SUM [i = 0...(k-1)] ( y_i * Lagrange_i(x) )
		// where Lagrange_i(x) = PROD [j = 0..(i-1)(i+1)..(k-1)] (x - x_j) * (x_i - x_j)^-1 [in the field F_p]
		// See: http://www.math.usm.edu/lambers/mat772/fall10/lecture5.pdf for more details
		lagrange := finitefield.NewFpInt(big.NewInt(1), fieldOrder)
		lagNum := finitefield.NewFpInt(big.NewInt(1), fieldOrder)
		lagDenum := finitefield.NewFpInt(big.NewInt(1), fieldOrder)
		lagDenumInverse := finitefield.NewFpInt(big.NewInt(1), fieldOrder)
		term := finitefield.NewFpInt(big.NewInt(1), fieldOrder)

		for j := 0; j < len(c); j++ {
			if i == j {
				continue
			}

			deltaNum := finitefield.NewFpInt(big.NewInt(0), fieldOrder)
			deltaDenum := finitefield.NewFpInt(big.NewInt(0), fieldOrder)

			deltaNum.Sub(
				deltaNum,
				finitefield.NewFpInt(c[j].Value.X, fieldOrder))

			lagNum.Mul(lagNum, deltaNum)

			deltaDenum.Sub(
				finitefield.NewFpInt(c[i].Value.X, fieldOrder),
				finitefield.NewFpInt(c[j].Value.X, fieldOrder))

			lagDenum.Mul(lagDenum, deltaDenum)
		}
		lagDenumInverse.ModInv(lagDenum)
		lagrange.Mul(lagNum, lagDenumInverse)
		term.Mul(finitefield.NewFpInt(c[i].Value.Y, fieldOrder), lagrange)
		pooledSecret.Add(pooledSecret, term)
	}
	return pooledSecret.Value
}
