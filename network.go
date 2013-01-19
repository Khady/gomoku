package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
)

func playOpponentTurn(words []string, game *Gomoku) int {
	if len(words) >= 3 {
		
		
	}
}

func checkPrise(words []string, game *Gomoku) int {
	if len(words) >= 4 {
		if words[1] == "1" {
			game.doubleThree = 1
		} else { game.doubleThree = 0 }
		if words[2] == "1" {
			game.endgameTake = 1
		} else { game.endgameTake = 0 }
	} 
	return 42
}

func learnRules(words []string, game *Gomoku) int {
	if len(words) >= 4 {
		if words[1] == "1" {
			game.doubleThree = 1
		} else { game.doubleThree = 0 }
		if words[2] == "1" {
			game.endgameTake = 1
		} else { game.endgameTake = 0 }
	} 
	return 42
}

func fillDictionary(dico map[string]func([]string, *Gomoku) int) (map[string]func([]string, *Gomoku) int) {
	dico["RULES"] = learnRules
	dico["REM"] = checkPrise
	dico["ADD"] = playOpponentTurn
	dico["WIN"] = victory
	dico["LOSE"] = defeat
	dico["YOURTURN\n"] = playIATurn
	return dico
}

func gameLoop(conn net.Conn) {
	var game Gomoku = Gomoku{make([]int, 361), true, false, false, 1, [2]int{10, 10}}
	var functionsDictionary map[string]func([]string, *Gomoku) int

	functionsDictionary = make(map[string]func([]string, *Gomoku) int)
	functionsDictionary = fillDictionary(functionsDictionary)
	game.playerTurn = 1
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err == nil {
			println(message)
			fmt.Fprintf(conn, "Y A BON BANANIA\n")
		} else {
			fmt.Println(err)
			break
		}
	}
}

func main() {
	var buffer bytes.Buffer
	
	if len(os.Args) < 3 {
		fmt.Println("Usage:", os.Args[0], "server port")
	} else {
		buffer.WriteString(os.Args[1])
		buffer.WriteString(":")
		buffer.WriteString(os.Args[2])
		println(buffer.String())
		conn, err := net.Dial("tcp", buffer.String())
		if err == nil {
			fmt.Fprintf(conn, "CONNECT CLIENT\n")
			gameLoop(conn)
		} else {
			fmt.Println(err)
		}
	}
}