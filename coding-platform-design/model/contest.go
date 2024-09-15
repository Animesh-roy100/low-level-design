package model

type Contest struct {
	ContestID    int
	Name         string
	Level        DifficultyLevel
	Questions    []Question
	Creator      User
	Participants []User
}
