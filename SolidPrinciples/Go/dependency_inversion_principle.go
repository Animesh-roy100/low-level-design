package main

import "fmt"

/*

High-level modules should not depend on low-level modules; both should depend on abstractions.
Abstractions should not depend on details; details should depend on abstractions.

Dependency Inversion Principle ensures that high-level business logic depends on interfaces instead of concrete implementations,
making the system loosely coupled, testable, and easy to extend.

Before - Violates DIP


type MySQLDB struct{}

func (db *MySQLDB) Save(data string) {
	fmt.Println("Saving data to MySQL:", data)
}

type OrderService struct {
	db *MySQLDB    ------> This is a concrete dependency, should depend on abstraction instead
}

func (os *OrderService) CreateOrder(order string) {
	os.db.Save(order)
}

func main() {
	db := &MySQLDB{}
	service := &OrderService{db: db}
	service.CreateOrder("Order#123")
}

*/

type Repository interface {
	Save(data string)
}

type MySQLDB struct{}

func (db *MySQLDB) Save(data string) {
	fmt.Println("Saving data to MySQL:", data)
}

type PostgresDB struct{}

func (db *PostgresDB) Save(data string) {
	fmt.Println("Saving data to Postgres:", data)
}

type OrderService struct {
	repo Repository // Depends on abstraction
}

func NewOrderService(r Repository) *OrderService {
	return &OrderService{repo: r}
}

func (os *OrderService) CreateOrder(order string) {
	os.repo.Save(order)
}

// func main() {
// 	mysqlDB := &MySQLDB{}
// 	postgresDB := &PostgresDB{}

// 	mysqlService := NewOrderService(mysqlDB)
// 	postgresService := NewOrderService(postgresDB)

// 	mysqlService.CreateOrder("Order#123")
// 	postgresService.CreateOrder("Order#456")
// }
