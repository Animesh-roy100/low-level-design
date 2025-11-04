package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

// =====================
// Domain Models (Simplified with essential fields)
// =====================

type SeatStatus string

const (
	SeatAvailable SeatStatus = "AVAILABLE"
	SeatReserved  SeatStatus = "RESERVED"
	SeatBooked    SeatStatus = "BOOKED"
)

type Seat struct {
	SeatID        string
	EventID       int64
	SeatNumber    string
	Status        SeatStatus
	Price         int64 // in minor units
	ReservedBy    *string
	ReservedUntil *time.Time
	BookingID     *string
}

type ReservationStatus string

const (
	ReservationActive    ReservationStatus = "ACTIVE"
	ReservationExpired   ReservationStatus = "EXPIRED"
	ReservationConfirmed ReservationStatus = "CONFIRMED"
)

type Reservation struct {
	ReservationID string
	SeatID        string
	EventID       int64
	UserID        string
	ExpiresAt     time.Time
	Status        ReservationStatus
}

type Booking struct {
	BookingID        string
	EventID          int64
	UserID           string
	TotalAmount      int64
	Status           string // CONFIRMED, CANCELLED
	PaymentID        string
	PaymentStatus    string // SUCCESS, FAILED, REFUNDED
	BookingReference string
	ConfirmedAt      *time.Time
}

type BookingSeat struct {
	BookingID string
	SeatID    string
	Price     int64
}

type Event struct {
	EventID        int64
	AvailableSeats int
	Version        int64
}

// =====================
// Errors (Simplified)
// =====================

var (
	ErrSeatNotFound        = errors.New("seat not found")
	ErrSeatNotAvailable    = errors.New("seat not available")
	ErrReservationNotFound = errors.New("reservation not found")
	ErrReservationExpired  = errors.New("reservation expired")
	ErrInvalidReservation  = errors.New("invalid reservation state")
	ErrInvalidSeatState    = errors.New("seat state changed")
	ErrPaymentFailed       = errors.New("payment failed")
)

// =====================
// Repository Interfaces (Repository Pattern)
// =====================

type SeatRepository interface {
	FindByEventAndNumber(ctx context.Context, eventID int64, seatNumber string) (*Seat, error)
	FindByID(ctx context.Context, seatID string) (*Seat, error)
	Save(ctx context.Context, seat *Seat) error
}

type ReservationRepository interface {
	FindByID(ctx context.Context, reservationID string) (*Reservation, error)
	Save(ctx context.Context, r *Reservation) error
	FindExpired(ctx context.Context, now time.Time) ([]*Reservation, error)
}

type BookingRepository interface {
	FindByID(ctx context.Context, bookingID string) (*Booking, error)
	Save(ctx context.Context, b *Booking) error
}

type BookingSeatRepository interface {
	Save(ctx context.Context, bs *BookingSeat) error
}

type EventRepository interface {
	FindByID(ctx context.Context, eventID int64) (*Event, error)
	UpdateAvailableSeatsCAS(ctx context.Context, eventID int64, expectedVersion int64, newCount int) (updated bool, err error)
}

// =====================
// LockManager and IdempotencyStore Interfaces
// =====================

type LockManager interface {
	TryLock(ctx context.Context, key, value string, ttl time.Duration) (bool, error)
	Unlock(ctx context.Context, key, value string) error
}

type IdempotencyStore interface {
	Get(ctx context.Context, key string) (string, bool, error)
	SetNX(ctx context.Context, key, val string, ttl time.Duration) (bool, error)
}

// =====================
// PaymentService Interface
// =====================

type PaymentRequest struct {
	Amount int64
	UserID string
}

type PaymentResponse struct {
	Success      bool
	PaymentID    string
	ErrorMessage string
}

type PaymentService interface {
	Process(ctx context.Context, req PaymentRequest) (PaymentResponse, error)
	Refund(ctx context.Context, paymentID string) error
}

// =====================
// In-Memory Implementations (Simplified for demo)
// =====================

// InMemoryLockManager (Simulates distributed lock)
type InMemoryLockManager struct {
	mu    sync.Mutex
	locks map[string]struct {
		val string
		exp time.Time
	}
}

func NewInMemoryLockManager() *InMemoryLockManager {
	return &InMemoryLockManager{locks: make(map[string]struct {
		val string
		exp time.Time
	})}
}

func (m *InMemoryLockManager) TryLock(_ context.Context, key, value string, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for k, v := range m.locks {
		if now.After(v.exp) {
			delete(m.locks, k)
		}
	}
	if _, ok := m.locks[key]; ok {
		return false, nil
	}
	m.locks[key] = struct {
		val string
		exp time.Time
	}{value, now.Add(ttl)}
	return true, nil
}

func (m *InMemoryLockManager) Unlock(_ context.Context, key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.locks[key]; ok && v.val == value {
		delete(m.locks, key)
	}
	return nil
}

// InMemoryIdempotencyStore
type InMemoryIdempotencyStore struct {
	mu   sync.Mutex
	data map[string]struct {
		val string
		exp time.Time
	}
}

func NewInMemoryIdempotencyStore() *InMemoryIdempotencyStore {
	return &InMemoryIdempotencyStore{data: make(map[string]struct {
		val string
		exp time.Time
	})}
}

func (s *InMemoryIdempotencyStore) Get(_ context.Context, key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for k, v := range s.data {
		if now.After(v.exp) {
			delete(s.data, k)
		}
	}
	if v, ok := s.data[key]; ok {
		return v.val, true, nil
	}
	return "", false, nil
}

func (s *InMemoryIdempotencyStore) SetNX(_ context.Context, key, val string, ttl time.Duration) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[key]; ok {
		return false, nil
	}
	s.data[key] = struct {
		val string
		exp time.Time
	}{val, time.Now().Add(ttl)}
	return true, nil
}

// MockPaymentService
type MockPaymentService struct{}

func (m MockPaymentService) Process(_ context.Context, req PaymentRequest) (PaymentResponse, error) {
	if rand.Float64() < 0.1 { // 10% failure rate
		return PaymentResponse{Success: false, ErrorMessage: "payment failed"}, ErrPaymentFailed
	}
	return PaymentResponse{Success: true, PaymentID: fmt.Sprintf("pay_%d", rand.Int63())}, nil
}

func (m MockPaymentService) Refund(_ context.Context, paymentID string) error {
	log.Printf("Refunded %s", paymentID)
	return nil
}

// InMemorySeatRepository
type InMemorySeatRepository struct {
	mu    sync.Mutex
	seats map[string]*Seat  // by SeatID
	byKey map[string]string // by "EventID:SeatNumber" -> SeatID
}

func NewInMemorySeatRepository() *InMemorySeatRepository {
	return &InMemorySeatRepository{
		seats: make(map[string]*Seat),
		byKey: make(map[string]string),
	}
}

func (r *InMemorySeatRepository) FindByEventAndNumber(_ context.Context, eventID int64, seatNumber string) (*Seat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%d:%s", eventID, seatNumber)
	if id, ok := r.byKey[key]; ok {
		if seat, ok := r.seats[id]; ok {
			return copySeat(seat), nil
		}
	}
	return nil, ErrSeatNotFound
}

func (r *InMemorySeatRepository) FindByID(_ context.Context, seatID string) (*Seat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if seat, ok := r.seats[seatID]; ok {
		return copySeat(seat), nil
	}
	return nil, ErrSeatNotFound
}

func (r *InMemorySeatRepository) Save(_ context.Context, seat *Seat) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if seat.SeatID == "" {
		seat.SeatID = fmt.Sprintf("seat_%d", rand.Int63())
	}
	r.seats[seat.SeatID] = copySeat(seat)
	key := fmt.Sprintf("%d:%s", seat.EventID, seat.SeatNumber)
	r.byKey[key] = seat.SeatID
	return nil
}

func copySeat(s *Seat) *Seat {
	cp := *s
	if s.ReservedBy != nil {
		rb := *s.ReservedBy
		cp.ReservedBy = &rb
	}
	if s.ReservedUntil != nil {
		ru := *s.ReservedUntil
		cp.ReservedUntil = &ru
	}
	if s.BookingID != nil {
		bi := *s.BookingID
		cp.BookingID = &bi
	}
	return &cp
}

// Similar simplified in-memory repos for others...

type InMemoryReservationRepository struct {
	mu           sync.Mutex
	reservations map[string]*Reservation
}

func NewInMemoryReservationRepository() *InMemoryReservationRepository {
	return &InMemoryReservationRepository{reservations: make(map[string]*Reservation)}
}

func (r *InMemoryReservationRepository) FindByID(_ context.Context, id string) (*Reservation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if res, ok := r.reservations[id]; ok {
		cp := *res
		return &cp, nil
	}
	return nil, ErrReservationNotFound
}

func (r *InMemoryReservationRepository) Save(_ context.Context, res *Reservation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if res.ReservationID == "" {
		res.ReservationID = fmt.Sprintf("res_%d", rand.Int63())
	}
	cp := *res
	r.reservations[res.ReservationID] = &cp
	return nil
}

func (r *InMemoryReservationRepository) FindExpired(_ context.Context, now time.Time) ([]*Reservation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var expired []*Reservation
	for _, res := range r.reservations {
		if res.Status == ReservationActive && now.After(res.ExpiresAt) {
			cp := *res
			expired = append(expired, &cp)
		}
	}
	return expired, nil
}

type InMemoryBookingRepository struct {
	mu       sync.Mutex
	bookings map[string]*Booking
}

func NewInMemoryBookingRepository() *InMemoryBookingRepository {
	return &InMemoryBookingRepository{bookings: make(map[string]*Booking)}
}

func (r *InMemoryBookingRepository) FindByID(_ context.Context, id string) (*Booking, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if b, ok := r.bookings[id]; ok {
		cp := *b
		return &cp, nil
	}
	return nil, errors.New("booking not found")
}

func (r *InMemoryBookingRepository) Save(_ context.Context, b *Booking) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if b.BookingID == "" {
		b.BookingID = fmt.Sprintf("book_%d", rand.Int63())
	}
	cp := *b
	r.bookings[b.BookingID] = &cp
	return nil
}

type InMemoryBookingSeatRepository struct {
	mu           sync.Mutex
	bookingSeats []*BookingSeat
}

func NewInMemoryBookingSeatRepository() *InMemoryBookingSeatRepository {
	return &InMemoryBookingSeatRepository{}
}

func (r *InMemoryBookingSeatRepository) Save(_ context.Context, bs *BookingSeat) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *bs
	r.bookingSeats = append(r.bookingSeats, &cp)
	return nil
}

type InMemoryEventRepository struct {
	mu     sync.Mutex
	events map[int64]*Event
}

func NewInMemoryEventRepository() *InMemoryEventRepository {
	return &InMemoryEventRepository{events: make(map[int64]*Event)}
}

func (r *InMemoryEventRepository) FindByID(_ context.Context, id int64) (*Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e, ok := r.events[id]; ok {
		cp := *e
		return &cp, nil
	}
	return nil, errors.New("event not found")
}

func (r *InMemoryEventRepository) UpdateAvailableSeatsCAS(_ context.Context, eventID int64, expectedVersion int64, newCount int) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e, ok := r.events[eventID]; ok {
		if e.Version != expectedVersion {
			return false, nil
		}
		e.AvailableSeats = newCount
		e.Version++
		return true, nil
	}
	return false, errors.New("event not found")
}

// =====================
// TicketBookingService (Service Layer with Patterns)
// =====================
// Uses:
// - Repository Pattern for data access abstraction
// - Distributed Locking for concurrency (via LockManager)
// - Idempotency Pattern for safe retries
// - Saga Pattern for handling payment + booking with compensation
// - Optimistic Concurrency Control for event updates

type TicketBookingService struct {
	seatRepo           SeatRepository
	reservationRepo    ReservationRepository
	bookingRepo        BookingRepository
	bookingSeatRepo    BookingSeatRepository
	eventRepo          EventRepository
	lockManager        LockManager
	paymentService     PaymentService
	idempotencyStore   IdempotencyStore
	reservationTimeout time.Duration
}

func NewTicketBookingService(
	seatRepo SeatRepository,
	reservationRepo ReservationRepository,
	bookingRepo BookingRepository,
	bookingSeatRepo BookingSeatRepository,
	eventRepo EventRepository,
	lockManager LockManager,
	paymentService PaymentService,
	idempotencyStore IdempotencyStore,
	reservationTimeout time.Duration,
) *TicketBookingService {
	if reservationTimeout == 0 {
		reservationTimeout = 10 * time.Minute
	}
	return &TicketBookingService{
		seatRepo:           seatRepo,
		reservationRepo:    reservationRepo,
		bookingRepo:        bookingRepo,
		bookingSeatRepo:    bookingSeatRepo,
		eventRepo:          eventRepo,
		lockManager:        lockManager,
		paymentService:     paymentService,
		idempotencyStore:   idempotencyStore,
		reservationTimeout: reservationTimeout,
	}
}

const lockPrefix = "seat:lock:"

// ReserveSeats (Supports multiple seats with sorted locking to avoid deadlocks)
func (s *TicketBookingService) ReserveSeats(ctx context.Context, eventID int64, seatNumbers []string, userID string) ([]*Reservation, time.Time, error) {
	if len(seatNumbers) == 0 {
		return nil, time.Time{}, errors.New("no seats specified")
	}
	sortedSeats := make([]string, len(seatNumbers))
	copy(sortedSeats, seatNumbers)
	sort.Strings(sortedSeats)

	var acquiredLocks []struct{ key, val string }
	defer s.releaseLocks(ctx, acquiredLocks)

	for _, seatNum := range sortedSeats {
		lockKey := fmt.Sprintf("%s%d:%s", lockPrefix, eventID, seatNum)
		lockVal := randString(16)
		ok, err := s.lockManager.TryLock(ctx, lockKey, lockVal, 30*time.Second)
		if err != nil {
			return nil, time.Time{}, err
		}
		if !ok {
			return nil, time.Time{}, ErrSeatNotAvailable
		}
		acquiredLocks = append(acquiredLocks, struct{ key, val string }{lockKey, lockVal})
	}

	reservedUntil := time.Now().Add(s.reservationTimeout)
	var reservations []*Reservation
	for _, seatNum := range sortedSeats {
		seat, err := s.seatRepo.FindByEventAndNumber(ctx, eventID, seatNum)
		if err != nil {
			return nil, time.Time{}, err
		}
		if seat.Status != SeatAvailable {
			return nil, time.Time{}, ErrSeatNotAvailable
		}

		seat.Status = SeatReserved
		seat.ReservedBy = &userID
		seat.ReservedUntil = &reservedUntil
		if err := s.seatRepo.Save(ctx, seat); err != nil {
			return nil, time.Time{}, err
		}

		res := &Reservation{
			SeatID:    seat.SeatID,
			EventID:   eventID,
			UserID:    userID,
			ExpiresAt: reservedUntil,
			Status:    ReservationActive,
		}
		if err := s.reservationRepo.Save(ctx, res); err != nil {
			return nil, time.Time{}, err
		}
		reservations = append(reservations, res)
	}

	return reservations, reservedUntil, nil
}

// ConfirmBooking (For single reservation; extend for multi)
func (s *TicketBookingService) ConfirmBooking(ctx context.Context, reservationID string, payReq PaymentRequest) (*Booking, error) {
	res, err := s.reservationRepo.FindByID(ctx, reservationID)
	if err != nil {
		return nil, err
	}
	if time.Now().After(res.ExpiresAt) {
		return nil, ErrReservationExpired
	}
	if res.Status != ReservationActive {
		return nil, ErrInvalidReservation
	}

	seat, err := s.seatRepo.FindByID(ctx, res.SeatID)
	if err != nil {
		return nil, err
	}

	lockKey := fmt.Sprintf("%s%d:%s", lockPrefix, res.EventID, seat.SeatNumber)
	lockVal := randString(16)
	ok, err := s.lockManager.TryLock(ctx, lockKey, lockVal, 30*time.Second)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrSeatNotAvailable
	}
	defer s.lockManager.Unlock(ctx, lockKey, lockVal)

	payResp, err := s.paymentService.Process(ctx, payReq)
	if err != nil {
		return nil, err
	}
	if !payResp.Success {
		return nil, fmt.Errorf("payment failed: %s", payResp.ErrorMessage)
	}

	// Saga: If subsequent steps fail, refund
	defer func() {
		if err != nil {
			s.paymentService.Refund(ctx, payResp.PaymentID)
		}
	}()

	if seat.Status != SeatReserved || seat.ReservedBy == nil || *seat.ReservedBy != res.UserID {
		return nil, ErrInvalidSeatState
	}

	now := time.Now()
	booking := &Booking{
		EventID:          res.EventID,
		UserID:           res.UserID,
		TotalAmount:      seat.Price,
		Status:           "CONFIRMED",
		PaymentID:        payResp.PaymentID,
		PaymentStatus:    "SUCCESS",
		BookingReference: generateBookingReference(),
		ConfirmedAt:      &now,
	}
	if err = s.bookingRepo.Save(ctx, booking); err != nil {
		return nil, err
	}

	bs := &BookingSeat{
		BookingID: booking.BookingID,
		SeatID:    seat.SeatID,
		Price:     seat.Price,
	}
	if err = s.bookingSeatRepo.Save(ctx, bs); err != nil {
		return nil, err
	}

	seat.Status = SeatBooked
	seat.BookingID = &booking.BookingID
	seat.ReservedBy = nil
	seat.ReservedUntil = nil
	if err = s.seatRepo.Save(ctx, seat); err != nil {
		return nil, err
	}

	res.Status = ReservationConfirmed
	if err = s.reservationRepo.Save(ctx, res); err != nil {
		return nil, err
	}

	s.updateEventAvailableSeats(ctx, res.EventID, -1)

	return booking, nil
}

// ConfirmBookingWithIdempotency (Idempotency Pattern)
func (s *TicketBookingService) ConfirmBookingWithIdempotency(ctx context.Context, idempotencyKey, reservationID string, payReq PaymentRequest) (*Booking, error) {
	if idempotencyKey == "" {
		return s.ConfirmBooking(ctx, reservationID, payReq)
	}

	key := "idempotency:" + idempotencyKey
	val, found, err := s.idempotencyStore.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if found {
		return s.bookingRepo.FindByID(ctx, val)
	}

	booking, err := s.ConfirmBooking(ctx, reservationID, payReq)
	if err != nil {
		return nil, err
	}

	_, err = s.idempotencyStore.SetNX(ctx, key, booking.BookingID, 24*time.Hour)
	if err != nil {
		// Log but continue
		log.Printf("Failed to set idempotency key: %v", err)
	}

	return booking, nil
}

// CleanupExpiredReservations (Scheduled task)
func (s *TicketBookingService) CleanupExpiredReservations(ctx context.Context) {
	expired, err := s.reservationRepo.FindExpired(ctx, time.Now())
	if err != nil {
		log.Printf("Error finding expired reservations: %v", err)
		return
	}

	for _, res := range expired {
		s.releaseExpiredReservation(ctx, res)
	}
}

func (s *TicketBookingService) releaseExpiredReservation(ctx context.Context, res *Reservation) {
	seat, err := s.seatRepo.FindByID(ctx, res.SeatID)
	if err != nil {
		return
	}

	lockKey := fmt.Sprintf("%s%d:%s", lockPrefix, res.EventID, seat.SeatNumber)
	lockVal := randString(16)
	ok, _ := s.lockManager.TryLock(ctx, lockKey, lockVal, 5*time.Second)
	if !ok {
		return
	}
	defer s.lockManager.Unlock(ctx, lockKey, lockVal)

	if seat.Status == SeatReserved && seat.ReservedUntil != nil && time.Now().After(*seat.ReservedUntil) {
		seat.Status = SeatAvailable
		seat.ReservedBy = nil
		seat.ReservedUntil = nil
		s.seatRepo.Save(ctx, seat)

		res.Status = ReservationExpired
		s.reservationRepo.Save(ctx, res)

		s.updateEventAvailableSeats(ctx, res.EventID, 1)
	}
}

func (s *TicketBookingService) releaseLocks(ctx context.Context, locks []struct{ key, val string }) {
	for _, l := range locks {
		s.lockManager.Unlock(ctx, l.key, l.val)
	}
}

func (s *TicketBookingService) updateEventAvailableSeats(ctx context.Context, eventID int64, delta int) {
	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {
		event, err := s.eventRepo.FindByID(ctx, eventID)
		if err != nil {
			log.Printf("Error finding event: %v", err)
			return
		}
		newCount := event.AvailableSeats + delta
		updated, err := s.eventRepo.UpdateAvailableSeatsCAS(ctx, eventID, event.Version, newCount)
		if err != nil {
			log.Printf("Error updating event: %v", err)
			return
		}
		if updated {
			return
		}
	}
	log.Println("Failed to update event available seats after retries")
}

// Utilities
func randString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateBookingReference() string {
	return strings.ToUpper(fmt.Sprintf("BK-%d", rand.Intn(1000000)))
}

// =====================
// Main Demo
// =====================

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()

	// Initialize repositories and services
	seatRepo := NewInMemorySeatRepository()
	reservationRepo := NewInMemoryReservationRepository()
	bookingRepo := NewInMemoryBookingRepository()
	bookingSeatRepo := NewInMemoryBookingSeatRepository()
	eventRepo := NewInMemoryEventRepository()
	lockManager := NewInMemoryLockManager()
	idempotencyStore := NewInMemoryIdempotencyStore()
	paymentService := MockPaymentService{}

	// Seed data
	eventRepo.events[1] = &Event{EventID: 1, AvailableSeats: 2, Version: 0}
	seatRepo.Save(ctx, &Seat{EventID: 1, SeatNumber: "A1", Status: SeatAvailable, Price: 10000})
	seatRepo.Save(ctx, &Seat{EventID: 1, SeatNumber: "A2", Status: SeatAvailable, Price: 10000})

	service := NewTicketBookingService(
		seatRepo,
		reservationRepo,
		bookingRepo,
		bookingSeatRepo,
		eventRepo,
		lockManager,
		paymentService,
		idempotencyStore,
		10*time.Minute,
	)

	// Reserve
	reservations, until, err := service.ReserveSeats(ctx, 1, []string{"A1", "A2"}, "user123")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Reserved until %v: %+v", until, reservations)

	// Confirm one
	booking, err := service.ConfirmBookingWithIdempotency(ctx, "key123", reservations[0].ReservationID, PaymentRequest{Amount: 10000, UserID: "user123"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Booked: %+v", booking)

	// Cleanup (simulate)
	service.CleanupExpiredReservations(ctx)
}
