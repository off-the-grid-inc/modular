package modular

import (
	"sync"
	"crypto/rand"
	"errors"
	"math/big"
)

type Int big.Int

// First, create a struct that contains the value we want to return
// along with a mutex instance
type Prime struct {
	value *Int
	m     sync.Mutex
}

// Default prime (256-bit secp256k1 EC order)
var (
	val, _ = IntFromString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	p = &Prime {
		value: val,
	}
)

func GetP() *Int {
	p.m.Lock()

	defer p.m.Unlock()
	return p.value
}

func SetP(val *Int) *Int {
	p.m.Lock()

	p.value = val
	defer p.m.Unlock()

	return p.value
}

// ARITHMETIC OPERATIONS

// Add - Modular addition of (an arbitrary number of) field elements
func (n *Int) Add(nums ...*Int) *Int {
	out := big.NewInt(0)
	for _, n := range nums {
		out.Add(out, (*big.Int)(n))
	}
	out.Mod(out, (*big.Int)(GetP()))
	*n = (Int)(*out)
	return n
}

// Mul - Modular multiplication of two field elements
func (n *Int) Mul(x, y *Int) *Int {
	out := new(big.Int).Mul((*big.Int)(x), (*big.Int)(y))
	out.Mod(out, (*big.Int)(GetP()))
	*n = (Int)(*out)
	return n
}

// Sub - Modular Subtraction of two field elements
func (n *Int) Sub(x, y *Int) *Int {
	out := new(big.Int).Sub((*big.Int)(x), (*big.Int)(y))
	out.Mod(out, (*big.Int)(GetP()))
	*n = (Int)(*out)
	return n
}

// LinearCombination is the dot product of two vectors in a finite field. Vectors should have the same dimensionality for proper behavior
func (n *Int) LinearCombination(vec1 []*Int, vec2 []*Int) *Int {
	out := big.NewInt(0)
	for i := range vec1 {
		out.Add(out, new(big.Int).Mul((*big.Int)(vec1[i]), (*big.Int)(vec2[i])))
	}
	out.Mod(out, (*big.Int)(GetP()))
	*n = (Int)(*out)
	return n
}

// ModInverse is a custom implementation for the modular inverse of a field element
func ModInverse(number *Int) *Int {
	copy := big.NewInt(0).Set((*big.Int)(number))
	pcopy := big.NewInt(0).Set((*big.Int)(GetP()))
	x := big.NewInt(0)
	y := big.NewInt(0)

	copy.GCD(x, y, pcopy, copy)

	result := big.NewInt(0).Set((*big.Int)(GetP()))

	result.Add(result, y)
	result.Mod(result, (*big.Int)(GetP()))
	return (*Int)(result)
}

// Exp - exponentiate in a finite field
func (n *Int) Exp(base, exp *Int) *Int {
	*n = (Int)(*new(big.Int).Exp((*big.Int)(base), (*big.Int)(exp), (*big.Int)(GetP())))
	return n
}

// Cmp compares two finite field elements
func (n *Int) Cmp(x *Int) int {
	return (*big.Int)(n).Cmp((*big.Int)(x))
}

// Bytes returs the byte array representation of a field element
func (n *Int) Bytes() []byte {
	return (*big.Int)(n).Bytes()
}

// String returns the string representation
func (n *Int) String() string {
	return (*big.Int)(n).String()
}

// Helpers

// NewInt creates a modular Int from int64
func NewInt(i int64) *Int {
	num := big.NewInt(i)
	return (*Int)(num).Mod()
}

// Note: IntFrom... methods are not automatically reduced mod P

// IntFromBytes creates a modular Int from a byte array
func IntFromBytes(b []byte) *Int {
	num := new(big.Int).SetBytes(b)
	return (*Int)(num)
}

// IntFromBig creates a modular Int from a big Int
func IntFromBig(num *big.Int) *Int {
	return (*Int)(new(big.Int).Set(num))
}

// IntFromString creates a modular Int from a string representation of an integer in a specific base.
func IntFromString(str string, base int) (*Int, error) {
	num, err := new(big.Int).SetString(str, base)
	if !err {
		return nil, errors.New("Could not set string")
	}
	return (*Int)(num), nil
}

// IsModP checks whether a modular Int actually lies within the field order
func (n *Int) IsModP() bool {
	if n.Cmp(GetP()) == -1 {
		return true
	}
	return false
}

func (n *Int) AsBig() *big.Int {
	return (*big.Int)(n)
}

func (n *Int) Mod() *Int {
	*n = (Int)(*new(big.Int).Mod((*big.Int)(n), (*big.Int)(GetP())))
	return n
}

// RandInt creates a new random modular Int within the range [0, p)
func RandInt() (*Int, error) {
	max := big.NewInt(0).Set((*big.Int)(GetP()))
	max.Sub(max, big.NewInt(1))
	result, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, err
	}
	return (*Int)(result), nil
}
