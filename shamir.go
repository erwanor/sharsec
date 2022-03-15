package sharsec

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/erwanor/sharsec/curvewrapper"
	"github.com/erwanor/sharsec/finitefield"
)

type Shamir struct {
	ec elliptic.Curve
}

type ClearShare struct {
	SID   *finitefield.FpInt
	Value curvewrapper.Point
}

type EncryptedShare = ClearShare

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
		/*
		 Note: modular exponentiation is not a cryptographically constant-time operation.
		 Idea: https://www.cse.buffalo.edu/srds2009/escs2009_submission_Gopal.pdj
		 Language proposal to track: https://github.com/golang/go/issues/20655
		*/

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

func (s *Shamir) Split(secret []byte, threshold int, pubkeys []Key) []ClearShare {
	sec := new(big.Int).SetBytes(secret)

	polynomial := s.GeneratePolynomial(threshold)
	polynomial[0] = sec

	var shares []ClearShare
	for i := 0; i < len(pubkeys); i++ {
		order := s.ec.Params().N
		// Take N arbitrary points
		// The secret is hidden in the constant term of the polynomial.
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

func (s *Shamir) Combine(c []ClearShare) []byte {
	fieldOrder := s.ec.Params().N
	field := finitefield.NewField(fieldOrder)
	pooledSecret := field.Zero()

	for i := 0; i < len(c); i++ {
		/*
			Lagrange interpolation:

			Quick summary:
			    Given K-1 points, we want to obtain a polynomial P
			    such that P(x_1) = y_1 ...P(x_{k-1}) = y_{k-1}.

			    We seek to compute P(0) (recall the secret is hidden in the constant term.
			    f(x) = SUM [i = 0 -> (k-1)] ( y_i * Lagrange_i(x) )
			    where Lagrange_i (x) = PROD [j = 0 -> (k-1) with j =/= i] (x - x_j) * (x_i - x_j)^-1 [in the field F_p] */

		lagrange := field.One()
		lagNum := field.One()
		lagDenum := field.One()
		lagDenumInverse := field.One()
		term := field.One()

		for j := 0; j < len(c); j++ {
			if i == j {
				continue
			}

			deltaNum := field.Zero()
			deltaDenum := field.Zero()

			// We want to interpolate at value x = 0
			// Therefore, x - x_j = 0 - x_j = -x_j
			deltaNum.Sub(deltaNum, field.NewInt(c[j].Value.X))

			lagNum.Mul(lagNum, deltaNum)

			deltaDenum.Sub(field.NewInt(c[i].Value.X), field.NewInt(c[j].Value.X))

			lagDenum.Mul(lagDenum, deltaDenum)
		}
		lagDenumInverse.ModInv(lagDenum)
		lagrange.Mul(lagNum, lagDenumInverse)
		term.Mul(field.NewInt(c[i].Value.Y), lagrange)
		pooledSecret.Add(pooledSecret, term)
	}
	return pooledSecret.Value.Bytes()
}
