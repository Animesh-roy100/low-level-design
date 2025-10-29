package main

import (
	"errors"
	"fmt"
	"sync"
)

// ──────────────────────────────────────────────────────────────
// 1. PRODUCT & CART ITEM
// ──────────────────────────────────────────────────────────────
type Product struct {
	ID                string
	Name              string
	Price             float64
	InventoryQuantity int
}

type CartItem struct {
	Product  *Product
	Quantity int
}

// ──────────────────────────────────────────────────────────────
// 2. SINGLETON SERVICES
// ──────────────────────────────────────────────────────────────
type InventoryService struct {
	products map[string]Product
	mu       sync.Mutex
}

var (
	inventoryOnce sync.Once
	inventoryInst *InventoryService
)

func Inventory() *InventoryService {
	inventoryOnce.Do(func() {
		inventoryInst = &InventoryService{
			products: make(map[string]Product),
		}
	})
	return inventoryInst
}

func (is *InventoryService) AddProduct(p Product) { is.products[p.ID] = p }
func (is *InventoryService) Get(id string) (*Product, error) {
	p, ok := is.products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return &p, nil
}
func (is *InventoryService) Available(id string, qty int) bool {
	p, _ := is.Get(id)
	return p != nil && p.InventoryQuantity >= qty
}
func (is *InventoryService) Reserve(id string, qty int) error {
	is.mu.Lock()
	defer is.mu.Unlock()
	p, err := is.Get(id)
	if err != nil {
		return err
	}
	if p.InventoryQuantity < qty {
		return errors.New("insufficient stock")
	}
	p.InventoryQuantity -= qty
	is.products[id] = *p
	return nil
}
func (is *InventoryService) Release(id string, qty int) {
	is.mu.Lock()
	defer is.mu.Unlock()
	if p, ok := is.products[id]; ok {
		p.InventoryQuantity += qty
		is.products[id] = p
	}
}

// ──────────────────────────────────────────────────────────────
// 3. DISCOUNT SERVICE (Singleton + Strategy for rules)
// ──────────────────────────────────────────────────────────────
type DiscountRule struct {
	Code     string
	Percent  float64 // 0.1 = 10%
	MinTotal float64 // optional condition
}

type DiscountService struct {
	rules []DiscountRule
	mu    sync.Mutex
}

var (
	discountOnce sync.Once
	discountInst *DiscountService
)

func Discount() *DiscountService {
	discountOnce.Do(func() {
		discountInst = &DiscountService{
			rules: []DiscountRule{
				{Code: "SWIGGY10", Percent: 0.10},
				{Code: "FIRST50", Percent: 0.50, MinTotal: 100},
			},
		}
	})
	return discountInst
}

func (ds *DiscountService) Apply(subtotal float64, code string) float64 {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for _, r := range ds.rules {
		if r.Code == code && subtotal >= r.MinTotal {
			return subtotal * (1 - r.Percent)
		}
	}
	return subtotal // no discount
}

// ──────────────────────────────────────────────────────────────
// 4. PAYMENT STRATEGY
// ──────────────────────────────────────────────────────────────
type PaymentService interface {
	Process(amount float64) (bool, error)
}

// Mock (for demo)
type MockPayment struct{}

func (m MockPayment) Process(amount float64) (bool, error) {
	fmt.Printf("[Payment] Charged %.2f\n", amount)
	return true, nil
}

// Real gateway (example)
type StripePayment struct{}

func (s StripePayment) Process(amount float64) (bool, error) {
	// call Stripe API …
	return true, nil
}

// ──────────────────────────────────────────────────────────────
// 5. NOTIFICATION OBSERVER
// ──────────────────────────────────────────────────────────────
type Observer interface {
	Update(message string)
}

type NotificationService struct {
	subscribers []Observer
	mu          sync.Mutex
}

var (
	notifyOnce sync.Once
	notifyInst *NotificationService
)

func Notifier() *NotificationService {
	notifyOnce.Do(func() {
		notifyInst = &NotificationService{}
	})
	return notifyInst
}

func (ns *NotificationService) Subscribe(o Observer) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.subscribers = append(ns.subscribers, o)
}

func (ns *NotificationService) Send(msg string) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	for _, s := range ns.subscribers {
		s.Update(msg)
	}
}

// Simple console observer
type ConsoleObserver struct{}

func (c ConsoleObserver) Update(msg string) {
	fmt.Println("Notification:", msg)
}

// ──────────────────────────────────────────────────────────────
// 6. CART (core domain) – uses all services
// ──────────────────────────────────────────────────────────────
type Cart struct {
	UserID string
	Items  []CartItem
	Total  float64

	// injected services (DI)
	inventory *InventoryService
	discount  *DiscountService
	payment   PaymentService
	notifier  *NotificationService
}

func NewCart(userID string, payment PaymentService) *Cart {
	c := &Cart{
		UserID:    userID,
		inventory: Inventory(),
		discount:  Discount(),
		payment:   payment,
		notifier:  Notifier(),
	}
	// auto-subscribe cart as observer for its own events
	c.notifier.Subscribe(c) // Cart implements Observer for internal logging
	return c
}

// Implement Observer to log its own events
func (c *Cart) Update(msg string) {
	fmt.Printf("[Cart %s] %s\n", c.UserID, msg)
}

// ----- Cart operations -------------------------------------------------
func (c *Cart) Add(productID string, qty int) error {
	if !c.inventory.Available(productID, qty) {
		c.notifier.Send("Insufficient stock for " + productID)
		return errors.New("stock unavailable")
	}
	if err := c.inventory.Reserve(productID, qty); err != nil {
		return err
	}
	p, _ := c.inventory.Get(productID)

	for i := range c.Items {
		if c.Items[i].Product.ID == productID {
			c.Items[i].Quantity += qty
			c.recalculate("")
			c.notifier.Send(fmt.Sprintf("Added %d more %s", qty, p.Name))
			return nil
		}
	}
	c.Items = append(c.Items, CartItem{Product: p, Quantity: qty})
	c.recalculate("")
	c.notifier.Send(fmt.Sprintf("Added %s (x%d)", p.Name, qty))
	return nil
}

func (c *Cart) Remove(productID string, qty int) error {
	for i := range c.Items {
		if c.Items[i].Product.ID == productID {
			if c.Items[i].Quantity < qty {
				return errors.New("cannot remove more than present")
			}
			c.Items[i].Quantity -= qty
			c.inventory.Release(productID, qty)
			if c.Items[i].Quantity == 0 {
				c.Items = append(c.Items[:i], c.Items[i+1:]...)
			}
			c.recalculate("")
			c.notifier.Send(fmt.Sprintf("Removed %d of %s", qty, c.Items[i].Product.Name))
			return nil
		}
	}
	return errors.New("item not in cart")
}

func (c *Cart) Delete(productID string) error {
	for i := range c.Items {
		if c.Items[i].Product.ID == productID {
			c.inventory.Release(productID, c.Items[i].Quantity)
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			c.recalculate("")
			c.notifier.Send("Deleted " + c.Items[i].Product.Name + " from cart")
			return nil
		}
	}
	return errors.New("item not found")
}

func (c *Cart) recalculate(promo string) {
	sub := 0.0
	for _, it := range c.Items {
		sub += it.Product.Price * float64(it.Quantity)
	}
	c.Total = c.discount.Apply(sub, promo)
}

// Checkout uses Builder-like fluent API (optional)
func (c *Cart) Checkout(promo string) error {
	c.recalculate(promo)
	ok, err := c.payment.Process(c.Total)
	if err != nil || !ok {
		c.notifier.Send("Payment failed")
		return errors.New("payment failed")
	}
	c.notifier.Send(fmt.Sprintf("Order placed! Total: %.2f", c.Total))
	c.Items = nil
	c.Total = 0
	return nil
}

// ──────────────────────────────────────────────────────────────
// 7. DEMO MAIN
// ──────────────────────────────────────────────────────────────
func main() {
	// seed inventory
	Inventory().AddProduct(Product{ID: "p1", Name: "Margherita Pizza", Price: 12.99, InventoryQuantity: 10})
	Inventory().AddProduct(Product{ID: "p2", Name: "Veggie Burger", Price: 5.99, InventoryQuantity: 20})

	// register console observer (could be email, push, etc.)
	Notifier().Subscribe(ConsoleObserver{})

	// create cart with mock payment
	cart := NewCart("u123", MockPayment{})

	_ = cart.Add("p1", 2) // success
	_ = cart.Add("p2", 3) // success
	_ = cart.Add("p1", 9) // fail → insufficient stock

	cart.recalculate("SWIGGY10")
	fmt.Printf("Subtotal after discount: %.2f\n", cart.Total)

	_ = cart.Remove("p2", 1)
	_ = cart.Delete("p1")

	_ = cart.Checkout("SWIGGY10")
}