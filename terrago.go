package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

type Grid [][]float64

const NCPU = 4

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func randIter(iter int) float64 {
	return (rand.Float64()*2.0 - 1.0) * math.Pow(2, -0.8*float64(iter))
}

// Grid functions

// Creates a new grid of size n, with random heights.
func initGrid(n int) Grid {
	grid := make(Grid, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]float64, n)
		for y := 0; y < n; y++ {
			grid[i][y] = 0
		}
	}

	return grid
}

// n must be >= 1.
func iterGrid(grid Grid, n int, c chan int) Grid {
	oldLen := len(grid)
	newLen := (oldLen-1)*2 + 1 // must be of form 2**n + 1
	newGrid := initGrid(newLen)

	// copy over old values
	for y := 0; y < oldLen; y++ {
		for x := 0; x < oldLen; x++ {
			expand(newGrid, grid, x, y)
		}
	}

	// square step
	for i := 0; i < NCPU; i++ {
		go squareSegment(newGrid, i, n, c)
	}

	// wait for all calculations to finish
	for i := 0; i < NCPU; i++ {
		<-c
	}

	// diamond step
	for i := 0; i < NCPU; i++ {
		go diamondSegment(newGrid, i, n, c)
	}

	// wait for all calculations to finish
	for i := 0; i < NCPU; i++ {
		<-c
	}

	return newGrid
}

func expand(newGrid Grid, oldGrid Grid, x int, y int) {
	newGrid[2*x][2*y] = oldGrid[x][y]
}

// Performs the diamond step of the algorithm for all diamonds in the grid.
func diamondSegment(grid Grid, offset int, n int, c chan int) {
	length := len(grid)

	// (x,y) with offset (dx, dy) pairs that form a diamond
	dxList := []int{-1, 1, 0, 0}
	dyList := []int{0, 0, -1, 1}

	// refactor similarites of these for loops
	for y := 0; y < length; y += 2 {
		for x := 1 + 2*offset; x < length; x += 2 * NCPU {
			calcCenter(grid, dxList, dyList, x, y, n)
		}
	}
	for y := 1; y < length; y += 2 {
		for x := 0 + 2*offset; x < length; x += 2 * NCPU {
			calcCenter(grid, dxList, dyList, x, y, n)
		}
	}

	// send segment finished to channel
	c <- 1
}

// Performs the square step of the algorithm for all squares in the grid.
func squareSegment(grid Grid, offset int, n int, c chan int) {
	length := len(grid)

	// dx, dy pairs that form a square
	dxList := []int{-1, 1, -1, 1}
	dyList := []int{1, -1, -1, 1}

	for y := 1; y < length; y += 2 {
		for x := 1 + 2*offset; x < length; x += 2 * NCPU {
			calcCenter(grid, dxList, dyList, x, y, n)
		}
	}

	// send segment finished to channel
	c <- 1
}

// Takes the average of points and updates center point
// (dxList[i], dyList[i]), i=0..3 are points surrounding center
func calcCenter(grid Grid, dxList, dyList []int, x, y, n int) {
	var sum, num float64
	var length int = len(grid)

	// we sum corners around (x,y) to get average height of that area
	for i := range dxList {
		dx, dy := dxList[i], dyList[i]
		if x+dx >= 0 && x+dx < length && y+dy >= 0 && y+dy < length {
			sum += grid[x+dx][y+dy]
			num++
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
			case -0.5 <= cell && cell < 0.1:
				sym = ". "
			case 0.1 <= cell && cell < .2:
				sym = "+ "
			case .2 <= cell:
				sym = "# "
			}
			fmt.Print(sym)
		}
		fmt.Println()
	}
}

func printHeights(grid Grid) {
	var sym string
	n := len(grid)
	n = int(math.Min(float64(n), 10))
	for x := 0; x < n; x++ {
		for y := 0; y < n; y++ {
			sym = strconv.FormatFloat(grid[x][y], 'f', 8, 64)
			sym = sym[0:7] + " "
			fmt.Print(sym)
		}
		fmt.Println()
	}
}
func prettyPrintCompare(grid Grid, c chan int) {
	prevGrid := grid
	newGrid := iterGrid(prevGrid, 1, c)
	prettyPrint(prevGrid)
	fmt.Println("-------------------------------------")
	prettyPrint(newGrid)
}

// Renders a 2D image of a grid.
func render2D(grid Grid) {
	img := image.NewNRGBA(image.Rect(0, 0, len(grid), len(grid[0])))
	bounds := img.Bounds()
	min := takeMin(grid)
	max := takeMax(grid)
	// image's bounds doesn't have to start at (0,0)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		realY := y - bounds.Min.Y
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			realX := x - bounds.Min.X
			img.SetNRGBA(x, y, calcColor(grid[realX][realY], min, max))
		}
	}

	// create PNG
	w, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	err = png.Encode(w, img)
	if err != nil {
		log.Fatal(err)
	}
}

// Find minimum value in grid.
func takeMin(grid Grid) float64 {
	min := grid[0][0]
	for y := 0; y < len(grid[0]); y++ {
		for x := 0; x < len(grid); x++ {
			if grid[x][y] < min {
				min = grid[x][y]
			}
		}
	}
	return min
}

// Find maximum value in grid.
func takeMax(grid Grid) float64 {
	max := grid[0][0]
	for y := 0; y < len(grid[0]); y++ {
		for x := 0; x < len(grid); x++ {
			if grid[x][y] > max {
				max = grid[x][y]
			}
		}
	}
	return max
}

// Calculates color for a value in the grid by normalizing it.
func calcColor(val float64, min float64, max float64) color.NRGBA {
	var r, g, b uint8

	// TODO should not be in this function
	delta := max - min
	normalized := (val - min) / delta // we want (0,1)

	//exponorm := math.Pow(normalized, 0.7) // is this ok?

	switch {
	case normalized < .2:
		b = uint8(255 - normalized*255/.4)
	default:
		g = uint8(normalized * 255)
	}

	return color.NRGBA{R: r, G: g, B: b, A: 255}
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	c := make(chan int, NCPU)
	//	prettyPrintCompare(initGrid(9), c)
	rand.Seed(time.Now().UnixNano())

	grid := initGrid(2)
	printHeights(grid)
	println()
	println()
	t0 := time.Now()
	for i := 1; i <= 11; i++ {
		grid = iterGrid(grid, i, c)
	}
	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))

	render2D(grid)
	fmt.Println("Created PNG.")

	printHeights(grid)

	//prettyPrint(grid)

	// ~ 1.6 secs for n=2 and 10 iterations
}
