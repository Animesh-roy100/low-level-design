package main

import (
	"math"
	"time"
)

// Placeholder for generateId function
func generateId() string {
	return "generated-id" // Implement actual ID generation
}

// Placeholder for isPeakHour function
func isPeakHour(t time.Time) bool {
	// Implement peak hour logic
	return false
}

// 1. User Management

type User struct {
	UserId        string
	Name          string
	Email         string
	Phone         string
	PaymentMethod string
}

type UserRepository interface {
	Save(user User)
	FindById(userId string) User
}

type UserService struct {
	userRepo UserRepository
}

func (us *UserService) RegisterUser(name, email, phone string) User {
	// Validation logic
	user := User{UserId: generateId(), Name: name, Email: email, Phone: phone}
	us.userRepo.Save(user)
	return user
}

// 2. Vehicle Management (Strategy Pattern)

type VehicleStatus string

const (
	Available   VehicleStatus = "AVAILABLE"
	Booked      VehicleStatus = "BOOKED"
	InService   VehicleStatus = "IN_SERVICE"
	Maintenance VehicleStatus = "MAINTENANCE"
)

type VehicleType struct {
	TypeId     string
	Name       string
	BaseFare   float64
	PerKmRate  float64
	PerMinRate float64
}

type Vehicle interface {
	CalculateFare(route Route) float64
	GetCurrentLocation() Location // Added for VehicleAllocator
}

type BaseVehicle struct {
	VehicleId       string
	LicensePlate    string
	Model           string
	Status          VehicleStatus
	Type            VehicleType
	CurrentLocation Location // Added for tracking and allocation
}

type AutonomousVehicle struct {
	BaseVehicle
	SoftwareVersion string
}

func (av *AutonomousVehicle) CalculateFare(route Route) float64 {
	return av.Type.BaseFare +
		(route.Distance * av.Type.PerKmRate) +
		(route.EstimatedDuration * av.Type.PerMinRate)
}

func (av *AutonomousVehicle) GetCurrentLocation() Location {
	return av.CurrentLocation
}

// 3. Booking System (State Pattern)

type BookingState interface {
	ConfirmBooking(booking *Booking)
	StartRide(booking *Booking)
	CompleteRide(booking *Booking)
	CancelBooking(booking *Booking)
}

type ConfirmedState struct{}

func (cs *ConfirmedState) ConfirmBooking(booking *Booking) {
	// Implement state transition (e.g., do nothing if already confirmed)
}

func (cs *ConfirmedState) StartRide(booking *Booking) {
	// Implement state transition to in-progress state
}

func (cs *ConfirmedState) CompleteRide(booking *Booking) {
	// Implement state transition
}

func (cs *ConfirmedState) CancelBooking(booking *Booking) {
	// Implement state transition to cancelled
}

type Booking struct {
	BookingId string
	User      User
	Vehicle   Vehicle
	Route     Route
	Payment   Payment
	State     BookingState
	StartTime time.Time
	EndTime   time.Time
}

func (b *Booking) StartRide() {
	b.State.StartRide(b)
}

// Other state transition methods can be added similarly

// 4. Pricing (Strategy Pattern)

type PricingStrategy interface {
	CalculateFare(route Route, typ VehicleType) float64
}

type StandardPricing struct{}

func (sp *StandardPricing) CalculateFare(route Route, typ VehicleType) float64 {
	return typ.BaseFare +
		(route.Distance * typ.PerKmRate)
}

type PeakHourPricing struct{}

func (php *PeakHourPricing) CalculateFare(route Route, typ VehicleType) float64 {
	base := typ.BaseFare * 1.2
	return base + (route.Distance * typ.PerKmRate * 1.5)
}

func GetStrategy(t time.Time) PricingStrategy {
	if isPeakHour(t) {
		return &PeakHourPricing{}
	}
	return &StandardPricing{}
}

// 5. Route Management

type Location struct {
	Latitude  float64
	Longitude float64
}

func (l Location) DistanceTo(other Location) float64 {
	// Haversine formula implementation
	const R = 6371.0 // Earth radius in km

	lat1 := l.Latitude * math.Pi / 180
	lat2 := other.Latitude * math.Pi / 180
	dlat := lat2 - lat1
	dlon := (other.Longitude - l.Longitude) * math.Pi / 180

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

type Route struct {
	RouteId           string
	Start             Location
	End               Location
	Distance          float64 // in km
	EstimatedDuration float64 // in minutes
}

func (r *Route) CalculateRoute(start, end Location) {
	r.Start = start
	r.End = end
	r.Distance = start.DistanceTo(end)
	// Calculate estimated duration based on traffic (placeholder)
	r.EstimatedDuration = r.Distance / 50 * 60 // Assume average speed 50 km/h
}

// 6. Payment Processing (Bridge Pattern)

type Payment struct {
	// Placeholder fields for payment details
}

type PaymentProcessor interface {
	ProcessPayment(amount float64, paymentMethod string) Payment
}

type CreditCardProcessor struct{}

func (ccp *CreditCardProcessor) ProcessPayment(amount float64, paymentMethod string) Payment {
	// Integrate with payment gateway (placeholder)
	return Payment{}
}

type PaymentService struct {
	processor PaymentProcessor
}

func NewPaymentService(processor PaymentProcessor) *PaymentService {
	return &PaymentService{processor: processor}
}

func (ps *PaymentService) CreatePayment(booking *Booking) Payment {
	amount := booking.Vehicle.CalculateFare(booking.Route)
	return ps.processor.ProcessPayment(amount, booking.User.PaymentMethod)
}

// 7. Vehicle Allocation

type VehicleRepository interface {
	FindAvailableVehicles(typ VehicleType) []Vehicle
}

type NoVehicleAvailableException struct{}

func (e *NoVehicleAvailableException) Error() string {
	return "no vehicle available"
}

type VehicleAllocator struct {
	vehicleRepo VehicleRepository
}

func (va *VehicleAllocator) FindNearestAvailable(pickup Location, typ VehicleType) Vehicle {
	available := va.vehicleRepo.FindAvailableVehicles(typ)
	if len(available) == 0 {
		panic(&NoVehicleAvailableException{})
	}

	var minVehicle Vehicle
	minDist := math.MaxFloat64
	for _, v := range available {
		dist := v.GetCurrentLocation().DistanceTo(pickup)
		if dist < minDist {
			minDist = dist
			minVehicle = v
		}
	}
	return minVehicle
}

// 8. Tracking System (Observer Pattern)

type TrackingObserver interface {
	UpdateVehicleLocation(vehicleId string, location Location)
}

type VehicleTracker struct {
	observers map[string][]TrackingObserver
}

func NewVehicleTracker() *VehicleTracker {
	return &VehicleTracker{observers: make(map[string][]TrackingObserver)}
}

func (vt *VehicleTracker) UpdateLocation(vehicleId string, location Location) {
	vt.notifyObservers(vehicleId, location)
}

func (vt *VehicleTracker) RegisterObserver(vehicleId string, observer TrackingObserver) {
	vt.observers[vehicleId] = append(vt.observers[vehicleId], observer)
}

func (vt *VehicleTracker) notifyObservers(vehicleId string, location Location) {
	for _, obs := range vt.observers[vehicleId] {
		obs.UpdateVehicleLocation(vehicleId, location)
	}
}

type BookingService struct {
	// Other fields
}

func (bs *BookingService) UpdateVehicleLocation(vehicleId string, location Location) {
	// Update UI with vehicle position (placeholder)
}

// For Singleton Pattern in service classes, Go typically uses package-level variables or dependency injection.
// Example for VehicleAllocator singleton (if needed):
var vehicleAllocator *VehicleAllocator

func GetVehicleAllocator(repo VehicleRepository) *VehicleAllocator {
	if vehicleAllocator == nil {
		vehicleAllocator = &VehicleAllocator{vehicleRepo: repo}
	}
	return vehicleAllocator
}
