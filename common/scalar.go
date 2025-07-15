package common

import "math/big"

type Scalar struct {
	field Field
	value *big.Int
}

func (s Scalar) Value() *big.Int {
	return s.value
}

func (s Scalar) Field() Field {
	return s.field
}

func (s Scalar) IsZero() bool {
	return s.value.Cmp(s.field.zero) == 0
}

func (s Scalar) Equals(a *Scalar) bool {
	return s.value.Cmp(a.value) == 0 && s.field.modulus.Cmp(a.field.modulus) == 0
}

func (s *Scalar) Add(a *Scalar) *Scalar {
	s.value.Add(s.value, a.value)
	if s.value.Cmp(s.field.modulus) >= 0 {
		s.value.Sub(s.value, s.field.modulus)
	}
	return s
}

func (s *Scalar) Sub(a *Scalar) *Scalar {
	s.value.Sub(s.value, a.value)
	if s.value.Sign() < 0 {
		s.value.Add(s.value, s.field.modulus)
	}
	return s
}

func (s *Scalar) Mul(a *Scalar) *Scalar {
	s.value.Mul(s.value, a.value)
	s.field.Mod(s.value)
	return s
}

func (s *Scalar) Div(a *Scalar) *Scalar {
	inv := new(big.Int).ModInverse(a.value, s.field.modulus)
	if inv == nil {
		panic("division by zero or non-invertible scalar in field")
	}
	s.value.Mul(s.value, inv)
	s.field.Mod(s.value)
	return s
}

func (s *Scalar) Inv() *Scalar {
	inv := new(big.Int).ModInverse(s.value, s.field.modulus)
	if inv == nil {
		panic("no modular inverse exists for this scalar")
	}
	s.value.Set(inv)
	return s
}
