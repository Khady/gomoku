package main

import (
	"strconv"
	"strings"
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
)

// g_MaxDepth est une globale definie dans ia.go
func setMaxDepth(timeout int) {
	if timeout <= 4 {
		g_MaxDepth = 1
	} else if timeout <= 250 {
		g_MaxDepth = 2
	} else if timeout <= 4000 {
		g_MaxDepth = 3
	} else {
		g_MaxDepth = 4
	}
	println("depth", g_MaxDepth)
}

func playIATurn(conn net.Conn, words []string, game *Gomoku) bool {
	x, y := IATurn(game)
	fmt.Fprintf(conn, "PLAY %d %d\n", x, y)
	game.Play(x, y)
	return true
}

func playOpponentTurn(conn net.Conn, words []string, game *Gomoku) bool {
	x, _ := strconv.Atoi(words[1])	
	y, _ := strconv.Atoi(words[2])
	game.Play(x, y)
	return true
}

func endGame(conn net.Conn, words []string, game *Gomoku) bool {
	reasons := make(map[string]string)
	
	reasons["CAPTURE\n"]   = "capture de 10 pierres ou plus"
	reasons["FIVEALIGN\n"] = "alignement de 5 pierres"
	reasons["RULEERR\n"]   = "non-respect des regles"
	reasons["TIMEOUT\n"]   = "timeout"
	println(words[0], "! Reason:",reasons[words[1]])
	return false
} 

// Etant donne qu'on gere deja la prise des pierres a l'interieur de l'arbitre, cette commande ne nous sert a rien. Il faudrait a la rigueur check que les pierres annoncées comme prises le sont bien sur notre goban...
func checkPrise(conn net.Conn, words []string, game *Gomoku) bool {

	return true
}

func learnRules(conn net.Conn, words []string, game *Gomoku) bool {
	if words[1] == "1" {
		game.doubleThree = true
	} else { game.doubleThree = false }
	if words[2] == "1" {
		game.endgameTake = true
	} else { game.endgameTake = false }
	timeout, _ := strconv.Atoi(strings.Split(words[3], "\n")[0])
	setMaxDepth(timeout)
	return true
}

func getProtocolDictionary() map[string]func(net.Conn, []string, *Gomoku) bool {
	protoDico := make(map[string]func(net.Conn, []string, *Gomoku) bool)

	protoDico["RULES"] = learnRules
	protoDico["REM"] = checkPrise
	protoDico["ADD"] = playOpponentTurn
	protoDico["WIN"] = endGame
	protoDico["LOSE"] = endGame
	protoDico["YOURTURN\n"] = playIATurn
	return protoDico
}

func gameLoop(conn net.Conn) {
	var game Gomoku = Gomoku{make([]int, 361), true, false, false, 1, [2]int{10, 10}}
	var canContinue bool
	
	functionsDictionary := getProtocolDictionary()
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err == nil {
			words := strings.Split(message, " ")
			canContinue = functionsDictionary[words[0]](conn, words, &game)
		} else {
			fmt.Println("Error:", err)
			canContinue = false
		}
		if !canContinue {
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
		conn, err := net.Dial("tcp", buffer.String())
		if err == nil {
			fmt.Fprintf(conn, "CONNECT CLIENT\n")
			gameLoop(conn)
		} else {
			fmt.Println(err)
		}
	}
}