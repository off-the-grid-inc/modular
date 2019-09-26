# Modular
Modular is a golang math package for cyclic group arithmetic. We create a modular.Int type (in the style of big.Int), but automatically works within a prime order group, handling all modular reductions for convenience. We also built a modest modular.Matrix type with a few rudimentary matrix operations for working with discrete matrices, since most linear algebra packages work with rational/floating point types and don't seamlessly support modular arithmetic in their matrix operations. 

`$ go get "github.com/off-the-grid-inc/modular"`

## Usage

```
package main

import (
  "github.com/off-the-grid-inc/modular"
  "math/big"
)

func main() {
    prime := big.NewInt(23)
    fieldElement := modular.NewInt(25, prime) // i.e. 2 mod 23
    randomElement := modular.RandInt(prime)  // generate random field elements
    fieldElement.Mul(fieldElement, randomElement) // modular arithmetic
    fieldElement.Add(fieldElement, modular.NewInt(1, prime))
    fieldElement.Exp(fieldElement, randomElement)
    invElement := modular.ModInverse(fieldElement) // modular inverse
    newMatrix := modular.NewMatrix(3,3,nil) // 3x3 finite feild matrix with nil values
    row := []*modular.Int{modular.NewInt(0, prime), modular.NewInt(1, prime), modular.NewInt(2, prime)}
    for i := 1; i<4; i++ {
      newMatrix.SetRow(i, row)
    }
    newMatrix.GetColumn(1) // i.e. [0, 0, 0]
    invMatrix, err := newMatrix.Inverse()  // inverted finite field matrix
}
```
