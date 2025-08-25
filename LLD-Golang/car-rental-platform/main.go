package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Customers can do:
// 1. Search available cars
// 2. Reserver a car
// 3. Cancel a reservation
// 4. Manage billing & payments
// 5. Track rental history

// Requirements:
// 1. User should be able to book/cancel cars
// 2. Cars should be managed by stores with their own inventory
// 3. Payments should support different modes (e.g. credit, wallet)
// 4. The system should support different vehicle types
// 5. Notification should be sent on booking/cancellation

// Entity Classes:
// 1. User → represents customers with basic details (id, name, licence)
// 2. Car/Vehicle → Abstract class for different vehicle types (SUV, Sedan etc)
// 3. VehicleInventoryManagement → Manages availability of cars in a store
// 4. Store → Represents a physical location with cars & reservations
// 5. ReservationManager → Handles lifecycle of reservations (create, update, cancel)
// 6. Bill → Calculates rental charges
// 7. Payment/PaymentService → Handles different payment modes
// 8. NotificationService → Sends booking confirmation/cancellation alerts.

// Design Patterns Used:
// 1. Factory Pattern → Used for creating different vehicle objects dynamically (SUV, Sedan etc)
// 2. Strategy Pattern → Used in Bill/Pricing for flexible billing strategies (HOURLY, DAILY, WEEKLY)
// 3. Observer Pattern → NotificationService observes reservation events (confirmation, cancellation) and triggers notifications
// 4. Singleton Pattern → ReservationManager can be singleton to avoid duplicate booking handlings

// Enums (using string constants for simplicity)
type ReservationStatus string

const (
	Initiated  ReservationStatus = "INITIATED"
	Scheduled  ReservationStatus = "SCHEDULED"
	InProgress ReservationStatus = "INPROGRESS"
	Completed  ReservationStatus = "COMPLETED"
	Cancelled  ReservationStatus = "CANCELLED"
)

type PaymentMode string

const (
	Online PaymentMode = "ONLINE"
	Cash   PaymentMode = "CASH"
)

type ReservationType string

const (
	Hourly ReservationType = "HOURLY"
	Daily  ReservationType = "DAILY"
)

type CarType string

const (
	Minivan   CarType = "MINIVAN"
	SUV       CarType = "SUV"
	Sedan     CarType = "SEDAN"
	Sport     CarType = "SPORT"
	Hatchback CarType = "HATCHBACK"
)

// Location
type Location struct {
	Address string
	Pincode int
	City    string
	State   string
	Country string
}

// User
type User struct {
	UserID         int
	UserName       string
	DrivingLicense string
}

// Vehicle (merged with Car; no need for separate struct)
type Vehicle struct {
	ID          int
	Make        string
	Model       string
	Year        int
	PricePerDay float64
	NumberPlate string
	Type        CarType
}

// Reservation
type Reservation struct {
	ReservationID     int
	User              *User
	Vehicle           *Vehicle
	StartTime         time.Time
	EndTime           time.Time
	Status            ReservationStatus
	Location          *Location
	ReservationType   ReservationType // Added for billing
	Bill              *Bill           // Link to Bill
	reservationsMutex sync.Mutex      // For internal overlap check
}

// Check if this reservation overlaps with any existing (dummy; in real, check against all for vehicle)
func (r *Reservation) Overlaps(start, end time.Time) bool {
	return !(r.EndTime.Before(start) || r.StartTime.After(end))
}

// Bill
type Bill struct {
	Reservation    *Reservation
	TotalAmount    float64
	IsPaid         bool
	PaymentDetails *PaymentDetails // Added link
}

func (b *Bill) ComputeBillAmount() float64 {
	duration := b.Reservation.EndTime.Sub(b.Reservation.StartTime)
	var rate float64
	switch b.Reservation.ReservationType {
	case Hourly:
		rate = b.Reservation.Vehicle.PricePerDay / 24 // Derive hourly from daily
		b.TotalAmount = rate * duration.Hours()
	case Daily:
		rate = b.Reservation.Vehicle.PricePerDay
		b.TotalAmount = rate * (duration.Hours()/24 + 1) // Ceiling days
	}
	return b.TotalAmount
}

// PaymentDetails
type PaymentDetails struct {
	PaymentID     int
	AmountPaid    float64
	DateOfPayment time.Time
	IsRefundable  bool
	PaymentMode   PaymentMode
}

// Payment (struct with method)
type Payment struct{}

func (p *Payment) PayBill(bill *Bill, mode PaymentMode) error {
	if bill.IsPaid {
		return errors.New("bill already paid")
	}
	bill.ComputeBillAmount() // Ensure computed
	bill.PaymentDetails = &PaymentDetails{
		PaymentID:     1, // Dummy ID
		AmountPaid:    bill.TotalAmount,
		DateOfPayment: time.Now(),
		IsRefundable:  true, // Default
		PaymentMode:   mode,
	}
	bill.IsPaid = true
	fmt.Printf("Paid bill of $%.2f via %s\n", bill.TotalAmount, mode)
	return nil
}

// VehicleInventoryManagement
type VehicleInventoryManagement struct {
	Vehicles []Vehicle
	mutex    sync.Mutex
}

func (vim *VehicleInventoryManagement) GetVehicles() []Vehicle {
	return vim.Vehicles
}

func (vim *VehicleInventoryManagement) AddVehicles(vehicles ...Vehicle) {
	vim.mutex.Lock()
	defer vim.mutex.Unlock()
	vim.Vehicles = append(vim.Vehicles, vehicles...)
}

// ReservationManager
type ReservationManager struct {
	Reservations []Reservation
	mutex        sync.Mutex
	counter      int
}

func (rm *ReservationManager) CreateReservation(res *Reservation) (*Reservation, error) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Check availability (simple overlap check on all reservations for this vehicle)
	for _, existing := range rm.Reservations {
		if existing.Vehicle.ID == res.Vehicle.ID && existing.Status != Cancelled && existing.Overlaps(res.StartTime, res.EndTime) {
			return nil, errors.New("vehicle not available")
		}
	}

	rm.counter++
	res.ReservationID = rm.counter
	res.Status = Initiated
	rm.Reservations = append(rm.Reservations, *res)
	return res, nil
}

func (rm *ReservationManager) ChangeStatusToScheduled(id int) error {
	return rm.updateStatus(id, Scheduled)
}

func (rm *ReservationManager) ChangeStatusToInProgress(id int) error {
	return rm.updateStatus(id, InProgress)
}

func (rm *ReservationManager) CompleteReservation(id int) error {
	return rm.updateStatus(id, Completed)
}

func (rm *ReservationManager) CancelReservation(id int) error {
	return rm.updateStatus(id, Cancelled)
}

func (rm *ReservationManager) updateStatus(id int, status ReservationStatus) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	for i := range rm.Reservations {
		if rm.Reservations[i].ReservationID == id {
			rm.Reservations[i].Status = status
			return nil
		}
	}
	return errors.New("reservation not found")
}

// Store
type Store struct {
	StoreID            int
	Location           Location
	InventoryManager   VehicleInventoryManagement
	ReservationManager ReservationManager
}

func (s *Store) GetVehiclesByType(t CarType) []Vehicle {
	var result []Vehicle
	for _, v := range s.InventoryManager.Vehicles {
		if v.Type == t {
			result = append(result, v)
		}
	}
	return result
}

func (s *Store) SearchAvailableVehicles(t CarType, start, end time.Time) []Vehicle {
	var available []Vehicle
	for _, v := range s.GetVehiclesByType(t) {
		avl := true
		for _, res := range s.ReservationManager.Reservations {
			if res.Vehicle.ID == v.ID && res.Status != Cancelled && res.Overlaps(start, end) {
				avl = false
				break
			}
		}
		if avl {
			available = append(available, v)
		}
	}
	return available
}

func (s *Store) UpdateOrCreateReservation(res *Reservation) error {
	// If ID exists, update; else create
	if res.ReservationID > 0 {
		return s.ReservationManager.updateStatus(res.ReservationID, res.Status) // Example update
	}
	_, err := s.ReservationManager.CreateReservation(res)
	return err
}

// VehicleRentalSystem
type VehicleRentalSystem struct {
	StoreList []Store
	UserList  []User
	mutex     sync.Mutex
}

func (vrs *VehicleRentalSystem) GetStore(loc Location) *Store {
	for i := range vrs.StoreList {
		if vrs.StoreList[i].Location == loc {
			return &vrs.StoreList[i]
		}
	}
	return nil
}

func (vrs *VehicleRentalSystem) AddStore(store Store) {
	vrs.mutex.Lock()
	defer vrs.mutex.Unlock()
	vrs.StoreList = append(vrs.StoreList, store)
}

func (vrs *VehicleRentalSystem) AddUser(user User) {
	vrs.mutex.Lock()
	defer vrs.mutex.Unlock()
	vrs.UserList = append(vrs.UserList, user)
}

// Demo
func main() {
	// Setup
	vrs := VehicleRentalSystem{}

	loc := Location{Address: "123 Main St", Pincode: 10001, City: "New York", State: "NY", Country: "USA"}
	store := Store{
		StoreID:            1,
		Location:           loc,
		InventoryManager:   VehicleInventoryManagement{},
		ReservationManager: ReservationManager{},
	}
	store.InventoryManager.AddVehicles(
		Vehicle{ID: 1, Make: "Toyota", Model: "Camry", Year: 2020, PricePerDay: 50, NumberPlate: "ABC123", Type: Sedan},
		Vehicle{ID: 2, Make: "Ford", Model: "Explorer", Year: 2021, PricePerDay: 80, NumberPlate: "DEF456", Type: SUV},
	)
	vrs.AddStore(store)

	user := User{UserID: 1, UserName: "John Doe", DrivingLicense: "LIC123"}
	vrs.AddUser(user)

	// 1. Search available vehicles
	start := time.Now().Add(24 * time.Hour)
	end := start.Add(48 * time.Hour)
	available := store.SearchAvailableVehicles(SUV, start, end)
	fmt.Printf("Available SUVs: %d\n", len(available)) // 1

	// 2. Create reservation
	res := &Reservation{
		User:            &user,
		Vehicle:         &store.InventoryManager.Vehicles[1], // SUV
		StartTime:       start,
		EndTime:         end,
		Location:        &loc,
		ReservationType: Daily,
		Bill:            &Bill{},
	}
	_, err := store.ReservationManager.CreateReservation(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	store.ReservationManager.ChangeStatusToScheduled(res.ReservationID)

	// 3. Compute and pay bill
	bill := res.Bill
	bill.Reservation = res
	amount := bill.ComputeBillAmount()
	fmt.Printf("Bill amount: $%.2f\n", amount)

	payment := Payment{}
	err = payment.PayBill(bill, Online)
	if err != nil {
		fmt.Println(err)
	}

	// 4. Cancel reservation
	err = store.ReservationManager.CancelReservation(res.ReservationID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Reservation status:", res.Status)
}
