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
	"time"
)

type Grid [][]float64

const NCPU = 4

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func randIter(iter int) float64 {
	return (rand.Float64()*2.0 - 1.0) * math.Pow(2, 0.8*float64(iter))
}

// Grid functions

// Creates a new grid of size n, with random heights.
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

	// diamond step
	for i := 0; i < NCPU; i++ {
		go diamondSegment(newGrid, i, n, c)
	}

	// wait for all calculations to finish
	for i := 0; i < NCPU; i++ {
		<-c
	}

	// square step
	for i := 0; i < NCPU; i++ {
		go squareSegment(newGrid, i, n, c)
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

func diamondSegment(grid Grid, offset int, n int, c chan int) {
	length := len(grid)
	for y := 1; y < length; y += 2 {
		for x := 1 + 2*offset; x < length; x += 2 * NCPU {
			diamond(grid, x, y, n)
		}
	}

	// send segment finished to channel
	c <- 1
}

func squareSegment(grid Grid, offset int, n int, c chan int) {
	length := len(grid)
	for y := 0; y < length; y += 2 {
		for x := 1 + 2*offset; x < length; x += 2 * NCPU {
			square(grid, x, y, n)
		}
	}
	for y := 1; y < length; y += 2 {
		for x := 0 + 2*offset; x < length; x += 2 * NCPU {
			square(grid, x, y, n)
		}
	}

	// send segment finished to channel
	c <- 1
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

func prettyPrintCompare(grid Grid, c chan int) {
	prevGrid := grid
	newGrid := iterGrid(prevGrid, 1, c)
	prettyPrint(prevGrid)
	fmt.Println("-------------------------------------")
	prettyPrint(newGrid)
}

func render2D(grid Grid) {
	img := image.NewNRGBA(image.Rect(0, 0, len(grid), len(grid[0])))
	bounds := img.Bounds()
	min := takeMin(grid)
	max := takeMax(grid)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		realY := y - bounds.Min.Y
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			realX := x - bounds.Min.X
			img.SetNRGBA(x, y, calcColor(grid, min, max, grid[realX][realY]))
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
func calcColor(grid Grid, min float64, max float64, val float64) color.NRGBA {
	var r, g, b uint8
	diff := max - min
	normalized := float64(val-min) / float64(diff)

	switch {
	case normalized < .1:
		b = 255
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
	grid := initGrid(3)

	t0 := time.Now()
	for i := 1; i <= 10; i++ {
		grid = iterGrid(grid, i, c)
	}
	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))

	render2D(grid)
	fmt.Println("Created PNG.")

	// ~ 1.6 secs for n=2 and 10 iterations
}
