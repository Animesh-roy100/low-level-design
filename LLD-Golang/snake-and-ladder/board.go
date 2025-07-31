package main

import "math/rand"

type Board struct {
	Size  int
	Cells [][]*Cell
}

func NewBoard(size, snakes, ladders int) *Board {
	board := &Board{Size: size}
	board.initializeCells()
	board.addSnakesAndLadders(snakes, ladders)
	return board
}

func (b *Board) initializeCells() {
	b.Cells = make([][]*Cell, b.Size)
	for i := range b.Cells {
		b.Cells[i] = make([]*Cell, b.Size)
		for j := range b.Cells[i] {
			b.Cells[i][j] = &Cell{}
		}
	}
}

func (b *Board) addSnakesAndLadders(snakes, ladders int) {
	totalCells := b.Size * b.Size

	for snakes > 0 {
		start := rand.Intn(totalCells-2) + 1
		end := rand.Intn(totalCells-2) + 1

		if start <= end || b.getCell(start).Jump != nil {
			continue
		}

		b.getCell(start).Jump = &Jump{Start: start, End: end}
		snakes--
	}

	for ladders > 0 {
		start := rand.Intn(totalCells-2) + 1
		end := rand.Intn(totalCells-2) + 1

		if start >= end || b.getCell(start).Jump != nil {
			continue
		}

		b.getCell(start).Jump = &Jump{Start: start, End: end}
		ladders--
	}
}

func (b *Board) getCell(position int) *Cell {
	row := position / b.Size
	col := position % b.Size
	return b.Cells[row][col]
}
