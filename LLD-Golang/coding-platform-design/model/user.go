package model

type User struct {
	UserID               int
	Name                 string
	Email                string
	ParticipatedContests []string
	Score                int
}
