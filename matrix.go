package modular

import (
	"errors"
)

type Matrix struct {
	nRow      int
	nCol      int
	values   []*Int
}

// NewMatrix creates a new (unexported) matrix struct.
func NewMatrix(r, c int, vals []*Int) *Matrix {
	space := r*c - len(vals)
	for space > 0 {
		vals = append(vals, NewInt(0))
		space--
	}
	return &Matrix{
		nRow: r,
		nCol: c,
		values: vals,
	}
}

// GetRow returns a rwo of a matrix. It starts from 1 rather than 0.
func (m *Matrix) GetRow(r int) []*Int {
	return m.values[(r-1)*m.nCol:r*m.nCol]
}

// SetRow resets a row. It also starts from 1 rather than 0.
func (m *Matrix) SetRow(r int, row []*Int) *Matrix {
	i := 0
	for i < m.nCol {
		m.values[(r-1)*m.nCol + i] = row[i]
		i++
	}
	return m
}

// GetCol returns a column of a matrix. It starts from 1 rather than 0.
func (m *Matrix) GetCol(c int) []*Int {
	c--
	res := make([]*Int, m.nRow)
	for i := range res {
		res[i] = m.GetRow(i+1)[c]
	}
	return res
}

func (m *Matrix) SetCol(c int, col []*Int) *Matrix {
	i := 0
	for i < m.nRow {
		m.values[i*m.nCol+(c-1)] = col[i]
		i++
	}
	return m
}

// ScalarMul multiplies a matrix by a scalar
func (m *Matrix) ScalarMul(i *Int) *Matrix {
	for _, v := range m.values {
		v.Mul(v, i)
	}
	return m
}

func (m *Matrix) Represent2D() [][]*Int {
	mat := make([][]*Int, m.nRow)
	for i := range mat {
		row := m.GetRow(i+1)
		mat[i] = make([]*Int, len(row))
		for j := range mat[i] {
			mat[i][j] = IntFromBig(row[j].AsBig())
		}
	}
	return mat
}

func (m *Matrix) Copy() *Matrix {
	vals := make([]*Int, len(m.values))
	for i := range vals {
		vals[i] = IntFromBig(m.values[i].AsBig())
	}
	return NewMatrix(m.nRow, m.nCol, vals)
}

// Mul does matrix multiplication, returns a new matrix if the dimensions are acceptable for multiplication, else an error.
func (m *Matrix) Mul(x, y *Matrix) (*Matrix, error) {
	if x.nCol != y.nRow {
		return nil, errors.New("mismatched dimensions, cannot multiply")
	}
	vals := make([]*Int, x.nRow*y.nCol)
	i := 1
	j := 1
	idex := 0
	for i <= x.nRow {
		j = 1
		for j <= y.nCol {
			vals[idex] = new(Int).LinearCombination(x.GetRow(i),  y.GetCol(j))
			idex++
			j++
		}
		i++
	}
	m = &Matrix{
		nRow: y.nCol,
		nCol: x.nRow,
		values: vals,
	}
	return m, nil
}

func (m *Matrix) Inverse() (*Matrix, error) {
	if m.nRow != m.nCol {
		return nil, errors.New("only square matrices are invertible")
	}
	inverse := NewMatrix(m.nRow, m.nCol, []*Int{})
	i := 1
	for i < m.nRow+1 {
		col, err := GaussJordan(m.Represent2D(), GetI(m.nRow).GetRow(i))
		if err != nil {
			return nil, err
		}
		inverse.SetCol(i, col)
		i++
	}

	return inverse, nil
}

func GetI(c int) *Matrix {
	I := NewMatrix(c,c, []*Int{})
	i := 0
	idex := 0
	for i < c {
		idex = c*i + i
		I.values[idex] = NewInt(1)
		i++
	}
	return I
}