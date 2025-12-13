package main

import "fmt"

/*
A class should have only one reason to change.
This means a class should be responsible for a single, specific piece of functionality.
*/

// ---------- BEFORE (Violates SRP) ----------

type User struct {
	Name  string
	Email string
}

func (u *User) Save() {
	fmt.Println("User saved to database")
}

func (u *User) SendEmail(message string) {
	fmt.Printf("Email sent to %s: %s\n", u.Email, message)
}

// ---------- AFTER (Follows SRP) ----------

// Entity: holds only user data
type UserSRP struct {
	Name  string
	Email string
}

// Repository: responsible only for persistence
type UserRepository struct{}

func (ur *UserRepository) Save(u *UserSRP) {
	fmt.Println("User saved to database")
}

// Service: responsible only for email communication
type EmailService struct{}

func (es *EmailService) SendEmail(email, message string) {
	fmt.Printf("Email sent to %s: %s\n", email, message)
}

func main() {
	// Violates SRP
	user := &User{Name: "John Doe", Email: "john@example.com"}
	user.Save()
	user.SendEmail("Hello!")

	// SRP-compliant
	userSRP := &UserSRP{Name: "John Doe", Email: "john@example.com"}

	userRepo := &UserRepository{}
	emailService := &EmailService{}

	userRepo.Save(userSRP)
	emailService.SendEmail(userSRP.Email, "Hello!")
}
