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
		game := Gomoku{make([]int, 361), true, true, true, 1, [2]int{10, 10}}
		mode := game_mode(WINNER)
		if mode == 2 {
			return
		}
		fmt.Println("Game's mode is", mode)
		board_display(game)
	}
}
