package main

import (
  "fmt"
  "math"
  "math/rand"
)

// type for it?


// grid input, formated as a
// can we have [][] slices?

func createGrid(n int) [][]float64 {
  grid := make([][]float64, n)
  for i := 0; i < n; i++ {
    grid[i] = make([]float64, n)
    for y := 0; y < n; y++ {
      grid[i][y] = randVal()
    }
  }
  return grid
}

func randVal() float64 {
  return randIter(0)
}

func randIter(iter int) float64 {
  return (rand.Float64() * 2.0 - 1.0) * math.Pow(2, 0.8 * float64(iter))
}

func main() {
 fmt.Println("grid", createGrid(10))
}
