package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

/*
At Urban Company we have both Customers and Partners making Payments on the platform.
Customers make payments for services they book and Partners for purchasing products used during Job delivery.

To handle these transactions in the system, we are supposed to set up a Payments Platform.
The platform needs to have intelligent logic to make sure the success rate of transactions at an overall level is always high.
Therefore, we have onboarded 3 Payment gateways (PayU, Paytm, RazorPay).
While initiating a Transaction, Payment system should choose the best Gateway at that time.

Users should be able to choose their preferred payment method for a transaction - Card, UPI, NetBanking.
It’s possible that not all gateways have support for each payment method.
For example PayU may only support Cards and NB. RazorPay may have support for all three etc.
Additionally it's possible that one gateway may be performing good for Cards recently but the NB performance has degraded.
In such a case our payment system should keep routing Cards but a different gateway for NB.

You are required to come up with a low level Design for the Payment system - initiating a transaction, updating status of transaction and choosing best Gateway
*/

type PaymentMethod string

const (
	Card       PaymentMethod = "CARD"
	UPI        PaymentMethod = "UPI"
	NetBanking PaymentMethod = "NETBANKING"
)

type UserType string

const (
	Customer UserType = "CUSTOMER"
	Partner  UserType = "PARTNER"
)

type TransactionStatus string

const (
	Pending    TransactionStatus = "PENDING"
	Initiated  TransactionStatus = "INITIATED"
	Success    TransactionStatus = "SUCCESS"
	Failed     TransactionStatus = "FAILED"
	Processing TransactionStatus = "PROCESSING"
)

// Core Domain Entities -----------------------------------------

type Transaction struct {
	ID                   string
	UserType             UserType
	Amount               float64
	PaymentMethod        PaymentMethod
	GatewayUsed          string
	GatewayTransactionID string
	Status               TransactionStatus
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// Gateway Interface (Strategy Pattern) --------------------------------------------

type Gateway interface {
	Initiate(tx *Transaction) (*GatewayResponse, error)
	GetStatus(gatewayTxID string) (*GatewayStatusResponse, error)
	SupportedMethods() []PaymentMethod
	Name() string
}

// GatewayResponse from initiation
type GatewayResponse struct {
	GatewayTxID string `json:"gateway_tx_id"`
	Status      string `json:"status"`
}

// GatewayStatusResponse
type GatewayStatusResponse struct {
	GatewayTxID string `json:"gateway_tx_id"`
	Status      string `json:"status"` // SUCCESS, FAILED, etc.
}

// Concrete Gateways ------------------------------------------------

type PayU struct{}

func (p *PayU) Initiate(tx *Transaction) (*GatewayResponse, error) {
	if !contains(p.SupportedMethods(), tx.PaymentMethod) {
		return nil, errors.New("unsupported payment method")
	}
	gatewayTxID := fmt.Sprintf("payu_%s", generateID())
	return &GatewayResponse{GatewayTxID: gatewayTxID, Status: "INITIATED"}, nil
}

func (p *PayU) GetStatus(gatewayTxID string) (*GatewayStatusResponse, error) {
	// Mock: 90% success for CARD, 80% for NB
	if rand.Float64() > 0.85 { // Simulate failure
		return &GatewayStatusResponse{GatewayTxID: gatewayTxID, Status: "SUCCESS"}, nil
	}
	return &GatewayStatusResponse{GatewayTxID: gatewayTxID, Status: "FAILED"}, nil
}

func (p *PayU) SupportedMethods() []PaymentMethod {
	return []PaymentMethod{Card, NetBanking}
}

func (p *PayU) Name() string { return "PayU" }

// ------------------------------------------------

type Paytm struct{}

func (pt *Paytm) Initiate(tx *Transaction) (*GatewayResponse, error) {
	if !contains(pt.SupportedMethods(), tx.PaymentMethod) {
		return nil, errors.New("unsupported payment method")
	}
	gatewayTxID := fmt.Sprintf("paytm_%s", generateID())
	return &GatewayResponse{GatewayTxID: gatewayTxID, Status: "INITIATED"}, nil
}

func (pt *Paytm) GetStatus(gatewayTxID string) (*GatewayStatusResponse, error) {
	// Mock: 85% success for all
	if rand.Float64() > 0.15 {
		return &GatewayStatusResponse{GatewayTxID: gatewayTxID, Status: "SUCCESS"}, nil
	}
	return &GatewayStatusResponse{GatewayTxID: gatewayTxID, Status: "FAILED"}, nil
}

func (pt *Paytm) SupportedMethods() []PaymentMethod {
	return []PaymentMethod{Card, UPI}
}

func (pt *Paytm) Name() string { return "Paytm" }

// ------------------------------------------------

type RazorPay struct{}

func (r *RazorPay) Initiate(tx *Transaction) (*GatewayResponse, error) {
	if !contains(r.SupportedMethods(), tx.PaymentMethod) {
		return nil, errors.New("unsupported payment method")
	}
	gatewayTxID := fmt.Sprintf("razor_%s", generateID())
	return &GatewayResponse{GatewayTxID: gatewayTxID, Status: "INITIATED"}, nil
}

func (r *RazorPay) GetStatus(gatewayTxID string) (*GatewayStatusResponse, error) {
	// Mock: 95% success for all
	if rand.Float64() > 0.05 {
		return &GatewayStatusResponse{GatewayTxID: gatewayTxID, Status: "SUCCESS"}, nil
	}
	return &GatewayStatusResponse{GatewayTxID: gatewayTxID, Status: "FAILED"}, nil
}

func (r *RazorPay) SupportedMethods() []PaymentMethod {
	return []PaymentMethod{Card, UPI, NetBanking}
}

func (r *RazorPay) Name() string { return "RazorPay" }

// Helper: generateID generates a unique ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Helper: contains
func contains(methods []PaymentMethod, method PaymentMethod) bool {
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}

// Gateway Factory (Factory Pattern) --------------------------------------------

type GatewayFactory struct{}

func NewGatewayFactory() *GatewayFactory {
	return &GatewayFactory{}
}

func (gf *GatewayFactory) CreateGateway(name string) Gateway {
	switch name {
	case "PayU":
		return &PayU{}
	case "Paytm":
		return &Paytm{}
	case "RazorPay":
		return &RazorPay{}
	default:
		return nil
	}
}

// Metrics and Gateway Selection Logic --------------------------------------------

type Metrics struct {
	successRates map[string]map[PaymentMethod]float64 // gateway -> method -> rate
	mu           sync.RWMutex
	totalTxs     map[string]map[PaymentMethod]int // for updating rates
	successTxs   map[string]map[PaymentMethod]int
}

func NewMetrics() *Metrics {
	return &Metrics{
		successRates: make(map[string]map[PaymentMethod]float64),
		totalTxs:     make(map[string]map[PaymentMethod]int),
		successTxs:   make(map[string]map[PaymentMethod]int),
	}
}

func (m *Metrics) UpdateSuccess(gateway string, method PaymentMethod, isSuccess bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.totalTxs[gateway] == nil {
		m.totalTxs[gateway] = make(map[PaymentMethod]int)
		m.successTxs[gateway] = make(map[PaymentMethod]int)
		m.successRates[gateway] = make(map[PaymentMethod]float64)
	}

	m.totalTxs[gateway][method]++
	if isSuccess {
		m.successTxs[gateway][method]++
	}

	total := m.totalTxs[gateway][method]
	success := m.successTxs[gateway][method]
	if total > 0 {
		m.successRates[gateway][method] = float64(success) / float64(total)
	} else {
		m.successRates[gateway][method] = 0.5 // Default
	}
}

func (m *Metrics) GetSuccessRate(gateway string, method PaymentMethod) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if rates, ok := m.successRates[gateway]; ok {
		if rate, ok := rates[method]; ok {
			return rate
		}
	}
	return 0.5 // Default
}

// router ----------------------------------------------

type Router interface {
	SelectGateway(method PaymentMethod) (Gateway, error)
}

type DynamicRouter struct {
	factory  *GatewayFactory
	metrics  *Metrics
	gateways []string // Available gateways
}

func NewDynamicRouter(factory *GatewayFactory, metrics *Metrics, gateways []string) *DynamicRouter {
	return &DynamicRouter{
		factory:  factory,
		metrics:  metrics,
		gateways: gateways,
	}
}

func (dr *DynamicRouter) SelectGateway(method PaymentMethod) (Gateway, error) {
	bestGateway := ""
	bestRate := -1.0

	for _, gwName := range dr.gateways {
		gw := dr.factory.CreateGateway(gwName)
		if gw == nil || !contains(gw.SupportedMethods(), method) {
			continue
		}
		rate := dr.metrics.GetSuccessRate(gw.Name(), method)
		if rate > bestRate {
			bestRate = rate
			bestGateway = gw.Name()
		}
	}

	if bestGateway == "" {
		return nil, errors.New("no suitable gateway found")
	}

	return dr.factory.CreateGateway(bestGateway), nil
}

// TransactionRepository (Repository Pattern for persistence)
type TransactionRepository interface {
	Save(tx *Transaction) error
	UpdateStatus(txID string, status TransactionStatus) error
	GetByID(txID string) (*Transaction, error)
}

type InMemoryTxRepo struct {
	txs map[string]*Transaction
	mu  sync.RWMutex
}

func NewInMemoryTxRepo() *InMemoryTxRepo {
	return &InMemoryTxRepo{txs: make(map[string]*Transaction)}
}

func (r *InMemoryTxRepo) Save(tx *Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()
	r.txs[tx.ID] = tx
	return nil
}

func (r *InMemoryTxRepo) UpdateStatus(txID string, status TransactionStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if tx, ok := r.txs[txID]; ok {
		tx.Status = status
		tx.UpdatedAt = time.Now()
		return nil
	}
	return errors.New("transaction not found")
}

func (r *InMemoryTxRepo) GetByID(txID string) (*Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if tx, ok := r.txs[txID]; ok {
		return tx, nil
	}
	return nil, errors.New("transaction not found")
}

// PaymentService (Facade/Orchestrator)
type PaymentService struct {
	repo    TransactionRepository
	router  Router
	metrics *Metrics
}

func NewPaymentService(repo TransactionRepository, router Router, metrics *Metrics) *PaymentService {
	return &PaymentService{
		repo:    repo,
		router:  router,
		metrics: metrics,
	}
}

// InitiateTransaction
func (ps *PaymentService) InitiateTransaction(userType UserType, amount float64, method PaymentMethod) (*Transaction, error) {
	tx := &Transaction{
		ID:            generateID(),
		UserType:      userType,
		Amount:        amount,
		PaymentMethod: method,
		Status:        Pending,
	}

	gw, err := ps.router.SelectGateway(method)
	if err != nil {
		return nil, err
	}

	resp, err := gw.Initiate(tx)
	if err != nil {
		return nil, err
	}

	tx.GatewayUsed = gw.Name()
	tx.GatewayTransactionID = resp.GatewayTxID
	tx.Status = Initiated

	if err := ps.repo.Save(tx); err != nil {
		return nil, err
	}

	// Async status polling or webhook handling can be added here
	go ps.pollStatus(tx)

	return tx, nil
}

// UpdateStatus (Webhook or Poll handler)
func (ps *PaymentService) UpdateStatus(gatewayTxID string, gatewayStatus string) error {
	// In real: Find tx by gatewayTxID (assume we have a reverse index)
	// For LLD, assume we pass txID or query repo
	// Mock: Assume we have txID from context

	var status TransactionStatus
	switch gatewayStatus {
	case "SUCCESS":
		status = Success
	case "FAILED":
		status = Failed
	default:
		status = Processing
	}

	// Assume txID is derived or passed; in real, use a map or DB query
	// For demo, skip full impl
	if err := ps.repo.UpdateStatus("dummy_tx_id", status); err != nil { // Replace with actual
		return err
	}

	// Update metrics (need gateway name and method; assume derived)
	// ps.metrics.UpdateSuccess("RazorPay", Card, status == Success)

	return nil
}

func (ps *PaymentService) pollStatus(tx *Transaction) {
	// Mock polling every 5s up to 1min
	gw := ps.router.(*DynamicRouter).factory.CreateGateway(tx.GatewayUsed)
	for i := 0; i < 12; i++ { // 1 min
		time.Sleep(5 * time.Second)
		if tx.Status != Pending && tx.Status != Initiated {
			return
		}
		resp, err := gw.GetStatus(tx.GatewayTransactionID)
		if err != nil {
			continue
		}
		var newStatus TransactionStatus
		switch resp.Status {
		case "SUCCESS":
			newStatus = Success
		case "FAILED":
			newStatus = Failed
		default:
			newStatus = Processing
		}
		ps.repo.UpdateStatus(tx.ID, newStatus)
		ps.metrics.UpdateSuccess(tx.GatewayUsed, tx.PaymentMethod, newStatus == Success)
		if newStatus == Success || newStatus == Failed {
			return
		}
	}
}

// Example Usage (main for demo)
func main() {
	factory := NewGatewayFactory()
	metrics := NewMetrics()
	repo := NewInMemoryTxRepo()
	router := NewDynamicRouter(factory, metrics, []string{"PayU", "Paytm", "RazorPay"})
	service := NewPaymentService(repo, router, metrics)

	// Simulate init
	tx, err := service.InitiateTransaction(Customer, 100.0, Card)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Initiated: %+v\n", tx)

	// Simulate update (e.g., via webhook)
	service.UpdateStatus(tx.GatewayTransactionID, "SUCCESS")

	// Get updated
	updatedTx, _ := repo.GetByID(tx.ID)
	fmt.Printf("Updated: %+v\n", updatedTx)
}
