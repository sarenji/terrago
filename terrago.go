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

// n must be >= 1.
func iterGrid(grid [][]float64, n int) [][]float64 {
  oldLen := len(grid)
  newLen := oldLen * 2
  newGrid := initGrid(newLen) // TODO: randVal should take iteration count.

  gridEach(grid, evenEven(newGrid, grid, n))
  gridEach(grid, evenOdd(newGrid, grid, n))
  return newGrid
}

func evenEven(newGrid [][]float64, grid [][]float64, n int) func(float64, int, 
  int) {
  return func(cell float64, x int, y int) {
    newGrid[2 * x][2 * y] = cell
  }
}

func evenOdd(newGrid [][]float64, grid [][]float64, n int) func(float64, int, 
  int) {
  return func(cell float64, x int, y int) {
    indices := [][]int{
      []int{x, y},
      []int{x, y + 1},
      []int{x + 1, y},
      []int{x + 1, y + 1},
    }
    newGrid[2 * x + 1][2 * y + 1] = gridAverage(grid, indices) + randIter(n)
  }
}

func gridEach(grid [][]float64, callback func(float64, int, int)) [][]float64 {
  length := len(grid)
  for y := 0; y < length; y++ {
    for x := 0; x < length; x++ {
      callback(grid[x][y], x, y)
    }
  }
  return grid
}

func gridAverage(grid [][]float64, indices [][]int) float64 {
  var sum float64 = 0.0
  length := len(grid)
  indicesLen := len(indices)
  cells := make([]float64, 0, 4)
  for i := 0; i < indicesLen; i++ {
    x := indices[i][0]
    y := indices[i][1]
    if x >= 0 && y >= 0 && x < length && y < length {
      cells = append(cells, grid[x][y])
    }
  }
  for _, cell := range cells {
    sum += cell
  }
  return sum / float64(len(cells))
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

  grid := iterGrid(prevGrid, 1)

  prettyPrint(prevGrid)
  fmt.Println("-------------------------------------")
  prettyPrint(grid)
}

// a(k+1) [2i, 2j] = a(k) [i, j]
