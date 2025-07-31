package main

type Player struct {
	ID              string
	CurrentPosition int
}

func NewPlayer(id string) *Player {
	return &Player{
		ID:              id,
		CurrentPosition: 0,
	}
}
