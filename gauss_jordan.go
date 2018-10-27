package modular

import (
	"errors"
)

// NonZeroEntry returns the index of the first non-zero entry of the array, along with such value.
// If no such value exists then return -1 as the index and nil as the value
func NonZeroEntry(arr []*Int) (int, *Int, error) {
	if len(arr) == 0 {
		return 0, nil, errors.New("Empty array is not valid")
	}
	value := arr[0]
	index := -1
	for i, a := range arr {
		if a.Value.Cmp(ZERO) != 0 {
			value = a
			index = i
			break
		}
	}
	if index == -1 {
		return index, nil, nil
	}
	return index, value, nil
}

// Appends arr as the last column to matrix
func createExtendedArray(matrix [][]*Int, arr []*Int, n int) [][]*Int {
	nrows := len(matrix)
	ncols := len(matrix[0])

	extended := make([][]*Int, nrows)
	for i := 0; i < nrows; i++ {
		extended[i] = make([]*Int, ncols+1)

		for j := 0; j < ncols; j++ {
			extended[i][j] = matrix[i][j]
		}
		extended[i][ncols] = arr[i]
	}
	return extended
}

// Extracts column j from matrix
func ExtractColumn(matrix [][]*Int, j int) []*Int {
	column := make([]*Int, len(matrix))
	for i, a := range matrix {
		column[i] = a[j]
	}
	return column
}

// Swaps rows i and j of matrix
func SwapRows(matrix [][]*Int, i int, j int) {
	matrix[i], matrix[j] = matrix[j], matrix[i]
}

// GaussJordan solves the linear system linearSystem*x = linearSystemResult
// It uses Gauss Jordan algorithm
// An error is returned if no solutions are found
// If multiple solutions exist, the one with 0's for the free variables is chosen
func GaussJordan(linearSystem [][]*Int, linearSystemResult []*Int) ([]*Int, error) {
	nrows := len(linearSystem)
	// Check that matrix is non-empty
	if nrows == 0 {
		return nil, errors.New("Empty matrix is not allowed")
	}
	ncols := len(linearSystem[0])
	// Check that all rows have the same length
	for _, a := range linearSystem {
		if len(a) != ncols {
			return nil, errors.New("There are rows with different length")
		}
	}
	dimension_result := len(linearSystemResult)
	// Check the dimensions of the matrix vs the result vector
	if (nrows*ncols == 0) || (ncols != dimension_result) {
		return nil, errors.New("Wrong dimensions for the system")
	}

	extendedMatrix := createExtendedArray(linearSystem, linearSystemResult, ncols)

	// The matrix has been extended, the number of columns increase
	ncols++

	index_current_row := 0
	index_current_column := 0
	current_column := make([]*Int, len(extendedMatrix))

	for (index_current_row < nrows) && (index_current_column < ncols) {
		// Find the k-th pivot:
		current_column = ExtractColumn(extendedMatrix, index_current_column)
		prime := current_column[0].Base

		index_max, non_zero_entry, err := NonZeroEntry(current_column[index_current_row:])

		if err != nil {
			return nil, err
		}
		// Index is relative, we make it absolute
		index_max = index_max + index_current_row
		if non_zero_entry == nil {
			// No pivot in this column, pass to next column
			index_current_column++
		} else {
			SwapRows(extendedMatrix, index_current_row, index_max)
			// Do for all rows below pivot:
			for i := index_current_row + 1; i < nrows; i++ {
				factor := new(Int).Mul(extendedMatrix[i][index_current_column], ModInverse(extendedMatrix[index_current_row][index_current_column]))
				// Fill with zeros the lower part of pivot column:
				extendedMatrix[i][index_current_column] = NewInt(0, prime)
				// Do for all remaining elements in current row:
				for j := index_current_column + 1; j < ncols; j++ {
					num := new(Int).Mul(extendedMatrix[index_current_row][j], factor)
					extendedMatrix[i][j] = new(Int).Sub(extendedMatrix[i][j], num)
				}
			}
			// Increase pivot row and column
			index_current_row++
			index_current_column++
		}
	}

	// The matrix is upper triangular, but not in row echelon form yet
	// List holding the pivots positions for each row
	pivots := make([]int, nrows)
	for i := 0; i < nrows; i++ {
		pivots[i] = -1
	}
	for i := 0; i < nrows; i++ {
		for j := 0; j < ncols; j++ {
			if extendedMatrix[i][j].Value.Cmp(ZERO) != 0 {
				pivots[i] = j
				break
			}
		}
		// If the pivot is the entry in the last column then an error is returned since in the case no solutions exist
		if pivots[i] == ncols-1 {
			err := errors.New("There are no solutions")
			return nil, err
		}
	}

	result := make([]*Int, ncols-1)
	prime := extendedMatrix[0][0].Base
	for j := 0; j < ncols-1; j++ {
		result[j] = NewInt(0, prime)
	}

	// Calculate result in position pivots[i]
	for i := nrows - 1; i >= 0; i-- {
		if pivots[i] != -1 {
			result[pivots[i]] = extendedMatrix[i][ncols-1]
			for i_ := nrows - 1; i_ > i; i_-- {
				if pivots[i_] != -1 {
					mult := new(Int).Mul(result[pivots[i_]], extendedMatrix[i][pivots[i_]])

					result[pivots[i]] = new(Int).Sub(result[pivots[i]], mult)
				}

			}

			result[pivots[i]] = new(Int).Mul(result[pivots[i]], ModInverse(extendedMatrix[i][pivots[i]]))
		}
	}
	return result, nil
}
