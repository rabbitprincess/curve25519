package common

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

type Field struct {
	zero    *big.Int
	one     *big.Int
	modulus *big.Int
	factor  *big.Int
	shift   uint
}

func NewFile(modulus *big.Int, factor *big.Int, shift uint) *Field {
	return &Field{
		zero:    big.NewInt(0),
		one:     big.NewInt(1),
		modulus: modulus,
		factor:  factor,
		shift:   shift,
	}
}

func (f Field) Mod(a *big.Int) {
	var t big.Int

	t.Mul(a, f.factor)
	t.Rsh(&t, f.shift)
	t.Mul(&t, f.modulus)
	a.Sub(a, &t)
	if a.Cmp(f.modulus) >= 0 {
		a.Sub(a, f.modulus)
	}
}

func (f Field) Zero() Scalar {
	return Scalar{
		field: f,
		value: f.zero,
	}
}

func (f Field) One() Scalar {
	return Scalar{
		field: f,
		value: f.one,
	}
}

func (f Field) Modulus() *big.Int {
	return f.modulus
}

func (f Field) ByteLen() int {
	return (f.modulus.BitLen() + 7) / 8
}

func (f Field) EncodeScalar(scalar *Scalar) []byte {
	b := make([]byte, f.ByteLen())
	scalar.value.FillBytes(b)
	return b
}

func (f Field) DecodeScalar(b []byte) (*Scalar, error) {
	if len(b) != f.ByteLen() {
		return nil, fmt.Errorf("invalid scalar length")
	}
	var r big.Int
	r.SetBytes(b)
	if r.Cmp(f.modulus) >= 0 {
		return nil, fmt.Errorf("scalar larger than modulus")
	}
	return &Scalar{
		field: f,
		value: &r,
	}, nil
}

func (f Field) NewScalarWithRandom(r io.Reader) (*Scalar, error) {
	if r == nil {
		r = rand.Reader
	}
	n, err := rand.Int(r, f.modulus)
	if err != nil {
		return nil, err
	}
	return &Scalar{
		field: f,
		value: n,
	}, nil
}

func (f Field) NewScalarWithModularReduction(value *big.Int) *Scalar {
	var reduced big.Int
	if value.Sign() >= 0 && value.Cmp(f.modulus) < 0 {
		reduced.Set(value)
	} else {
		reduced.Mod(value, f.modulus)
	}
	return &Scalar{
		field: f,
		value: &reduced,
	}
}
