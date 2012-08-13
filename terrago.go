package main

import (
  "fmt"
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
      grid[i][y] = rand.Float64()
    }
  }
  return grid
}

func main() {
 fmt.Println("grid", createGrid(10))
}
