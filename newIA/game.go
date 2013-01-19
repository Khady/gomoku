package main

import (
	"errors"
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
	playerTurn  uint64
	stonesTaken [2]byte
	DoubleThreeCheckFuncs map[uint8]func (int, int) bool
}

const (
	EMPTY = 0
	WHITE = 1
	BLACK = 2

	UP        = 0x1
	LEFT      = 0x2
	DOWN      = 0x4
	RIGHT     = 0x8
	LEFTUP    = 0x10
	LEFTDOWN  = 0x20
	RIGHTDOWN = 0x40
	RIGHTUP   = 0x80
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
	dumpVerticalLines(game.VLines)
	dumpDiagonaleLines(game.DRLLines)
	dumpDiagonaleLines(game.DLRLines)
	dumpHorizontalLines(game.HLines)

}


func (game *GomokuBoard) clearBoard() {
	for i := 0; i < 37; i++ {
		if i < 19 {
			game.HLines[i] &= 0
			game.VLines[i] &= 0
		}
		game.DRLLines[i] &= 0
		game.DLRLines[i] &= 0
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

// Il semblerait qu'en Go, ce qu'on fait dans d'autre langages avec l'operateur NOT (~) se fait avc '^'...
func (game *GomokuBoard) emptyCase(x, y int) {
	color := game.whatIsOnCase(x, y)

	game.HLines[y] &= ^(color << (36 - uint64(2 * x)))

	game.VLines[x] &= ^(color << (36 - uint64(2 * x)))

	// Setting case on the diagonales lines (left to right and right to left)
	indexDRL := x + y
	if indexDRL <= 18 {
		game.DRLLines[indexDRL] &= ^(color << (36 - uint64(2 * y)))
	} else {
		game.DRLLines[indexDRL] &= ^(color << uint64(2 * x))
	}

	indexDLR := 18 - (y - x)
	if indexDLR >= 18 {
		game.DLRLines[indexDLR] &= (color << (36 - uint64(2 * y)))
	} else {
		game.DLRLines[indexDLR] |= (color << (36 - uint64(2 * x)))
	}

}

func (game *GomokuBoard) moveIsDTDiagonaleLRSafe(x, y int) bool {
	return true
}

func (game *GomokuBoard) moveIsDTDiagonaleRLSafe(x, y int) bool {
	return true

}

func (game *GomokuBoard) moveIsDTVerticallySafe(x, y int) bool {
	return true

}

func (game *GomokuBoard) checkDoubleThreeOnTheLeft(x, y int) bool {
	return (x >= 3 &&
		game.whatIsOnCase(x - 1, y) == game.playerTurn &&
		game.whatIsOnCase(x - 2, y) == game.playerTurn &&
		game.whatIsOnCase(x - 3, y) == EMPTY)
}

func (game *GomokuBoard) checkDoubleThreeOnTheRight(x, y int) bool {
	return (x <= 15 &&
		game.whatIsOnCase(x + 1, y) == game.playerTurn &&
		game.whatIsOnCase(x + 2, y) == game.playerTurn &&
		game.whatIsOnCase(x + 3, y) == EMPTY)
}

func (game *GomokuBoard) checkDoubleThreeOnTheDown(x, y int) bool {
	return (y <= 15 &&
		game.whatIsOnCase(x, y + 1) == game.playerTurn &&
		game.whatIsOnCase(x, y + 2) == game.playerTurn &&
		game.whatIsOnCase(x, y + 3) == EMPTY)
}

func (game *GomokuBoard) checkDoubleThreeOnTheUp(x, y int) bool {
	return (y >= 3 &&
		game.whatIsOnCase(x, y - 1) == game.playerTurn &&
		game.whatIsOnCase(x, y - 2) == game.playerTurn &&
		game.whatIsOnCase(x, y - 3) == EMPTY)
}

func (game *GomokuBoard) checkDoubleThreeOnTheLeftUp(x, y int) bool {
		return (x >= 3 &&
		y >= 3 &&
		game.whatIsOnCase(x - 1, y - 1) == game.playerTurn &&
		game.whatIsOnCase(x - 2, y - 2) == game.playerTurn &&
		game.whatIsOnCase(x - 3, y - 3) == EMPTY)
}

func (game *GomokuBoard) checkDoubleThreeOnTheLeftDown(x, y int) bool {
	return (x >= 3 && y <= 15 &&
		game.whatIsOnCase(x - 1, y + 1) == game.playerTurn &&
		game.whatIsOnCase(x - 2, y + 2) == game.playerTurn &&
		game.whatIsOnCase(x - 3, y + 3) == EMPTY)
}

func (game *GomokuBoard) checkDoubleThreeOnTheRightDown(x, y int) bool {
	return (x <= 15 && y <= 15 &&
		game.whatIsOnCase(x + 1, y + 1) == game.playerTurn &&
		game.whatIsOnCase(x + 2, y + 2) == game.playerTurn &&
		game.whatIsOnCase(x + 3, y + 3) == EMPTY)
}

func (game *GomokuBoard) checkDoubleThreeOnTheRightUp(x, y int) bool {
	return (x <= 15 && y >= 3 &&
		game.whatIsOnCase(x + 1, y - 1) == game.playerTurn &&
		game.whatIsOnCase(x + 2, y - 2) == game.playerTurn &&
		game.whatIsOnCase(x + 3, y - 3) == EMPTY)
}

func (game *GomokuBoard) initializeDTcheckFuncs() {
	game.DoubleThreeCheckFuncs = make(map[uint8]func (int, int) bool)


	game.DoubleThreeCheckFuncs[LEFT]	  = func (x, y int) bool { return game.checkDoubleThreeOnTheLeft(x, y) }
	game.DoubleThreeCheckFuncs[RIGHT]	  = func (x, y int) bool { return game.checkDoubleThreeOnTheRight(x, y) }
	game.DoubleThreeCheckFuncs[UP]	          = func (x, y int) bool { return game.checkDoubleThreeOnTheUp(x, y) }
	game.DoubleThreeCheckFuncs[DOWN]	  = func (x, y int) bool { return game.checkDoubleThreeOnTheDown(x, y) }
	game.DoubleThreeCheckFuncs[LEFTUP]	  = func (x, y int) bool { return game.checkDoubleThreeOnTheLeftUp(x, y) }
	game.DoubleThreeCheckFuncs[LEFTDOWN]      = func (x, y int) bool { return game.checkDoubleThreeOnTheLeftDown(x, y) }
	game.DoubleThreeCheckFuncs[RIGHTDOWN]     = func (x, y int) bool { return game.checkDoubleThreeOnTheRightDown(x, y) }
	game.DoubleThreeCheckFuncs[RIGHTUP]	  = func (x, y int) bool { return game.checkDoubleThreeOnTheRightUp(x, y) }
}

func (game *GomokuBoard) checkThreesAround(x, y int, checkFlags uint8) (Three bool) {

	Three = false
	for flag, checkFunc := range game.DoubleThreeCheckFuncs {
		if checkFlags & flag != 0 {
			Three = checkFunc(x, y)
			if Three == true {
				break
			}
		}
	}
	return
}

func (game *GomokuBoard) horizontalLeftDoubleThree(x, y int) bool {
	if x >= 3 && game.whatIsOnCase(x - 1, y) == game.playerTurn &&
		game.whatIsOnCase(x - 2, y) == game.playerTurn &&
		game.whatIsOnCase(x - 3, y) == EMPTY {
		// possible double free on the left : check nearly everything !
		return (game.checkThreesAround(x, y, UP | LEFTUP | LEFTDOWN | DOWN | RIGHTDOWN | RIGHTUP)     ||
			game.checkThreesAround(x - 1, y, UP | LEFTUP | LEFTDOWN | DOWN | RIGHTDOWN | RIGHTUP) ||
			game.checkThreesAround(x - 2, y, UP | LEFTUP | LEFTDOWN | DOWN | RIGHTDOWN | RIGHTUP))
	}

	return false
}

func (game *GomokuBoard) checkHorizontalFirstCase(x, y int) bool {
	if x >= 3 && game.horizontalLeftDoubleThree(x, y) {
		return false
	}
	
	return true
}

func (game *GomokuBoard) moveIsDTHorizontallySafe(x, y int) bool {
	return game.checkHorizontalFirstCase(x, y)
	
	return true
}

func (game *GomokuBoard) isMoveDoubleThreeSafe(x, y int) bool {
	return game.moveIsDTHorizontallySafe(x, y)
	/*&&
		game.moveIsDTVerticallySafe(x, y)   && 
		game.moveIsDTDiagonaleLRSafe(x, y)  &&
		game.moveIsDTDiagonaleRLSafe(x, y)) */
}

func (game *GomokuBoard) moveIsValid(x, y int) error {	
	if game.whatIsOnCase(x, y) != EMPTY {
		return errors.New("Invalid move (stone already on this position)")
	}
	if game.doubleThree == true {
		if game.isMoveDoubleThreeSafe(x, y) == false {
			return errors.New("Invalid move (creating a Double Three)")
		}
	}
	return nil
}

func (game *GomokuBoard) Play(x, y int) (int, [][2]int, error) {
	err := game.moveIsValid(x, y)
	if  err != nil {
		return 0, nil, err
	}

	return 42, nil, nil
}

func (game *GomokuBoard) setCase(x, y int, state uint64) {
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

	// Setting case on the diagonales lines (left to right and right to left)
	indexDRL := x + y
	if indexDRL <= 18 {
		game.DRLLines[indexDRL] |= (state << (36 - uint64(2 * y)))
	} else {
		game.DRLLines[indexDRL] |= (state << uint64(2 * x))
	}

	indexDLR := 18 - (y - x)
	if indexDLR >= 18 {
		game.DLRLines[indexDLR] |= (state << (36 - uint64(2 * y)))
	} else {
		game.DLRLines[indexDLR] |= (state << (36 - uint64(2 * x)))
	}
}

func main() {
	var game GomokuBoard
	
/*	i := 0
	for i < 19 {
		j := 0
		for j < 19 {
			game.setCase(i, j, 2)
			j++
		}
		i++
	}*/
//	game.dumpBoard()
//	game.clearBoard()
// UP | LEFTUP | LEFTDOWN | DOWN | LEFTDOWN | LEFTUP
	game.playerTurn = 1
	game.initializeDTcheckFuncs()
	game.setCase(9, 8, 1)
	game.setCase(10, 7, 1)
	game.setCase(8, 9, 1)
	game.setCase(9, 9, 1)
	game.setCase(7, 9, 1)
	game.setCase(1, 9, 2)
//	game.setCase(7, 9, 2)
	// NOTE POUR PLUS TARD : ON A TROUVE !
	// SI game.HLines[9] & (MASQUE << (36 - uint64(2 * 9))) == MASQUE && game.HLines[9] & (ENEMYMASK << (36 - uint64(2 * 9))) != ENEMYMASK
	// -> la voie est libre !
	game.dumpBoard()
	println("Case 9, 9 is", game.whatIsOnCase(9, 9), "is (10, 9) double three safe? ")
	tBefore := time.Now()
	res := game.isMoveDoubleThreeSafe(10, 9)
	fmt.Println(res, "(perf:", time.Since(tBefore), ")")
	fmt.Println("test", game.HLines[9] & (37 << (36 - uint64(2 * 9))), (37 << (36 - uint64(2 * 9))))
	
}