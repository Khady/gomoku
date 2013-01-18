package main

import (
	"fmt"
	"time"
)

type GomokuBoard struct {
	HLines      [19]uint64
	VLines      [19]uint64
	DLRLines    [37]uint64
	DRLLines    [37]uint64
	endgameTake bool
	doubleThree bool
}

const (
	EMPTY = 0
	WHITE = 1
	BLACK = 2
)

func dumpBit(val uint64, i uint64) {
	if (val & (1 << i)) != 0 {
		print("1")
	} else { print("0") }
}

func dumpHorizontalLines(lines [19]uint64) {
	for i := 0; i < 19; i++ {
		for j := uint64(37); j >= 0; j-- {
			dumpBit(lines[i], j)
			if j % 2 == 0 {
				print(" ")
			}
			if j == 0 { break }
		}
		print(" ", i, "\n")
	}
	print("\n")
}

func dumpDiagonaleLines(lines [37]uint64) {
	for i := 0; i < 37; i++ {
		for j := uint64(37); j >= 0; j-- {
			dumpBit(lines[i], j)
			if j % 2 == 0 {
				print(" ")
			}
			if j == 0 { break }
		}
		print(" ", i, "\n")
	}
	print("\n")
}


func dumpVerticalLines(lines [19]uint64) {
	for j := uint64(37); j >= 0; j-- {
		for i := 0; i < 19; i++ {
			dumpBit(lines[i], j)
			print(" ")
		}
		print("\n")
		if j % 2 == 0 {
			print("\n")
		}
		if j == 0 { break }
	}
	print("\n")
}


func (game *GomokuBoard) dumpBoard() {
//	dumpVerticalLines(game.VLines)
	dumpDiagonaleLines(game.DLRLines)
	//dumpDRLLines(game.DRLLines)
	//dumpHorizontalLines(game.HLines)

}

func (game *GomokuBoard) setCaseHorizontal(line *uint64, x int, state uint64) {
	if x == 18 {
		*line |= state
	} else {
		
		*line |= (state << (36 - uint64(2 * x)))
	}
}

func (game *GomokuBoard) setCaseVertical(line *uint64, y int, state uint64) {
	if y == 18 {
		*line |= state
	} else {
		
		*line |= (state << (36 - uint64(2 * y)))
	}
}


func (game *GomokuBoard) setCaseDiagonaleTopBottom(line *uint64, x, y, index int, state uint64) {
	if index <= 18 {
		*line |= (state << (36 - uint64(2 * y)))
	} else {
		*line |= (state << uint64(2 * x))
	}
}

func (game *GomokuBoard) setCaseDiagonaleBottomTop(line *uint64, x, y, index int, state uint64) {
	if index <= 18 {
		*line |= (state << (36 - uint64(2 * y)))
	} else {
		*line |= (state << (36 - uint64(2 * x)))
	}
}

func (game *GomokuBoard) whatIsOnCase(x, y int) uint64 {
	if (game.HLines[y] & (1 << (36 - uint64(2 * x)))) != 0 {
		return WHITE
	}
	if (game.HLines[y] & (2 << (36 - uint64(2 * x)))) != 0 {
		return BLACK
	}
	return EMPTY
}

func (game *GomokuBoard) emptyCaseHorizontal(line *uint64, x, y int) {
	*line &= ^(game.whatIsOnCase(x, y) << (36 - uint64(2 * x)))
}

func (game *GomokuBoard) emptyCase(x, y int) {
	game.emptyCaseHorizontal(&(game.HLines[y]), x, y)
	//game.emptyCaseVertical(&(game.VLines[x]), y)
	//game.emptyCaseDiagonaleTopBottom(&(game.DRLLines[index]), x, y, index)
	/*if x >= y {
		game.emptyCaseDiagonaleBottomTop(&(game.DLRLines[18 - (x - y)]), x, y, 18 - (x - y))
	} else {
		game.emptyCaseDiagonaleBottomTop(&(game.DLRLines[18 + (y - x)]), x, y, 18 + (y - x))
	}*/
}

func (game *GomokuBoard) setCase(x, y int, state uint64) {
	var index int
	
	// Setting case on horizontal lines
	if x == 18 {
		game.HLines[y] |= state
	} else {
		game.HLines[y] |= (state << (36 - uint64(2 * x)))
	}

	// Setting case on vertical lines
	if y == 18 {
		game.VLines[x] |= state
	} else {
		
		game.VLines[x] |= (state << (36 - uint64(2 * y)))
	}

	// Setting case on "right to left" diagonale lines
	index = x + y
	if index <= 18 {
		game.DRLLines[index] |= (state << (36 - uint64(2 * y)))
	} else {
		game.DRLLines[index] |= (state << uint64(2 * x))
	}

	// Setting case on "left to right" diagonale lines
	if x >= y {
		index = 18 - (x - y)
	} else {
		index = 18 + (y - x)
	}
	if index <= 18 {
		game.DLRLines[index] |= (state << (36 - uint64(2 * y)))
	} else {
		game.DLRLines[index] |= (state << (36 - uint64(2 * x)))
	}
}

func main() {
	var game GomokuBoard
	
	i := 0
	tBefore := time.Now()
	for i < 19 {
		j := 0
		for j < 19 {
			game.setCase(i, j, 2)
			j++
		}
		i++
	}
	fmt.Println("Le coup a pris", time.Since(tBefore))
//	game.setCase(2, 3, 1)
	game.emptyCase(1, 2)
	game.dumpBoard()
	println("Case 1, 2 is ", game.whatIsOnCase(1, 2))
}