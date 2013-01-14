package main

import (
	"fmt"
//	"math/rand"
)

const (
	EMPTY = 0
	PION_HUMAN = 1
	PION_IA = 2
	MIN = true
	MAX = false
	MAXDEPTH = 2
	MAXINT = int(^uint(0) >> 1)
	MININT = -MAXINT - 1
	)

func max(a, b int) int {
	if a <= b { return b }
	return a
}

func min(a, b int) int {
	if a >= b { return b }
	return a
}

// Returns the index in slice of int of a given value.
// -1 if the value is nowhere to be found.
func getValueIndex(slice []int, searched int) int {
	for index, value := range slice {
		if (value == searched) {
			return index
		}
	}
    return -1
}

// Returns minimum value of a slice of int arguments
func minOfSlice(v []int) (m int) {
	if len(v) > 0 {
		m = v[0]
	}
	for i := 1; i < len(v); i++ {
		if v[i] < m {
			m = v[i]
		}
	}
	return
}


// Same as above for maximum value of a slice of int arguments
func maxOfSlice(v []int) (m int) {
	if len(v) > 0 {
		m = v[0]
	}
	for i := 1; i < len(v); i++ {
		if v[i] > m {
			m = v[i]
		}
	}
	return
}

func notAlreadyInMoves(moves *[][2]int, x, y int) bool {
	for i := 0; i < len(*moves); i++ {
		if ((*moves)[i][0] == x && (*moves)[i][1] == y) {
			return false
		}
	}
	return true
}

// Fonction qui check toutes les cases vides autour d'une pierre a un etat donne.
// On lui donne le board, la liste des moves, le x, y et i (index direct dans le board) d'une pierre
// et cette fonction va checker les 8 cases autour de cette pierre. Si ce sont des
// cases vides, alors on les ajoute a la liste des mouvements possibles (pour economiser les perfs l'IA se concentrera
// sur des coups liés à ses propres pierres en priorité).
//
// Ordre de check: à gauche, à droite, en haut, en bas, haut gauche, bas gauche, haut droit, bas droit.
// Petit Schéma de l'ordre dans lequel les cases sont checkées (avec X la pierre dont le x et le y sont donnés en paramètre):
// 5 3 7
// 1 X 2
// 6 4 8
// NB : la fonction ne check pas a droite si on est le plus a droite possible, pas en bas si on est le plus bas possible, etc. etc..
// NB 2: la fonction vérifie que chaque case n'a pas déjà été signalée comme libre afin de
// ne pas faire de doublon.
func findMovesAroundPiece(board []int, moves *[][2]int, x, y, i int) {
	if x > 0 && board[i - 1] == EMPTY && notAlreadyInMoves(moves, x - 1, y) {
		*moves = append(*moves, [2]int{x - 1, y})
	}
	if x < 18 && board[i + 1] == EMPTY && notAlreadyInMoves(moves, x + 1, y) {
		*moves = append(*moves, [2]int{x + 1, y})
	}
	if y > 0 && board[i - 19] == EMPTY && notAlreadyInMoves(moves, x, y - 1) {
		*moves = append(*moves, [2]int{x, y - 1})
	}
	if y < 18 && board[i + 19] == EMPTY && notAlreadyInMoves(moves, x, y + 1) {
		*moves = append(*moves, [2]int{x, y + 1})
	}
	if x > 0 && y > 0 && board[i - 20] == EMPTY && notAlreadyInMoves(moves, x - 1, y - 1) {
		*moves = append(*moves, [2]int{x - 1, y - 1})
	}
	if x > 0 && y < 18 && board[i + 18] == EMPTY && notAlreadyInMoves(moves, x - 1, y + 1) {
		*moves = append(*moves, [2]int{x - 1, y + 1})
	}
	if x < 18 && y > 0 && board[i - 18] == EMPTY && notAlreadyInMoves(moves, x + 1, y - 1) {
		*moves = append(*moves, [2]int{x + 1, y - 1})
	}
	if x < 18 && y < 18 && board[i + 20] == EMPTY && notAlreadyInMoves(moves, x + 1, y + 1) {
		*moves = append(*moves, [2]int{x + 1, y + 1})
	}
}

// Fonction qui va retourner les mouvements possibles pour un etat donné du plateau de jeu
// On ne cherche qu'a placer des pieces autour de pieces deja existantes.
func getPossibleMoves(game *Gomoku) [][2]int {
	var moves [][2]int
	var i, x, y int = 0, 0, 0
	
	for ; i < 361; i++ {
		if game.board[i] != EMPTY {
			findMovesAroundPiece(game.board, &moves, x, y, i)
		}
		x++
		if (x == 19) {
			x = 0
			y++
		}
	}
	return moves
}

// Copie en profondeur d'un plateau
func copyGame(dest *Gomoku, src *Gomoku) {
	dest.board = make([]int, len(src.board)) // important
	copy(dest.board, src.board)
	dest.gameType = src.gameType
 	dest.endgameTake = src.endgameTake
 	dest.doubleThree = src.doubleThree
 	dest.playerTurn = src.playerTurn
 	dest.playerTurn = src.playerTurn
 	dest.countTake[0] = src.countTake[0]
 	dest.countTake[1] = src.countTake[1]
}

// FOnction du Minimax Algorithm.
// Plan de route :
// Pour chaque noeud:
// Si on est au MAXDEPTH: calculer le score du board ainsi obtenu. (Cela ne tient pas compte du fait ou par exemple ce coup permet de gagner...)
// Sinon:
// trouver tous les coups possibles pour le prochain joueur (que ce soit le joueur ou l'IA, peu importe).
// Par coups possibles on sous-entend tous ceux qui touchent une pierre, de la couleur du joueur ou pas. Si il n'y en a pas 
// a disposition, on va 
// pour chaque mouvement : créer une copie du board, et jouer ce coup.
// Ensuite relancer min-max avec le board ainsi obtenu, stocker le retour.
// si on est dans un node min -> prendre la valeur MAX des retours
// si on est dans un node max -> prendre la valeur MIN des retours
// remonter
func minMaxAlgorithm(game *Gomoku, depth, alpha, beta int, minmax bool) int {
	if depth == MAXDEPTH {
		return gameHeuristicScore(&(game.board))
	}
	moves := getPossibleMoves(game)
	for i := 0; i < len(moves); i++ {
		var gameCopy Gomoku
		
		copyGame(&gameCopy, game)
	 	victory, _, err := gameCopy.Play(moves[i][0], moves[i][1])
		if err == nil {
			if victory == 0 {
				if (minmax == MAX) {
					alpha = max(alpha, minMaxAlgorithm(&gameCopy, depth + 1, alpha, beta, !minmax))
					if beta <= alpha {
						return alpha
					}
				} else {
					beta = min(beta, minMaxAlgorithm(&gameCopy, depth + 1, alpha, beta, !minmax))
					if beta <= alpha {
						return beta
					}
				}
			} else {
				if minmax == MIN {
					return -42
					beta = -42
				} else {
					return 42
					alpha = 42
				}
			}
		}
	}
	if minmax == MIN {
		return beta
	} else {
		return alpha
	}
	return 42 // mandatory
}

func diagonaleBottomTopCheck(board *[]int, checked *[]int, i, x, y, player int) (score int) {
	var x2, y2, j int = x - 1, y + 1, i + 18

	score = 0
	for x2 > 0 && y2 < 19 && j < 361 && (*board)[j] == player {
		score += 1
		*checked = append(*checked, j)
		j += 18
		x2--
		y2++
	}
	if (score > 4) {
		score = 0
	} else if (score == 4) {
		score = 42
	}
	if (x < 19 && y > 0 && (*board)[i - 18] == EMPTY) && (x2 > 0 && y2 < 19 && (*board)[j] == EMPTY) && score > 0 {
		score *= 2
	}
	
	return
}

func diagonaleTopBottomCheck(board *[]int, checked *[]int, i, x, y, player int) (score int) {
	var x2, y2, j int = x + 1, y + 1, i + 20

	score = 0
	for x2 < 19 && y2 < 19 && j < 361 && (*board)[j] == player {
		score += 1
		*checked = append(*checked, j)
		j += 20
		x2++
		y2++
	}
	if (score > 4) {
		score = 0
	} else if (score == 4) {
		score = 42
	}
	if (x > 0 && y > 0 && (*board)[i - 20] == EMPTY) && (x2 < 19 && y2 < 19 && (*board)[j] == EMPTY) && score > 0 {
		score *= 2
	}
	
	return
}

func verticalCheck(board *[]int, checked *[]int, i, x, y, player int) (score int) {
	var y2, j int = y + 1, i + 19

	score = 0
	for y2 < 19 && j < 361 && (*board)[j] == player {
		score += 1
		*checked = append(*checked, j)
		j += 19
		y2++
	}
	if (score > 4) {
		score = 0
	} else if (score == 4) {
		score = 42
	}
	if (y > 0 && (*board)[i - 19] == EMPTY) && (y2 < 19 && (*board)[j] == EMPTY) && score > 0 {
		score *= 2
	}
	return
}

// Fonction qui va check si la pierre a l'index i (coordonnees x y) fait partie d'une ligne horizontale
// et si oui, combien de points vaut cette ligne
// NB : on ajoute +1 point à toute combinaison étant "open" (exemple: XOOOOX avec X des cases libres)
// NB 2 : une ligne détectée comme gagnante (5 pions alignés )vaut +42 points. En revanche, une ligne de + de 5 pions est totalement inutile : dans ce cas on ramène le score à zéro.
func horizontalCheck(board *[]int, HChecked *[]int, i, x, y, player int) (score int) {
	var x2, j int = x + 1, i + 1

	score = 0
	for x2 < 19 && j <= 360 && (*board)[j] == player {
		score += 1
		*HChecked = append(*HChecked, j)
		j++
		x2++
	}
	if (score > 4) {
		score = 0
	} else if (score == 4) {
		score = 42
	}
	if (x > 0 && (*board)[i - 1] == EMPTY) && (x2 < 19 && (*board)[j] == EMPTY) && score > 0 {
		score *= 2
	}
	
	return
}

func calculatePionValue(board *[]int, HChecked, VChecked, DTBChecked, DBTChecked *[]int, i, x, y int) (pionScore int) {
	var player = (*board)[i]

	pionScore = 0
	if getValueIndex(*HChecked, i) == -1 {
		pionScore += horizontalCheck(board, HChecked, i, x, y, player)
		*HChecked = append(*HChecked, i)
	}
	if getValueIndex(*VChecked, i) == -1 {
		pionScore += verticalCheck(board, VChecked, i, x, y, player)
		*VChecked = append(*VChecked, i)
	}
	if getValueIndex(*DTBChecked, i) == -1 {
		pionScore += diagonaleTopBottomCheck(board, DTBChecked, i, x, y, player)
		*DTBChecked = append(*DTBChecked, i)
	}
	if getValueIndex(*DBTChecked, i) == -1 {
		pionScore += diagonaleBottomTopCheck(board, DBTChecked, i, x, y, player)
		*DBTChecked = append(*DBTChecked, i)
	}
	return
}

func gameHeuristicScore(board *[]int) int {
 	var IAScore, HumanScore, x, y, i int = 0, 0, 0, 0, 0
	var HuHChecked, HuVChecked, HuDTBChecked, HuDBTChecked []int
	var IAHChecked, IAVChecked, IADTBChecked, IADBTChecked []int

	for i = 0; i < 361; i++ {
		if (*board)[i] == PION_IA {
			IAScore += calculatePionValue(board, &IAHChecked, &IAVChecked, &IADTBChecked, &IADBTChecked, i, x, y)
		} else if (*board)[i] == PION_HUMAN {
			HumanScore += calculatePionValue(board, &HuHChecked, &HuVChecked, &HuDTBChecked, &HuDBTChecked, i, x, y)
		}
		x++;
		if x == 19 {
			x = 0
			y++
		}
	}
	return (IAScore - HumanScore)
}

func findBestMoveAccordingScores(scores [][3]int) (int, int) {
	var bestMoveX, bestMoveY, bestScore int = 0, 0, 0

	if len(scores) > 0 {
		bestScore, bestMoveX, bestMoveY  = scores[0][0], scores[0][1], scores[0][2]
	} else {
		fmt.Println("This should not happen")
		// should not happen, throw a random number
	}
	for i := 0; i < len(scores); i++ {
		if scores[i][0] > bestScore {
			bestScore, bestMoveX, bestMoveY  = scores[i][0], scores[i][1], scores[i][2]
		}
	}
	return bestMoveX, bestMoveY
}

func checkImmediateThreats(board *[]int) (int, int, bool) {
	var HuHChecked, HuVChecked, HuDTBChecked, HuDBTChecked []int
	var x, y int = 0, 0 

	for i := 0; i < len(*board); i++ {
		if (*board)[i] == PION_HUMAN {
			fmt.Println("Valeur du pion en ", x, y, ":", calculatePionValue(board, &HuHChecked, &HuVChecked, &HuDTBChecked, &HuDBTChecked, i, x, y))
		}
		x++
		if x == 19 {
			x = 0
			y++
		}
	}
	return (*board)[0], (*board)[1], false 
}

func findBestMove(game *Gomoku) (bestX, bestY int) {
	var bestScore, alpha int = MININT, 0
	
	moves := getPossibleMoves(game)
	for i := 0; i < len(moves); i++ {
		var gameCopy Gomoku
		
		copyGame(&gameCopy, game)
		vic, _, err := gameCopy.Play(moves[i][0], moves[i][1])
		if err == nil {
			if vic != 0 {
				return moves[i][0], moves[i][1]
			}
			alpha = minMaxAlgorithm(&gameCopy, 0, MININT, MAXINT, MIN)
			if alpha > bestScore {
				bestX, bestY = moves[i][0], moves[i][1]
				bestScore = alpha
			}
		}
	}
	return
}

func IATurn(game *Gomoku) (x int, y int) {
	x, y = findBestMove(game)
	return
}