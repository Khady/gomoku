package main

import "fmt"

func main() {
	// game.Play(1, 0)
	// game.Play(2, 0)
	// game.Play(1, 1)
	// game.Play(6, 2)
	// game.Play(1, 2)
	// game.Play(5, 0)
	// game.Play(1, 3)
	// game.Play(7, 0)
	// game.Play(2, 2)
	// game.Play(8, 0)
	// vic := game.Play(1, 4)
	// if vic != 0 {
	// 	fmt.Println("player ",  vic, "win")
	// }
	// game.Play(9, 0)
	// if game.Play(12, 5) != 0 {
	// 	fmt.Println("win")
	// }
	// game.Debug_aff()
	for {
		mode, doubleThree, endGame := game_mode(WINNER)
		game := Gomoku{make([]int, 361), true, endGame, doubleThree, 1, [2]int{10, 10}}
		fmt.Println("Game's mode is", mode)
		if mode == 0 {
			return
		}
		board_display(game)
	}
}
