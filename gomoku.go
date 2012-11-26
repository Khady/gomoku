package main

import "fmt"

type Gomoku struct {
	board []int
	gameType bool
	endgameTake bool
	doublethree bool
	playerTurn int
	countTake[2] int
}

func (p *Gomoku)verifLine(x, y, count, time, varx, vary int) int {
	if x + varx >= 0 && y + vary >= 0 && x + varx <= 18 && y + vary <= 18 &&
		p.board[x + varx + (y + vary) * 19] == p.playerTurn && (p.endgameTake == false || p.verifNotTakable(x, y)) {
		if time == 4 {
			return count + 1
		} else {
			return p.verifLine(x + varx, y + vary, count + 1, time + 1, varx, vary)
		}
	}
	return count
}

func (p *Gomoku)victory(x, y int) bool {
	if p.countTake[p.playerTurn - 1] == 0 || (p.verifLine(x, y, p.verifLine(x, y, 1, 1, -1, 0), 1, 1, 0) >= 5 ||
		p.verifLine(x, y, p.verifLine(x, y, 1, 1, 0, 1), 1, 0, -1) >= 5 ||
		p.verifLine(x, y, p.verifLine(x, y, 1, 1, -1, -1), 1, 1, 1) >= 5 ||
		p.verifLine(x, y, p.verifLine(x, y, 1, 1, -1, +1), 1, 1, -1) >= 5) && (p.endgameTake == false || p.verifNotTakable(x, y)) {
		return true
	}
	return false
}

func (p *Gomoku)verifDoubleThree(x, y, player int) {
	
}

func (p *Gomoku)otherPlayer() int {
	if p.playerTurn == 1 {
		return 2
	}
	return 1
}

func (p *Gomoku)verifNotTakable(x, y int) bool {
	if x <= 16 && x >= 1 && p.board[x + 1 + y * 19] == p.playerTurn &&
		(p.board[x + 2 + y * 19] == 0 || p.board[x - 1 + y * 19] == 0) && 
		(p.board[x + 2 + y * 19] == p.otherPlayer() || p.board[x - 1 + y * 19] == p.otherPlayer()) {
		return false
	}
	if x <= 17 && x >= 2 && p.board[x - 1 + y * 19] == p.playerTurn &&
		(p.board[x - 2 + y * 19] == 0 || p.board[x + 1 + y * 19] == 0) &&
		(p.board[x - 2 + y * 19] == p.otherPlayer() || p.board[x + 1 + y * 19] == p.otherPlayer()) {
		return false
	}
	if y <= 17 && y >= 2 && p.board[x + (y - 1) * 19] == p.playerTurn &&
		(p.board[x + (y - 2) * 19] == 0 || p.board[x + (y + 1) * 19] == 0) &&
		(p.board[x + (y - 2) * 19] == p.otherPlayer() || p.board[x + (y + 1) * 19] == p.otherPlayer()) {
		return false
	}
	if y <= 16 && y >= 1 && p.board[x + (y + 1) * 19] == p.playerTurn &&
		(p.board[x + (y + 2) * 19] == p.otherPlayer() || p.board[x + (y - 1) * 19] == p.otherPlayer()) &&
		(p.board[x + (y + 2) * 19] == 0 || p.board[x + (y - 1) * 19] == 0) {
		return false
	}
	if y <= 17 && y >= 2 && x >= 2 && x <= 17  && p.board[x - 1 + (y - 1) * 19] == p.playerTurn &&
		(p.board[x - 2 + (y - 2) * 19] == 0 || p.board[x + 1 + (y + 1) * 19] == 0) &&
		(p.board[x - 2 + (y - 2) * 19] == p.otherPlayer() || p.board[x + 1 + (y + 1) * 19] == p.otherPlayer()) {
		return false
	}
	if y <= 16 && y >= 2 && x >= 2 && x <= 17 && p.board[x - 1 + (y + 1) * 19] == p.playerTurn &&
		(p.board[x - 2 + (y + 2) * 19] == 0 || p.board[x + 1 + (y - 1) * 19] == 0) &&
		(p.board[x - 2 + (y + 2) * 19] == p.otherPlayer() || p.board[x + 1 + (y - 1) * 19] == p.otherPlayer()) {
		return false
	}
	if y <= 16 && y >= 2 && x <= 16 && x >= 1 && p.board[x + 1 + (y - 1) * 19] == p.playerTurn &&
		(p.board[x + 2 + (y - 2) * 19] == 0 || p.board[x - 1 + (y + 1) * 19] == 0) &&
		(p.board[x + 2 + (y - 2) * 19] == p.otherPlayer() || p.board[x - 1 + (y + 1) * 19] == p.otherPlayer()) {
		return false
	}
	if y <= 16 && y >= 1 && x <= 16 && y >= 1 && p.board[x + 1 + (y + 1) * 19] == p.playerTurn &&
		(p.board[x + 2 + (y + 2) * 19] == 0 || p.board[x - 1 + (y - 1) * 19] == 0) &&
		(p.board[x + 2 + (y + 2) * 19] == p.otherPlayer() || p.board[x - 1 + (y - 1) * 19] == p.otherPlayer()) {
		return false
	}
	return true
}

func (p *Gomoku)Play(x, y int) int {
	if p.board[x + y * 19] == 0 {
		p.board[x + y * 19] = p.playerTurn
	} else {
		fmt.Println("already occupied")
	}
	p.prise(x, y)
	if p.victory(x, y) {
		return p.playerTurn
	}
	p.changePlayerTurn()
	return 0
}

func (p *Gomoku)changePlayerTurn() {
	if p.playerTurn == 1 {
		p.playerTurn = 2
	} else {
		p.playerTurn = 1
	}
}

func (p *Gomoku)prise(x, y int) {
	if x <= 15 && p.board[x + 1 + y * 19] != p.playerTurn &&
		p.board[x + 2 + y * 19] != p.playerTurn && p.board[x + 3 + y * 19] == p.playerTurn {
		p.board[x + 2 + y * 19] = 0
		p.board[x + 1 + y * 19] = 0
		p.countTake[p.playerTurn - 1] -= 2
	}
	if x >= 3 && p.board[x - 1 + y * 19] != p.playerTurn &&
		p.board[x - 2 + y * 19] != p.playerTurn && p.board[x - 3 + y * 19] == p.playerTurn {
		p.board[x - 2 + y * 19] = 0
		p.board[x - 1 + y * 19] = 0
		p.countTake[p.playerTurn - 1] -= 2
	}
	if y >= 3 && p.board[x + (y - 1) * 19] != p.playerTurn &&
		p.board[x + (y - 2) * 19] != p.playerTurn && p.board[x + (y - 3) * 19] == p.playerTurn {
		p.board[x + (y - 1) * 19] = 0
		p.board[x + (y - 2) * 19] = 0
		p.countTake[p.playerTurn - 1] -= 2
	}
	if y <= 15 && p.board[x + (y + 1) * 19] != p.playerTurn &&
		p.board[x + (y + 2) * 19] != p.playerTurn && p.board[x + (y + 3) * 19] == p.playerTurn {
		p.board[x + (y + 1) * 19] = 0
		p.board[x + (y + 2) * 19] = 0
		p.countTake[p.playerTurn - 1] -= 2
	}
	if y >= 3 && x >= 3 && p.board[x - 1 + (y - 1) * 19] != p.playerTurn &&
		p.board[x - 2 + (y - 2) * 19] != p.playerTurn && p.board[x - 3 + (y - 3) * 19] == p.playerTurn {
		p.board[x - 1 + (y - 1) * 19] = 0
		p.board[x - 2 + (y - 2) * 19] = 0
		p.countTake[p.playerTurn - 1] -= 2
	}
	if y <= 15 && x >= 3 && p.board[x - 1 + (y + 1) * 19] != p.playerTurn &&
		p.board[x - 2 + (y + 2) * 19] != p.playerTurn && p.board[x - 3 + (y + 3) * 19] == p.playerTurn {
		p.board[x - 1 + (y + 1) * 19] = 0
		p.board[x - 2 + (y + 2) * 19] = 0
		p.countTake[p.playerTurn - 1] -= 2
	}
	if y >= 3 && x <= 15 && p.board[x + 1 + (y - 1) * 19] != p.playerTurn &&
		p.board[x + 2 + (y - 2) * 19] != p.playerTurn && p.board[x + 3 + (y - 3) * 19] == p.playerTurn {
		p.board[x + 1 + (y - 1) * 19] = 0
		p.board[x + 2 + (y - 2) * 19] = 0
		p.countTake[p.playerTurn - 1] -= 2
	}
	if y <= 15 && x <= 15 && p.board[x + 1 + (y + 1) * 19] != p.playerTurn &&
		p.board[x + 2 + (y + 2) * 19] != p.playerTurn && p.board[x + 3 + (y + 3) * 19] == p.playerTurn {
		p.board[x + 1 + (y + 1) * 19] = 0
		p.board[x + 2 + (y + 2) * 19] = 0
		p.countTake[p.playerTurn - 1] -= 2
	}
}

func (p *Gomoku)Debug_aff() {
	for i := 0;i < 19; i++ {
		for n := 0;n < 19; n++ {
			fmt.Print(p.board[i * 19 + n])
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