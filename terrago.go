package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Grid [][]float64

// Utils

func randVal() float64 {
	return randIter(0)
}

func randIter(iter int) float64 {
	return (rand.Float64()*2.0 - 1.0) * math.Pow(2, 0.8*float64(iter))
}

// Grid functions

func gridAverage(grid Grid, indices [][]int) float64 {
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

func initGrid(n int) Grid {
	grid := make(Grid, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]float64, n)
		for y := 0; y < n; y++ {
			grid[i][y] = randVal()
		}
	}
	return grid
}

// n must be >= 1.
func iterGrid(grid Grid, n int) Grid {
	oldLen := len(grid)
	newLen := (oldLen-1)*2 + 1
	newGrid := initGrid(newLen) // TODO: randVal should take iteration count.

	// TODO array out of bounds - temp fix by +1 and -1
	for y := 1; y < oldLen-1; y++ {
		for x := 1; x < oldLen-1; x++ {
			switch {
			case x%2 == 0 && y%2 == 0: // Even|Even
				evenEven(newGrid, grid, x, y, n)
			case x%2 == 1 && y%2 == 1: // Odd|Odd
				oddOdd(newGrid, grid, x, y, n)
			case x%2 == 1 && y%2 == 0: // Odd|Even
				oddEven(newGrid, grid, x, y, n)
			case x%2 == 0 && y%2 == 1: // Even|Odd
				evenOdd(newGrid, grid, x, y, n)
			}
		}
	}

	return newGrid
}

// Even/Odd functions

// (1)
func evenEven(newGrid Grid, grid Grid, x int, y int, n int) {
	newGrid[2*x][2*y] = grid[x][y]
}

// (2)
func oddOdd(newGrid Grid, grid Grid, x int, y int, n int) {
	indices := [][]int{
		[]int{x, y},
		[]int{x, y + 1},
		[]int{x + 1, y},
		[]int{x + 1, y + 1},
	}
	newGrid[2*x+1][2*y+1] = gridAverage(grid, indices) + randIter(n)
}

// (3)
func evenOdd(newGrid Grid, grid Grid, x int, y int, n int) {
	sum := (3.0*grid[x][y] + 3.0*grid[x][y+1]) / 8.0
	sum += (grid[x-1][y] + grid[x-1][y+1] + grid[x+1][y] + grid[x+1][y+1]) / 16
	newGrid[2*x][2*y+1] = sum
}

// (4)
func oddEven(newGrid Grid, grid Grid, x int, y int, n int) {
	sum := (3.0*grid[x][y] + 3.0*grid[x+1][y]) / 8.0
	sum += (grid[x][y-1] + grid[x+1][y-1] + grid[x][y+1] + grid[x+1][y+1]) / 16
	newGrid[2*x+1][2*y] = sum
}

// Print functions

func prettyPrint(grid Grid) {
	var sym string
	n := len(grid)
	for x := 0; x < n; x++ {
		for y := 0; y < n; y++ {
			switch cell := grid[x][y]; {
			case cell < -0.5:
				sym = "  "
			case -0.5 <= cell && cell < 0.5:
				sym = ". "
			case 0.5 <= cell && cell < 1.5:
				sym = "+ "
			case 1.5 <= cell:
				sym = "# "
			}
			fmt.Print(sym)
		}
		fmt.Println()
	}
}

func prettyPrintCompare(grid Grid) {
	prevGrid := grid
	newGrid := iterGrid(prevGrid, 1)
	prettyPrint(prevGrid)
	fmt.Println("-------------------------------------")
	prettyPrint(newGrid)
}

func main() {
	prettyPrintCompare(initGrid(9))
}
