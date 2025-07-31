package main

import "math/rand"

type Dice struct {
	DiceCount int
}

func NewDice(diceCount int) *Dice {
	return &Dice{
		DiceCount: diceCount,
	}
}

func (d *Dice) Roll() int {
	total := 0
	for i := 0; i < d.DiceCount; i++ {
		total += rand.Intn(6) + 1
	}

	return total
}
