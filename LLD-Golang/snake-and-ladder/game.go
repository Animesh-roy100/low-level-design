package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Game struct {
	Board      *Board
	Dice       *Dice
	Players    []*Player
	PlayerTurn int
	Winner     *Player
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	game := &Game{
		Board:   NewBoard(10, 5, 4),
		Dice:    NewDice(1),
		Players: []*Player{NewPlayer("p1"), NewPlayer("p2")},
	}
	return game
}

func (g *Game) Start() {
	for g.Winner == nil {
		player := g.Players[g.PlayerTurn]
		fmt.Printf("%s's turn. Current position: %d\n", player.ID, player.CurrentPosition)

		// Roll Dice
		roll := g.Dice.Roll()
		newPosition := player.CurrentPosition + roll
		fmt.Printf("%s rolled %d. Moving to %d\n", player.ID, roll, newPosition)

		// Check board bounds
		if newPosition >= g.Board.Size*g.Board.Size {
			newPosition = player.CurrentPosition
		}

		// Check for jumps
		if cell := g.Board.getCell(newPosition); cell.Jump != nil {
			jumpType := "ladder"
			if cell.Jump.Start > cell.Jump.End {
				jumpType = "snake"
			}

			fmt.Printf("%s encountered a %s! ", player.ID, jumpType)
			newPosition = cell.Jump.End
			fmt.Printf("Moved to %d\n", newPosition)
		}

		// Update position
		player.CurrentPosition = newPosition

		// Check win condition
		if newPosition == g.Board.Size*g.Board.Size-1 {
			g.Winner = player
			fmt.Printf("%s wins!\n", player.ID)
			return
		}

		// Next player
		g.nextTurn()
	}
}

func (g *Game) nextTurn() {
	g.PlayerTurn = (g.PlayerTurn + 1) % len(g.Players)
}
