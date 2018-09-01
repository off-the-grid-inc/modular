package modular

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func PrintArray(arr []*Int) {
	for _, a := range arr {
		b := (*big.Int)(a)
		fmt.Printf("%d, ", b)
	}
	fmt.Println("")
}

func TestBasicMatrix(t *testing.T) {
	require := require.New(t)

	data := make([]*Int, 10)
	for i := range data {
		v, err := RandInt()
		require.NoError(err)
		data[i] = v
	}

	// Test Create Matrix
	m := NewMatrix(2, 5, data)

	// Test Get Row/Col
	row := m.GetRow(2)
	require.Equal(row[4].Cmp(data[9]), 0, "get row failed")
	col := m.GetCol(5)
	require.Equal(row[4].Cmp(col[1]), 0, "get column failed")

	// Test Scalar Multiplication
	m.ScalarMul(new(Int).Exp(NewInt(2), NewInt(256)))
	require.Equal(0, m.values[9].Cmp(data[9].Mul(data[9], new(Int).Exp(NewInt(2), NewInt(256)))), "scalar mult failed")

	// Test Set Row/Col
	m.SetRow(1, []*Int{NewInt(1), NewInt(1), NewInt(1), NewInt(1), NewInt(1)})
	require.Equal(0, m.values[0].Cmp(m.values[4]), "set row failed")
	m.SetCol(5, []*Int{NewInt(1), NewInt(1)})
	require.Equal(0, m.values[9].Cmp(m.values[4]), "set column failed")

}

func TestMultiplication(t *testing.T) {
	require := require.New(t)

	m1 := NewMatrix(2, 3, []*Int{NewInt(1), NewInt(2), NewInt(3), NewInt(4), NewInt(5), NewInt(6)})
	m2 := NewMatrix(3, 1, []*Int{NewInt(3), NewInt(2), NewInt(1)})
	res, err := new(Matrix).Mul(m1, m2)
	require.NoError(err)
	require.Equal(len(res.values), 2, "wrong structure")
	require.Equal(0, res.values[0].Cmp(NewInt(10)), "multiplication failed")
	require.Equal(0, res.values[1].Cmp(NewInt(28)), "multiplication failed")
}

func TestInverse(t *testing.T) {
	SetP(NewInt(5))
	require := require.New(t)

	// Test Gauss Jordan
	linearSystem := NewMatrix(2, 2, []*Int{NewInt(-1), NewInt(2), NewInt(1), NewInt(0)})
	ls := linearSystem.Represent2D()
	linearSystemResult := []*Int{NewInt(13), NewInt(1)}

	// Warning: the system is modified during the Gaussian reduction process
	// So it's better to pass a copy
	result, err := GaussJordan(linearSystem.Represent2D(), linearSystemResult)

	require.NoError(err, "gauss jordan failed")

	PrintArray(result)
	PrintArray(linearSystem.Represent2D()[0])
	lhs := NewInt(0)
	for i := 0; i < len(result); i++ {
		factor := new(Int).Mul(ls[0][i], result[i])
		fmt.Printf("%v * %v = %v\n", (*big.Int)(ls[0][i]), (*big.Int)(result[i]), (*big.Int)(factor))
		lhs = new(Int).Add(lhs, factor)
	}
	fmt.Printf("Result: %v\n", (*big.Int)(lhs))

	require.Equal(0, linearSystemResult[0].Cmp(lhs), "System 1 failed")

	// require.Equal(0, result[0].Cmp(NewInt(1)), "gauss jordan failed")
	// require.Equal(0, result[1].Cmp(NewInt(1).Mod()), "gauss jordan failed")

	// Second matrix
	linearSystem2 := NewMatrix(3, 3, []*Int{NewInt(1), NewInt(2), NewInt(3), NewInt(4), NewInt(5), NewInt(6), NewInt(7), NewInt(8), NewInt(9)})
	ls2 := linearSystem2.Represent2D()

	linearSystemResult2 := []*Int{NewInt(0), NewInt(1), NewInt(2)}
	_ = linearSystemResult2

	result2, err := GaussJordan(ls2, linearSystemResult2)

	lhs = NewInt(0)
	for i := 0; i < len(result2); i++ {
		factor := new(Int).Mul(ls2[0][i], result2[i])
		lhs = new(Int).Add(lhs, factor)
	}
	require.Equal(0, linearSystemResult2[0].Cmp(lhs), "System 2 failed")

	// Test Inverses
	m := NewMatrix(2, 2, []*Int{NewInt(7), NewInt(-3).Mod(), NewInt(-2).Mod(), NewInt(1)})
	inv, err := m.Inverse()
	require.NoError(err, "inverse failed")
	_ = inv
	require.Equal(0, inv.values[0].Cmp(NewInt(1)), "inverse failed")

}
