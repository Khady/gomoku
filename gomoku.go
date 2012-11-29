package main

import (
	"errors"
	"fmt"
)

type Gomoku struct {
	board       []int
	gameType    bool
	endgameTake bool
	doubleThree bool
	playerTurn  int
	countTake   [2]int
}

func (p *Gomoku) verifLine(x, y, count, time, varx, vary int) int {
	if x+varx >= 0 && y+vary >= 0 && x+varx <= 18 && y+vary <= 18 &&
		p.board[x+varx+(y+vary)*19] == p.playerTurn && (p.endgameTake == false || p.verifNotTakable(x, y)) {
		if time == 4 {
			return count + 1
		} else {
			return p.verifLine(x+varx, y+vary, count+1, time+1, varx, vary)
		}
	}
	return count
}

func (p *Gomoku) victory(x, y int) bool {
	if p.victoryPion(x, y) || p.victoryPion(x + 1, y) || p.victoryPion(x + 2, y) ||
		p.victoryPion(x, y + 1) || p.victoryPion(x, y + 2) || p.victoryPion(x + 1, y + 1) ||
		p.victoryPion(x + 2, y + 2) || p.victoryPion(x - 1, y - 1) || p.victoryPion(x + 1, y - 1) ||
		p.victoryPion(x + 2, y - 2) || p.victoryPion(x - 1, y + 1) || p.victoryPion(x - 2, y + 2) {
		return true
	}
	return false
}

func (p *Gomoku) victoryPion(x, y int) bool {
	if p.countTake[p.playerTurn-1] == 0 || (p.verifLine(x, y, p.verifLine(x, y, 1, 1, -1, 0), 1, 1, 0) >= 5 ||
		p.verifLine(x, y, p.verifLine(x, y, 1, 1, 0, 1), 1, 0, -1) >= 5 ||
		p.verifLine(x, y, p.verifLine(x, y, 1, 1, -1, -1), 1, 1, 1) >= 5 ||
		p.verifLine(x, y, p.verifLine(x, y, 1, 1, -1, +1), 1, 1, -1) >= 5) && (p.endgameTake == false || p.verifNotTakable(x, y)) {
		return true
	}
	return false
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (p *Gomoku) verifEnemy(x, y, varx, vary, flag int) bool {
	if flag == 1 {
		x = x + varx
		y = y + vary
	}
	if varx > 0 {
		varx = 1
	} else if varx < 0 {
		varx = -1
	}
	if vary > 0 {
		vary = 1
	} else if vary < 0 {
		vary = -1
	}
	if p.board[x+varx+(y+vary)*19] == p.otherPlayer() || p.board[x-varx+(y-vary)*19] == p.otherPlayer() {
		return false
	}
	return true
}

func (p *Gomoku) verifThree(x, y, prof, varx1, vary1, varx2, vary2 int) int { // oubli de verification des empty + probleme avec l'exemple en cours.
	if x > 0 && y > 0 && x < 18 && y < 18 && x > 0-Min(varx1, varx2) && x < 18-Max(varx1, varx2) && y > 0-Min(vary1, vary2) && y < 18-Max(vary1, vary2) &&
		p.board[x+varx1+(y+vary1)*19] == p.playerTurn && p.board[x+varx2+(y+vary2)*19] == p.playerTurn &&
		p.verifEnemy(x, y, varx1, vary1, 1) && p.verifEnemy(x, y, varx2, vary2, 1) && p.verifEnemy(x, y, varx1, vary1, 0) && p.verifEnemy(x, y, varx2, vary2, 0) {
		if prof == 0 && (p.verifDoubleThree(x+varx1, y+vary1, 1) || p.verifDoubleThree(x-varx2, y+vary2, 1)){
			return 2
		}
		return 1
	}
	return 0
}

func (p *Gomoku) verifDoubleThree(x, y, prof int) bool {
	count := 0
	verif1 := p.verifThree(x, y, prof, 1, 0, -1, 0)
	verif2 := p.verifThree(x, y, prof, 1, 0, 2, 0)
	verif3 := p.verifThree(x, y, prof, -1, 0, -2, 0)
	verif4 := p.verifThree(x, y, prof, -2, 0, -3, 0)
	verif5 := p.verifThree(x, y, prof, 2, 0, 3, 0)
	verif6 := p.verifThree(x, y, prof, -1, 0, -3, 0)
	verif7 := p.verifThree(x, y, prof, 1, 0, 3, 0)
	if verif1 == 2 || verif2 == 2 || verif3 == 2 || verif4 == 2 || verif5 == 2 || verif6 == 2 || verif7 == 2 {
		return true
	} else if verif1 == 1 || verif2 == 1 || verif3 == 1 || verif4 == 1 || verif5 == 1 || verif6 == 1 || verif7 == 1 {
		count += 1
	}
	verif1 = p.verifThree(x, y, prof, 0, 1, 0, -1)
	verif2 = p.verifThree(x, y, prof, 0, 1, 0, 2)
	verif3 = p.verifThree(x, y, prof, 0, -1, 0, -2)
	verif4 = p.verifThree(x, y, prof, 0, -2, 0, -3)
	verif5 = p.verifThree(x, y, prof, 0, 2, 0, 3)
	verif6 = p.verifThree(x, y, prof, 0, -1, 0, -3)
	verif7 = p.verifThree(x, y, prof, 0, 1, 0, 3)
	if verif1 == 2 || verif2 == 2 || verif3 == 2 || verif4 == 2 || verif5 == 2 || verif6 == 2 || verif7 == 2 {
		return true
	} else if verif1 == 1 || verif2 == 1 || verif3 == 1 || verif4 == 1 || verif5 == 1 || verif6 == 1 || verif7 == 1 {
		count += 1
	}
	if count > 1 {
		return true
	}
	verif1 = p.verifThree(x, y, prof, 1, 1, -1, -1)
	verif2 = p.verifThree(x, y, prof, 1, 1, 2, 2)
	verif3 = p.verifThree(x, y, prof, -1, -1, -2, -2)
	verif4 = p.verifThree(x, y, prof, -2, -2, -3, -3)
	verif5 = p.verifThree(x, y, prof, 2, 2, 3, 3)
	verif6 = p.verifThree(x, y, prof, -1, -1, -3, -3)
	verif7 = p.verifThree(x, y, prof, 1, 1, 3, 3)
	if verif1 == 2 || verif2 == 2 || verif3 == 2 || verif4 == 2 || verif5 == 2 || verif6 == 2 || verif7 == 2 {
		return true
	} else if verif1 == 1 || verif2 == 1 || verif3 == 1 || verif4 == 1 || verif5 == 1 || verif6 == 1 || verif7 == 1 {
		count += 1
	}
	if count > 1 {
		return true
	}
	verif1 = p.verifThree(x, y, prof, 1, -1, -1, 1)
	verif2 = p.verifThree(x, y, prof, 1, -1, 2, -2)
	verif3 = p.verifThree(x, y, prof, -1, 1, -2, 2)
	verif4 = p.verifThree(x, y, prof, -2, 2, -3, 3)
	verif5 = p.verifThree(x, y, prof, 2, -2, 3, -3)
	verif6 = p.verifThree(x, y, prof, 1, -1, 3, -3)
	verif7 = p.verifThree(x, y, prof, -1, 1, -3, 3)
	if verif1 == 2 || verif2 == 2 || verif3 == 2 || verif4 == 2 || verif5 == 2 || verif6 == 2 || verif7 == 2 {
		return true
	} else if verif1 == 1 || verif2 == 1 || verif3 == 1 || verif4 == 1 || verif5 == 1 || verif6 == 1 || verif7 == 1 {
		count += 1
	}
	if count > 1 {
		return true
	}
	return false
}

func (p *Gomoku) otherPlayer() int {
	if p.playerTurn == 1 {
		return 2
	}
	return 1
}

func (p *Gomoku) verifNotTakable(x, y int) bool {
	if x <= 16 && x >= 1 && p.board[x+1+y*19] == p.playerTurn &&
		(p.board[x+2+y*19] == 0 || p.board[x-1+y*19] == 0) &&
		(p.board[x+2+y*19] == p.otherPlayer() || p.board[x-1+y*19] == p.otherPlayer()) {
		return false
	}
	if x <= 17 && x >= 2 && p.board[x-1+y*19] == p.playerTurn &&
		(p.board[x-2+y*19] == 0 || p.board[x+1+y*19] == 0) &&
		(p.board[x-2+y*19] == p.otherPlayer() || p.board[x+1+y*19] == p.otherPlayer()) {
		return false
	}
	if y <= 17 && y >= 2 && p.board[x+(y-1)*19] == p.playerTurn &&
		(p.board[x+(y-2)*19] == 0 || p.board[x+(y+1)*19] == 0) &&
		(p.board[x+(y-2)*19] == p.otherPlayer() || p.board[x+(y+1)*19] == p.otherPlayer()) {
		return false
	}
	if y <= 16 && y >= 1 && p.board[x+(y+1)*19] == p.playerTurn &&
		(p.board[x+(y+2)*19] == p.otherPlayer() || p.board[x+(y-1)*19] == p.otherPlayer()) &&
		(p.board[x+(y+2)*19] == 0 || p.board[x+(y-1)*19] == 0) {
		return false
	}
	if y <= 17 && y >= 2 && x >= 2 && x <= 17 && p.board[x-1+(y-1)*19] == p.playerTurn &&
		(p.board[x-2+(y-2)*19] == 0 || p.board[x+1+(y+1)*19] == 0) &&
		(p.board[x-2+(y-2)*19] == p.otherPlayer() || p.board[x+1+(y+1)*19] == p.otherPlayer()) {
		return false
	}
	if y <= 16 && y >= 2 && x >= 2 && x <= 17 && p.board[x-1+(y+1)*19] == p.playerTurn &&
		(p.board[x-2+(y+2)*19] == 0 || p.board[x+1+(y-1)*19] == 0) &&
		(p.board[x-2+(y+2)*19] == p.otherPlayer() || p.board[x+1+(y-1)*19] == p.otherPlayer()) {
		return false
	}
	if y <= 16 && y >= 2 && x <= 16 && x >= 1 && p.board[x+1+(y-1)*19] == p.playerTurn &&
		(p.board[x+2+(y-2)*19] == 0 || p.board[x-1+(y+1)*19] == 0) &&
		(p.board[x+2+(y-2)*19] == p.otherPlayer() || p.board[x-1+(y+1)*19] == p.otherPlayer()) {
		return false
	}
	if y <= 16 && y >= 1 && x <= 16 && y >= 1 && p.board[x+1+(y+1)*19] == p.playerTurn &&
		(p.board[x+2+(y+2)*19] == 0 || p.board[x-1+(y-1)*19] == 0) &&
		(p.board[x+2+(y+2)*19] == p.otherPlayer() || p.board[x-1+(y-1)*19] == p.otherPlayer()) {
		return false
	}
	return true
}

func (p *Gomoku) Play(x, y int) (int, error) {

	if p.board[x+y*19] == 0 {
		p.board[x+y*19] = p.playerTurn
	} else {
		return 0, errors.New("move not valid")
	}
	if p.board[x+y*19] != 0 && (p.doubleThree == true && p.verifDoubleThree(x, y, 0) == true) {
		p.board[x+y*19] = 0
		return 0, errors.New("move not valid")
	}
	p.prise(x, y)
	if p.victory(x, y) {
		return p.playerTurn, nil
	}
	p.changePlayerTurn()
	return 0, nil
}

func (p *Gomoku) changePlayerTurn() {
	if p.playerTurn == 1 {
		p.playerTurn = 2
	} else {
		p.playerTurn = 1
	}
}

func (p *Gomoku) prise(x, y int) {
	if x <= 15 && p.board[x+1+y*19] != p.playerTurn &&
		p.board[x+2+y*19] != p.playerTurn && p.board[x+3+y*19] == p.playerTurn {
		p.board[x+2+y*19] = 0
		p.board[x+1+y*19] = 0
		p.countTake[p.playerTurn-1] -= 2
	}
	if x >= 3 && p.board[x-1+y*19] != p.playerTurn &&
		p.board[x-2+y*19] != p.playerTurn && p.board[x-3+y*19] == p.playerTurn {
		p.board[x-2+y*19] = 0
		p.board[x-1+y*19] = 0
		p.countTake[p.playerTurn-1] -= 2
	}
	if y >= 3 && p.board[x+(y-1)*19] != p.playerTurn &&
		p.board[x+(y-2)*19] != p.playerTurn && p.board[x+(y-3)*19] == p.playerTurn {
		p.board[x+(y-1)*19] = 0
		p.board[x+(y-2)*19] = 0
		p.countTake[p.playerTurn-1] -= 2
	}
	if y <= 15 && p.board[x+(y+1)*19] != p.playerTurn &&
		p.board[x+(y+2)*19] != p.playerTurn && p.board[x+(y+3)*19] == p.playerTurn {
		p.board[x+(y+1)*19] = 0
		p.board[x+(y+2)*19] = 0
		p.countTake[p.playerTurn-1] -= 2
	}
	if y >= 3 && x >= 3 && p.board[x-1+(y-1)*19] != p.playerTurn &&
		p.board[x-2+(y-2)*19] != p.playerTurn && p.board[x-3+(y-3)*19] == p.playerTurn {
		p.board[x-1+(y-1)*19] = 0
		p.board[x-2+(y-2)*19] = 0
		p.countTake[p.playerTurn-1] -= 2
	}
	if y <= 15 && x >= 3 && p.board[x-1+(y+1)*19] != p.playerTurn &&
		p.board[x-2+(y+2)*19] != p.playerTurn && p.board[x-3+(y+3)*19] == p.playerTurn {
		p.board[x-1+(y+1)*19] = 0
		p.board[x-2+(y+2)*19] = 0
		p.countTake[p.playerTurn-1] -= 2
	}
	if y >= 3 && x <= 15 && p.board[x+1+(y-1)*19] != p.playerTurn &&
		p.board[x+2+(y-2)*19] != p.playerTurn && p.board[x+3+(y-3)*19] == p.playerTurn {
		p.board[x+1+(y-1)*19] = 0
		p.board[x+2+(y-2)*19] = 0
		p.countTake[p.playerTurn-1] -= 2
	}
	if y <= 15 && x <= 15 && p.board[x+1+(y+1)*19] != p.playerTurn &&
		p.board[x+2+(y+2)*19] != p.playerTurn && p.board[x+3+(y+3)*19] == p.playerTurn {
		p.board[x+1+(y+1)*19] = 0
		p.board[x+2+(y+2)*19] = 0
		p.countTake[p.playerTurn-1] -= 2
	}
}

func (p *Gomoku) Debug_aff() {
	for i := 0; i < 19; i++ {
		for n := 0; n < 19; n++ {
			fmt.Print(p.board[i*19+n])
		}
		fmt.Println()
	}
	fmt.Println(p.countTake[0], p.countTake[1])
}

// func main() {
// 	game := Gomoku{make([]int, 361), true, true, true, 1, [2]int{10, 10}};
// 	game.play(1, 0)
// 	game.play(2, 0)
// 	game.play(1, 1)
// 	game.play(6, 2)
// 	game.play(1, 2)
// 	game.play(5, 0)
// 	game.play(1, 3)
// 	game.play(7, 0)
// 	game.play(2, 2)
// 	game.play(8, 0)
// 	if (vic := game.play(1, 4)) != 0 {
// 		fmt.Println("player ",  vic, "win")
// 	}
// 	game.play(9, 0)
// 	if game.play(12, 5) != 0 {
// 		fmt.Println("win")
// 	}
// 	game.debug_aff()
// }
