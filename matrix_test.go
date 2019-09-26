package modular

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func PrintArray(arr []*Int) {
	for _, a := range arr {
		b := a.Value
		fmt.Printf("%d, ", b)
	}
	fmt.Println("")
}

func TestBasicMatrix(t *testing.T) {
	require := require.New(t)

	data := make([]*Int, 10)
	for i := range data {
		v, err := RandInt(defaultP)
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
	m.ScalarMul(new(Int).Exp(NewInt(2, defaultP), NewInt(256, defaultP)))
	require.Equal(0, m.values[9].Cmp(data[9].Mul(data[9], new(Int).Exp(NewInt(2, defaultP), NewInt(256, defaultP)))), "scalar mult failed")

	// Test Set Row/Col
	m.SetRow(1, []*Int{NewInt(1, defaultP), NewInt(1, defaultP), NewInt(1, defaultP), NewInt(1, defaultP), NewInt(1, defaultP)})
	require.Equal(0, m.values[0].Cmp(m.values[4]), "set row failed")
	m.SetCol(5, []*Int{NewInt(1, defaultP), NewInt(1, defaultP)})
	require.Equal(0, m.values[9].Cmp(m.values[4]), "set column failed")

}

func TestMultiplication(t *testing.T) {
	require := require.New(t)

	m1 := NewMatrix(2, 3, []*Int{NewInt(1, defaultP), NewInt(2, defaultP), NewInt(3, defaultP), NewInt(4, defaultP), NewInt(5, defaultP), NewInt(6, defaultP)})
	m2 := NewMatrix(3, 1, []*Int{NewInt(3, defaultP), NewInt(2, defaultP), NewInt(1, defaultP)})
	res, err := new(Matrix).Mul(m1, m2)
	require.NoError(err)
	require.Equal(len(res.values), 2, "wrong structure")
	require.Equal(0, res.values[0].Cmp(NewInt(10, defaultP)), "multiplication failed")
	require.Equal(0, res.values[1].Cmp(NewInt(28, defaultP)), "multiplication failed")
}

func TestInverse(t *testing.T) {
	require := require.New(t)

	// Test Gauss Jordan
	linearSystem := NewMatrix(2, 2, []*Int{NewInt(-1, defaultP), NewInt(2, defaultP), NewInt(1, defaultP), NewInt(0, defaultP)})
	ls := linearSystem.Represent2D()
	linearSystemResult := []*Int{NewInt(13, defaultP), NewInt(1, defaultP)}

	// Warning: the system is modified during the Gaussian reduction process
	// So it's better to pass a copy
	result, err := GaussJordan(linearSystem.Represent2D(), linearSystemResult)

	require.NoError(err, "gauss jordan failed")

	lhs := NewInt(0, defaultP)
	for i := 0; i < len(result); i++ {
		factor := new(Int).Mul(ls[0][i], result[i])
		lhs = new(Int).Add(lhs, factor)
	}

	require.Equal(0, linearSystemResult[0].Cmp(lhs), "System 1 failed")

	// require.Equal(0, result[0].Cmp(NewInt(1)), "gauss jordan failed")
	// require.Equal(0, result[1].Cmp(NewInt(1).Mod()), "gauss jordan failed")

	// Second matrix
	linearSystem2 := NewMatrix(3, 3, []*Int{NewInt(1, defaultP), NewInt(2, defaultP), NewInt(3, defaultP), NewInt(4, defaultP), NewInt(5, defaultP), NewInt(6, defaultP), NewInt(7, defaultP), NewInt(8, defaultP), NewInt(9, defaultP)})
	ls2 := linearSystem2.Represent2D()

	linearSystemResult2 := []*Int{NewInt(0, defaultP), NewInt(1, defaultP), NewInt(2, defaultP)}
	_ = linearSystemResult2

	result2, err := GaussJordan(ls2, linearSystemResult2)

	lhs = NewInt(0, defaultP)
	for i := 0; i < len(result2); i++ {
		factor := new(Int).Mul(ls2[0][i], result2[i])
		lhs = new(Int).Add(lhs, factor)
	}
	require.Equal(0, linearSystemResult2[0].Cmp(lhs), "System 2 failed")

	// Test Inverses
	m := NewMatrix(2, 2, []*Int{NewInt(7, defaultP), NewInt(-3, defaultP), NewInt(-2, defaultP), NewInt(1, defaultP)})
	inv, err := m.Inverse()
	require.NoError(err, "inverse failed")
	_ = inv
	require.Equal(0, inv.values[0].Cmp(NewInt(1, defaultP)), "inverse failed")

}
