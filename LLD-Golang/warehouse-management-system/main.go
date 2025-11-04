package main

import (
	"errors"
	"fmt"
)

// =====================================================
// WAREHOUSE MANAGEMENT SYSTEM (Go, single file)
// Core entities + patterns requested:
// - Factory Pattern for Shipment creation
// - Observer Pattern for stock alerts (low inventory)
// =====================================================

// -----------------------------
// StorageLocation (1:1 with Item)
// -----------------------------
type StorageLocation struct {
	ID               string
	Capacity         float64
	CurrentOccupancy float64
	Type             string // shelf, bin, pallet
}

func NewStorageLocation(id string, capacity float64, typ string) *StorageLocation {
	return &StorageLocation{ID: id, Capacity: capacity, Type: typ}
}

func (s *StorageLocation) CanAccommodate(size float64) bool {
	return s.CurrentOccupancy+size <= s.Capacity
}

func (s *StorageLocation) AddOccupancy(size float64) error {
	if !s.CanAccommodate(size) {
		return fmt.Errorf("not enough space in location %s", s.ID)
	}
	s.CurrentOccupancy += size
	return nil
}

// -----
// Item
// -----
type Item struct {
	SKU          string
	Name         string
	Quantity     int
	Size         float64
	ReorderLevel int
	Location     *StorageLocation // one-to-one
}

func NewItem(sku, name string, qty int, size float64, reorder int) *Item {
	return &Item{SKU: sku, Name: name, Quantity: qty, Size: size, ReorderLevel: reorder}
}

func (i *Item) SetLocation(loc *StorageLocation) error {
	if err := loc.AddOccupancy(i.Size); err != nil {
		return err
	}
	i.Location = loc
	return nil
}

func (i *Item) UpdateStock(delta int) { i.Quantity += delta }

func (i *Item) IsReorderNeeded() bool { return i.Quantity < i.ReorderLevel }

// -------------
// User & Roles
// -------------
const (
	RoleAdmin   = "ADMIN"
	RoleManager = "MANAGER"
	RoleWorker  = "WORKER"
)

type User struct {
	ID    string
	Name  string
	Email string
	Role  string // ADMIN / MANAGER / WORKER
}

// ---------------
// Order & Items
// ---------------
type OrderItem struct {
	Item     *Item
	Quantity int
	Price    float64 // unit price
}

type Order struct {
	OrderNumber string
	Customer    string
	Status      string // pending, fulfilled, shipped
	Items       []OrderItem
	ManagedBy   *User // many orders can be managed by one user
}

func NewOrder(orderNumber, customer string, managedBy *User) *Order {
	return &Order{OrderNumber: orderNumber, Customer: customer, Status: "pending", ManagedBy: managedBy}
}

func (o *Order) AddOrderItem(oi OrderItem) { o.Items = append(o.Items, oi) }

func (o *Order) TotalCost() float64 {
	sum := 0.0
	for _, it := range o.Items {
		sum += float64(it.Quantity) * it.Price
	}
	return sum
}

func (o *Order) SetStatus(status string) { o.Status = status }

// =====================================================
// Observer Pattern (Stock Alerts)
// =====================================================
// Observers subscribe to low-stock notifications for Items.

type StockObserver interface{ Update(lowItem *Item) }

type StockNotifier struct{ observers []StockObserver }

func (n *StockNotifier) AddObserver(o StockObserver) { n.observers = append(n.observers, o) }
func (n *StockNotifier) Notify(item *Item) {
	for _, o := range n.observers {
		o.Update(item)
	}
}

// Concrete observer example: Manager gets alerted

type ManagerStockObserver struct{ Name string }

func (m ManagerStockObserver) Update(item *Item) {
	fmt.Printf("[ALERT] Manager %s: Low stock for %s (SKU=%s, Qty=%d, ReorderLevel=%d)\n",
		m.Name, item.Name, item.SKU, item.Quantity, item.ReorderLevel)
}

// Helper to check and notify after a stock change
func checkAndNotifyReorder(it *Item, notifier *StockNotifier) {
	if notifier == nil {
		return
	}
	if it.IsReorderNeeded() {
		notifier.Notify(it)
	}
}

// =====================================================
// Factory Pattern (Shipments)
// =====================================================
// Shipment interface with two concrete types: Incoming and Outgoing.

const (
	ShipmentIncoming = "INCOMING"
	ShipmentOutgoing = "OUTGOING"
)

type ShipmentItem struct {
	Item *Item
	Qty  int
}

type Shipment interface {
	ID() string
	Carrier() string
	Type() string
	Status() string
	AddItem(item *Item, qty int)
	Items() []ShipmentItem
	Process(notifier *StockNotifier) error // apply stock changes and optionally notify
}

type baseShipment struct {
	shipmentID string
	carrier    string
	typeLabel  string
	status     string // created, processed
	items      []ShipmentItem
}

func (b *baseShipment) ID() string            { return b.shipmentID }
func (b *baseShipment) Carrier() string       { return b.carrier }
func (b *baseShipment) Type() string          { return b.typeLabel }
func (b *baseShipment) Status() string        { return b.status }
func (b *baseShipment) Items() []ShipmentItem { return b.items }
func (b *baseShipment) AddItem(it *Item, qty int) {
	b.items = append(b.items, ShipmentItem{Item: it, Qty: qty})
}

// IncomingShipment increases stock

type IncomingShipment struct{ baseShipment }

func NewIncomingShipment(id, carrier string) *IncomingShipment {
	return &IncomingShipment{baseShipment{shipmentID: id, carrier: carrier, typeLabel: ShipmentIncoming, status: "created"}}
}

func (s *IncomingShipment) Process(notifier *StockNotifier) error {
	fmt.Printf("Processing incoming shipment %s\n", s.ID())
	for _, si := range s.items {
		si.Item.UpdateStock(si.Qty)
		checkAndNotifyReorder(si.Item, notifier)
	}
	s.status = "received"
	return nil
}

// OutgoingShipment decreases stock

type OutgoingShipment struct{ baseShipment }

func NewOutgoingShipment(id, carrier string) *OutgoingShipment {
	return &OutgoingShipment{baseShipment{shipmentID: id, carrier: carrier, typeLabel: ShipmentOutgoing, status: "created"}}
}

func (s *OutgoingShipment) Process(notifier *StockNotifier) error {
	fmt.Printf("Processing outgoing shipment %s\n", s.ID())
	// Validate stock first
	for _, si := range s.items {
		if si.Item.Quantity < si.Qty {
			return fmt.Errorf("insufficient stock for %s (have %d, need %d)", si.Item.SKU, si.Item.Quantity, si.Qty)
		}
	}
	// Deduct
	for _, si := range s.items {
		si.Item.UpdateStock(-si.Qty)
		checkAndNotifyReorder(si.Item, notifier)
	}
	s.status = "shipped"
	return nil
}

// ShipmentFactory creates concrete shipments dynamically
func ShipmentFactory(typ, id, carrier string) (Shipment, error) {
	switch typ {
	case ShipmentIncoming:
		return NewIncomingShipment(id, carrier), nil
	case ShipmentOutgoing:
		return NewOutgoingShipment(id, carrier), nil
	default:
		return nil, errors.New("invalid shipment type")
	}
}

// -------------
// Demo (main)
// -------------
func main() {
	// Locations
	shelfA := NewStorageLocation("SHELF-A1", 100, "Shelf")
	binB := NewStorageLocation("BIN-B1", 50, "Bin")

	// Items
	laptop := NewItem("SKU101", "Laptop", 6, 5, 5)
	mouse := NewItem("SKU102", "Mouse", 12, 1, 10)
	_ = laptop.SetLocation(shelfA)
	_ = mouse.SetLocation(binB)

	// Users
	manager := &User{ID: "U2", Name: "Bob", Email: "bob@example.com", Role: RoleManager}
	_ = manager

	// Order
	order := NewOrder("ORD-001", "John Doe", manager)
	order.AddOrderItem(OrderItem{Item: laptop, Quantity: 2, Price: 1000})
	order.AddOrderItem(OrderItem{Item: mouse, Quantity: 3, Price: 25})
	fmt.Printf("Order %s total: %.2f (status: %s, managed by: %s)\n", order.OrderNumber, order.TotalCost(), order.Status, order.ManagedBy.Role)

	// --- Observer setup ---
	notifier := &StockNotifier{}
	notifier.AddObserver(ManagerStockObserver{Name: "Bob"})

	// --- Factory + Shipment processing ---
	// Outgoing shipment: ship 2 laptops and 5 mice
	outShip, _ := ShipmentFactory(ShipmentOutgoing, "SHIP-001", "DHL")
	outShip.AddItem(laptop, 2)
	outShip.AddItem(mouse, 5)
	if err := outShip.Process(notifier); err != nil {
		fmt.Println("Outgoing shipment error:", err)
	} else {
		fmt.Println("Outgoing shipment status:", outShip.Status())
	}

	// Incoming shipment: restock 3 laptops
	inShip, _ := ShipmentFactory(ShipmentIncoming, "SHIP-002", "UPS")
	inShip.AddItem(laptop, 3)
	_ = inShip.Process(notifier)
	fmt.Println("Incoming shipment status:", inShip.Status())

	// Manual check for reorder after operations
	if laptop.IsReorderNeeded() {
		notifier.Notify(laptop)
	}
	if mouse.IsReorderNeeded() {
		notifier.Notify(mouse)
	}
}
