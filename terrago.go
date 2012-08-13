package main

import (
  "fmt"
  "math"
  "math/rand"
)

func initGrid(n int) [][]float64 {
  grid := make([][]float64, n)
  for i := 0; i < n; i++ {
    grid[i] = make([]float64, n)
    for y := 0; y < n; y++ {
      grid[i][y] = randVal()
    }
  }
  return grid
}

func iterGrid(grid [][]float64) [][]float64 {
  // create new grid double the size
  // start by making 2x with moving all i and j to 2i and 2j
  oldLen := len(grid)
  newLen := oldLen * 2
  newGrid := initGrid(newLen) // TODO: randVal should take iteration count.
  for y := 0; y < oldLen; y++ {
    for x := 0; x < oldLen; x++ {
      newGrid[x * 2][y * 2] = grid[x][y]
    }
  }
  return newGrid
}

func randVal() float64 {
  return randIter(0)
}

func randIter(iter int) float64 {
  return (rand.Float64() * 2.0 - 1.0) * math.Pow(2, 0.8 * float64(iter))
}

func prettyPrint(grid [][]float64) {
  var sym string
  n := len(grid)
  for x := 0; x < n; x++ {
    for y := 0; y < n; y++ {
      switch cell := grid[x][y]; {
      case cell < -0.5: sym = "  "
      case -0.5 <= cell && cell < 0.5: sym = ". "
      case 0.5 <= cell && cell < 1.5: sym = "+ "
      case 1.5 <= cell: sym = "# "
      }
      fmt.Print(sym)
    }
    fmt.Println()
  }
}

func main() {
  prevGrid := initGrid(9)

  grid := iterGrid(prevGrid)

  prettyPrint(prevGrid)
  fmt.Println("-------------------------------------")
  prettyPrint(grid)
}

// a(k+1) [2i, 2j] = a(k) [i, j]
