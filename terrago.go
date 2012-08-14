package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
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
	var wg sync.WaitGroup

	// copy over old values
	for y := 0; y < oldLen; y++ {
		for x := 0; x < oldLen; x++ {
			wg.Add(1)
			go expand(newGrid, grid, x, y, &wg)
		}
	}

	wg.Wait()

	// diamond step
	for y := 1; y < newLen; y += 2 {
		for x := 1; x < newLen; x += 2 {
			wg.Add(1)
			go diamond(newGrid, x, y, n, &wg)
		}
	}

	wg.Wait()

	// square step
	for y := 0; y < newLen; y += 2 {
		for x := 1; x < newLen; x += 2 {
			wg.Add(1)
			go square(newGrid, x, y, n, &wg)
		}
	}
	for y := 1; y < newLen; y += 2 {
		for x := 0; x < newLen; x += 2 {
			wg.Add(1)
			go square(newGrid, x, y, n, &wg)
		}
	}

	wg.Wait()

	return newGrid
}

func expand(newGrid Grid, oldGrid Grid, x int, y int, wg *sync.WaitGroup) {
	newGrid[2*x][2*y] = oldGrid[x][y]
	wg.Done()
}

func diamond(grid Grid, x int, y int, n int, wg *sync.WaitGroup) {
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
	wg.Done()
}

func square(grid Grid, x int, y int, n int, wg *sync.WaitGroup) {
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
	wg.Done()
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
	//	prettyPrintCompare(initGrid(9))
	grid := initGrid(3)

	t0 := time.Now()
	for i := 1; i <= 8; i++ {
		grid = iterGrid(grid, i)
	}
	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))

	// ~ 1.6 secs for n=2 and 10 iterations
}
