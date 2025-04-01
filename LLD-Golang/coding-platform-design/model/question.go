package model

type DifficultyLevel string

const (
	Easy   DifficultyLevel = "Easy"
	Medium DifficultyLevel = "Medium"
	Hard   DifficultyLevel = "Hard"
)

type Question struct {
	QuestionID int
	Level      DifficultyLevel
	Score      int
}
