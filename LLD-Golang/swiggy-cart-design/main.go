package main

import (
	"errors"
	"fmt"
)

// Product represents an item available for purchase.
type Product struct {
	ID                string
	Name              string
	Price             float64
	InventoryQuantity int
}

// CartItem represents an item in the cart with quantity.
type CartItem struct {
	Product  *Product
	Quantity int
}

// InventoryService interface for managing inventory (Strategy/Dependency Injection).
type InventoryService interface {
	CheckAvailability(productID string, quantity int) bool
	UpdateInventory(productID string, quantity int) error
}

// MockInventoryService is a simple in-memory implementation for demo.
type MockInventoryService struct {
	inventory map[string]int // productID -> quantity
}

func NewMockInventoryService() *MockInventoryService {
	return &MockInventoryService{
		inventory: make(map[string]int),
	}
}

func (m *MockInventoryService) CheckAvailability(productID string, quantity int) bool {
	avail, ok := m.inventory[productID]
	if !ok {
		return false
	}
	return avail >= quantity
}

func (m *MockInventoryService) UpdateInventory(productID string, quantity int) error {
	avail, ok := m.inventory[productID]
	if !ok {
		return errors.New("product not found in inventory")
	}
	if avail+quantity < 0 { // quantity can be negative for add (reserve) or positive for remove (release)
		return errors.New("insufficient inventory")
	}
	m.inventory[productID] += quantity
	return nil
}

// DiscountStrategy interface for applying discounts (Strategy Pattern).
type DiscountStrategy interface {
	ApplyDiscount(total float64) float64
}

// PercentageDiscount is an example strategy (e.g., 10% off).
type PercentageDiscount struct {
	Percentage float64
}

func (p *PercentageDiscount) ApplyDiscount(total float64) float64 {
	return total * (1 - p.Percentage/100)
}

// PaymentProcessor interface for handling payments (Strategy Pattern).
type PaymentProcessor interface {
	ProcessPayment(amount float64) error
}

// MockPaymentProcessor is a demo implementation.
type MockPaymentProcessor struct{}

func (m *MockPaymentProcessor) ProcessPayment(amount float64) error {
	fmt.Printf("Processing payment of %.2f\n", amount)
	return nil // Simulate success
}

// NotificationService interface for sending notifications (Observer-like).
type NotificationService interface {
	SendNotification(message string) error
}

// MockNotificationService is a demo implementation.
type MockNotificationService struct{}

func (m *MockNotificationService) SendNotification(message string) error {
	fmt.Println("Notification:", message)
	return nil
}

// Cart is the main struct managing the shopping cart.
type Cart struct {
	Items               []CartItem
	Total               float64
	inventoryService    InventoryService
	discountStrategy    DiscountStrategy
	paymentProcessor    PaymentProcessor
	notificationService NotificationService
}

// NewCart constructor injects dependencies (Dependency Injection).
func NewCart(inventory InventoryService, discount DiscountStrategy, payment PaymentProcessor, notification NotificationService) *Cart {
	return &Cart{
		Items:               make([]CartItem, 0),
		inventoryService:    inventory,
		discountStrategy:    discount,
		paymentProcessor:    payment,
		notificationService: notification,
	}
}

// AddItem adds a product to the cart, validates inventory, updates total, notifies.
func (c *Cart) AddItem(product *Product, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if !c.inventoryService.CheckAvailability(product.ID, quantity) {
		return errors.New("insufficient inventory")
	}
	// Reserve inventory (negative update)
	if err := c.inventoryService.UpdateInventory(product.ID, -quantity); err != nil {
		return err
	}

	// Check if item already exists
	for i, item := range c.Items {
		if item.Product.ID == product.ID {
			c.Items[i].Quantity += quantity
			c.notify("Added more of " + product.Name + " to cart")
			c.CalculateTotal()
			return nil
		}
	}
	// New item
	c.Items = append(c.Items, CartItem{Product: product, Quantity: quantity})
	c.notify("Added " + product.Name + " to cart")
	c.CalculateTotal()
	return nil
}

// RemoveItem reduces quantity, releases inventory, updates total, notifies.
func (c *Cart) RemoveItem(productID string, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	for i, item := range c.Items {
		if item.Product.ID == productID {
			if item.Quantity < quantity {
				return errors.New("not enough quantity in cart")
			}
			item.Quantity -= quantity
			// Release inventory (positive update)
			if err := c.inventoryService.UpdateInventory(productID, quantity); err != nil {
				return err
			}
			if item.Quantity == 0 {
				c.Items = append(c.Items[:i], c.Items[i+1:]...)
			}
			c.notify("Removed some of " + item.Product.Name + " from cart")
			c.CalculateTotal()
			return nil
		}
	}
	return errors.New("product not found in cart")
}

// DeleteItem removes entire item, releases all inventory, updates total, notifies.
func (c *Cart) DeleteItem(productID string) error {
	for i, item := range c.Items {
		if item.Product.ID == productID {
			// Release all inventory
			if err := c.inventoryService.UpdateInventory(productID, item.Quantity); err != nil {
				return err
			}
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			c.notify("Deleted " + item.Product.Name + " from cart")
			c.CalculateTotal()
			return nil
		}
	}
	return errors.New("product not found in cart")
}

// CalculateTotal computes the total with discount applied.
func (c *Cart) CalculateTotal() float64 {
	subtotal := 0.0
	for _, item := range c.Items {
		subtotal += item.Product.Price * float64(item.Quantity)
	}
	c.Total = c.discountStrategy.ApplyDiscount(subtotal)
	c.notify(fmt.Sprintf("Cart total updated to %.2f", c.Total))
	return c.Total
}

// ProcessPayment handles payment for the current total.
func (c *Cart) ProcessPayment() error {
	if c.Total <= 0 {
		return errors.New("cart is empty")
	}
	err := c.paymentProcessor.ProcessPayment(c.Total)
	if err == nil {
		c.notify("Payment successful!")
	}
	return err
}

// notify is a helper to send notifications (Observer pattern integration).
func (c *Cart) notify(message string) {
	_ = c.notificationService.SendNotification(message) // Ignore error for simplicity
}

// Demo main function to test the cart.
func main() {
	// Setup mocks
	inventory := NewMockInventoryService()
	inventory.inventory["p1"] = 10 // Add product inventory
	inventory.inventory["p2"] = 5

	discount := &PercentageDiscount{Percentage: 10}
	payment := &MockPaymentProcessor{}
	notification := &MockNotificationService{}

	cart := NewCart(inventory, discount, payment, notification)

	product1 := &Product{ID: "p1", Name: "Pizza", Price: 10.0, InventoryQuantity: 10}
	product2 := &Product{ID: "p2", Name: "Burger", Price: 5.0, InventoryQuantity: 5}

	// Add items
	cart.AddItem(product1, 2) // Expect notifications and total update
	cart.AddItem(product2, 1)

	// Remove
	cart.RemoveItem("p2", 1)

	// Delete
	cart.DeleteItem("p1")

	// Try add again and pay
	cart.AddItem(product1, 3)
	cart.ProcessPayment()
}
