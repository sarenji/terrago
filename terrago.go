package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Grid [][]float64

func randIter(iter int) float64 {
	return (rand.Float64()*2.0 - 1.0) * math.Pow(2, 0.8*float64(iter))
}

// Grid functions

func initGrid(n int) Grid {
	grid := make(Grid, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]float64, n)
		for y := 0; y < n; y++ {
			grid[i][y] = randIter(0)
		}
	}
	return grid
}

// n must be >= 1.
func iterGrid(grid Grid, n int) Grid {
	oldLen := len(grid)
	newLen := (oldLen-1)*2 + 1 // must be of form 2**n + 1
	newGrid := initGrid(newLen)

	// copy over old values
	for y := 0; y < oldLen; y++ {
		for x := 0; x < oldLen; x++ {
			newGrid[2*x][2*y] = grid[x][y]
		}
	}

	// diamond step
	for y := 1; y < newLen; y += 2 {
		for x := 1; x < newLen; x += 2 {
			diamond(newGrid, x, y, n)
		}
	}

	// square step
	for y := 0; y < newLen; y += 2 {
		for x := 1; x < newLen; x += 2 {
			square(newGrid, x, y, n)
		}
	}
	for y := 1; y < newLen; y += 2 {
		for x := 0; x < newLen; x += 2 {
			square(newGrid, x, y, n)
		}
	}

	return newGrid
}

func diamond(grid Grid, x int, y int, n int) {
	var sum, num float64
	var length int = len(grid)
	if x-1 >= 0 {
		sum, num = sum+grid[x-1][y], num+1
	}
	if x+1 < length {
		sum, num = sum+grid[x+1][y], num+1
	}
	if y-1 >= 0 {
		sum, num = sum+grid[x][y-1], num+1
	}
	if y+1 < length {
		sum, num = sum+grid[x][y+1], num+1
	}
	grid[x][y] = sum/num + randIter(n)
}

func square(grid Grid, x int, y int, n int) {
	var sum, num float64
	var length int = len(grid)
	if x-1 >= 0 {
		if y-1 >= 0 {
			sum, num = sum+grid[x-1][y-1], num+1
		}
		if y+1 < length {
			sum, num = sum+grid[x-1][y+1], num+1
		}
	}
	if x+1 < length {
		if y-1 >= 0 {
			sum, num = sum+grid[x+1][y-1], num+1
		}
		if y+1 < length {
			sum, num = sum+grid[x+1][y+1], num+1
		}
	}
	grid[x][y] = sum/num + randIter(n)
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
