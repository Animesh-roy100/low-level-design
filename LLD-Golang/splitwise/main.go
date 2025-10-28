package main

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
)

// ======================= Utilities =======================
//
// WHY: Money math should be rounded deterministically to avoid drifting
// balances due to floating-point precision. We centralize rounding in one
// helper and use it EVERYWHERE we touch a monetary value.
func round2(x float64) float64 { return math.Round(x*100) / 100 }

//
// ======================= Domain Types =======================
//

// SplitType expresses how an expense is split among users.
// WHY: Strategy pattern needs an enum-like discriminator to pick a strategy.
type SplitType int

const (
	EQUAL SplitType = iota
	EXACT
	PERCENTAGE
)

// Split is an atomic "who owes how much" record.
// WHY: Value object produced by strategies; simple and serializable.
type Split struct {
	User      string
	Amount    float64
	SplitType SplitType
}

// User models a participant and keeps running totals.
// WHY: We maintain both "pays" and "borrows" so we can compute Net on demand.
// Net = TotalPays - TotalBorrows (positive means others owe this user).
type User struct {
	UserID       string
	Name         string
	Email        string
	Mobile       string
	TotalBorrows float64
	TotalPays    float64
	Net          float64
}

var (
	// WHY: Make user-ID generation threadsafe in case the app becomes concurrent.
	counterMu sync.Mutex
	counter   int
)

// NewUser auto-assigns an incrementing numeric ID as a string.
// WHY: Mirrors your C# style, keeps demo simple, avoids external ID deps.
func NewUser(name, email, mobile string) *User {
	counterMu.Lock()
	counter++
	id := fmt.Sprintf("%d", counter)
	counterMu.Unlock()

	return &User{
		UserID:       id,
		Name:         name,
		Email:        email,
		Mobile:       mobile,
		TotalBorrows: 0,
		TotalPays:    0,
		Net:          0,
	}
}

// Expense represents one bill pay event with a chosen splitting method.
// WHY: Keeps both raw inputs (PaidBy, Amount, SplitType) and computed
// artifacts (Splits, NetBalance) for observability/debugging.
type Expense struct {
	PaidBy      string
	Amount      float64
	Splits      []Split
	SplitType   SplitType
	NetBalance  float64
	ExpenseName string
}

// NewExpense constructs and computes the net balance visible to the payer.
// WHY: calculate once at creation; immutable thereafter (by convention).
func NewExpense(expenseName, paidBy string, amount float64, splits []Split, splitType SplitType) *Expense {
	e := &Expense{
		ExpenseName: expenseName,
		PaidBy:      paidBy,
		Amount:      round2(amount), // WHY: normalize incoming amount immediately
		Splits:      splits,
		SplitType:   splitType,
	}
	e.NetBalance = e.calculateNetBalance()
	return e
}

// calculateNetBalance sums how much "others" owe in this expense.
// WHY: Payer often wants to see "how much should come back to me" at a glance.
func (e *Expense) calculateNetBalance() float64 {
	total := 0.0
	for _, split := range e.Splits {
		if split.User == e.PaidBy {
			continue
		}
		total += split.Amount
	}
	e.NetBalance = round2(total)
	return e.NetBalance
}

//
// ======================= UserManager =======================
//
// WHY: Aggregate root for users and their balances. Centralizes all updates
// to enforce invariants (e.g., always round, always recalc Net).
// Thread-safe so we can later expose HTTP handlers without rewrites.
//

type UserManager struct {
	mu    sync.Mutex
	users map[string]*User
}

var (
	userManagerInstance *UserManager
	userManagerOnce     sync.Once
)

// GetUserManager is a threadsafe, lazy singleton constructor.
// WHY: Keeps global state controlled; easy to swap for DI in the future.
func GetUserManager() *UserManager {
	userManagerOnce.Do(func() {
		userManagerInstance = &UserManager{
			users: make(map[string]*User),
		}
	})
	return userManagerInstance
}

func (um *UserManager) GetUser(userID string) *User {
	um.mu.Lock()
	defer um.mu.Unlock()
	return um.users[userID]
}

func (um *UserManager) AddBorrow(userID string, amount float64) bool {
	um.mu.Lock()
	defer um.mu.Unlock()
	if user, ok := um.users[userID]; ok {
		// WHY: Round each step to keep ledger stable across ops.
		user.TotalBorrows = round2(user.TotalBorrows + amount)
		user.Net = round2(user.TotalPays - user.TotalBorrows)
		return true
	}
	return false
}

func (um *UserManager) AddPay(userID string, amount float64) bool {
	um.mu.Lock()
	defer um.mu.Unlock()
	if user, ok := um.users[userID]; ok {
		user.TotalPays = round2(user.TotalPays + amount)
		user.Net = round2(user.TotalPays - user.TotalBorrows)
		return true
	}
	return false
}

func (um *UserManager) GetNet(userID string) float64 {
	um.mu.Lock()
	defer um.mu.Unlock()
	if user, ok := um.users[userID]; ok {
		return user.Net
	}
	return 0
}

func (um *UserManager) AddUser(user *User) *User {
	um.mu.Lock()
	defer um.mu.Unlock()
	if _, exists := um.users[user.UserID]; !exists {
		um.users[user.UserID] = user
	}
	return user
}

// UpdateUserBorrows applies the split to all participants and credits the payer.
// WHY: Single point of truth for applying financial effects of an expense.
// - Payer gets TotalPays += full amount (what they fronted)
// - Everyone (including payer) gets TotalBorrows += their share
// - Net is recomputed and rounded for each user.
func (um *UserManager) UpdateUserBorrows(splits []Split, payer string, amount float64) bool {
	um.mu.Lock()
	defer um.mu.Unlock()

	for _, sp := range splits {
		u, ok := um.users[sp.User]
		if !ok {
			// WHY: In production, you'd return error; here we skip gracefully.
			continue
		}
		if sp.User == payer {
			u.TotalPays = round2(u.TotalPays + amount)
			u.TotalBorrows = round2(u.TotalBorrows + sp.Amount)
		} else {
			u.TotalBorrows = round2(u.TotalBorrows + sp.Amount)
		}
		u.Net = round2(u.TotalPays - u.TotalBorrows)
	}
	return true
}

//
// ======================= Split Strategies =======================
//
// WHY: Strategy pattern cleanly separates the split math from orchestration.
// We also verify sums and update user balances atomically (verify-then-apply).
//

// SplitStrategy is the contract for all splitting algorithms.
type SplitStrategy interface {
	Split(userID string, users []string, amount float64, subAmounts []float64) []Split
}

// VerifySplit ensures numerical integrity and only then updates user balances.
// WHY: Prevents recording invalid expenses and keeps the ledger consistent.
func VerifySplit(split []Split, amount float64, userID string) bool {
	sum := 0.0
	for _, s := range split {
		sum += s.Amount
	}
	// WHY: Allow tiny floating error; 1 cent tolerance is typical for FP ops.
	if math.Abs(sum-amount) > 0.01 {
		fmt.Println("Verify split failed.")
		return false
	}
	GetUserManager().UpdateUserBorrows(split, userID, amount)
	return true
}

// EqualSplitStrategy divides equally; last share absorbs rounding residue.
// WHY: Avoids a "lost cent" by tail-adjusting final participant.
type EqualSplitStrategy struct{}

func (EqualSplitStrategy) Split(userID string, users []string, amount float64, _ []float64) []Split {
	n := len(users)
	each := round2(amount / float64(n))
	// tail: ensures the total of shares == amount after rounding
	tail := round2(amount - each*float64(n-1))

	splits := make([]Split, 0, n)
	for i, u := range users {
		share := each
		if i == n-1 {
			share = tail
		}
		splits = append(splits, Split{User: u, Amount: share, SplitType: EQUAL})
	}
	if !VerifySplit(splits, amount, userID) {
		return nil
	}
	return splits
}

// ExactSplitStrategy uses exact sub-amounts; we just round each provided value.
// WHY: Caller controls distribution; we still enforce total integrity.
type ExactSplitStrategy struct{}

func (ExactSplitStrategy) Split(userID string, users []string, amount float64, sub []float64) []Split {
	if len(sub) != len(users) {
		fmt.Println("exact split: subAmounts and users length mismatch")
		return nil
	}
	splits := make([]Split, 0, len(users))
	for i, u := range users {
		splits = append(splits, Split{User: u, Amount: round2(sub[i]), SplitType: EXACT})
	}
	if !VerifySplit(splits, amount, userID) {
		return nil
	}
	return splits
}

// PercentageSplitStrategy multiplies by percentages; last share tail-adjusts.
// WHY: Rounding + percentages can drift; tail ensures sum matches amount exactly.
type PercentageSplitStrategy struct{}

func (PercentageSplitStrategy) Split(userID string, users []string, amount float64, sub []float64) []Split {
	if len(sub) != len(users) {
		fmt.Println("percentage split: sub percents and users length mismatch")
		return nil
	}
	n := len(users)
	cur := 0.0
	splits := make([]Split, 0, n)
	for i, u := range users {
		share := round2(amount * sub[i])
		if i == n-1 {
			share = round2(amount - cur) // tail adjust to eliminate rounding residue
		}
		splits = append(splits, Split{User: u, Amount: share, SplitType: PERCENTAGE})
		cur += share
	}
	if !VerifySplit(splits, amount, userID) {
		return nil
	}
	return splits
}

// CreateSplitStrategy maps SplitType to a strategy instance.
// WHY: Factory isolates construction logic from callers.
func CreateSplitStrategy(splitType SplitType) SplitStrategy {
	switch splitType {
	case EQUAL:
		return EqualSplitStrategy{}
	case EXACT:
		return ExactSplitStrategy{}
	case PERCENTAGE:
		return PercentageSplitStrategy{}
	default:
		return EqualSplitStrategy{}
	}
}

// ======================= ExpenseManager =======================
//
// WHY: Aggregate root for expenses. We keep a per-payer index (map[payer][]Expense)
// to quickly list what someone paid. We also expose ShowExpenses(userID) to list
// any expense a user is involved in (payer or participant).
// RWMutex: writes are rare, reads are frequent.
type ExpenseManager struct {
	mu       sync.RWMutex
	expenses map[string][]*Expense
}

var (
	expenseManagerInstance *ExpenseManager
	expenseManagerOnce     sync.Once
)

func GetExpenseManager() *ExpenseManager {
	expenseManagerOnce.Do(func() {
		expenseManagerInstance = &ExpenseManager{
			expenses: make(map[string][]*Expense),
		}
	})
	return expenseManagerInstance
}

// AddExpense orchestrates: pick strategy -> split -> verify -> record expense.
// WHY: If verification fails, we DO NOT record the expense (nil returned).
func (em *ExpenseManager) AddExpense(expenseName, userID string, splitType SplitType, users []string, amount float64, subAmounts []float64) *Expense {
	strategy := CreateSplitStrategy(splitType)
	splits := strategy.Split(userID, users, round2(amount), subAmounts)
	if splits == nil {
		// invalid split; don't record expense
		return nil
	}

	exp := NewExpense(expenseName, userID, amount, splits, splitType)

	em.mu.Lock()
	em.expenses[userID] = append(em.expenses[userID], exp)
	em.mu.Unlock()
	return exp
}

// ShowExpenses returns every expense where userID appears in any split.
// WHY: Useful to render a user's dashboard of "things I’m part of".
func (em *ExpenseManager) ShowExpenses(userID string) []*Expense {
	em.mu.RLock()
	defer em.mu.RUnlock()

	var result []*Expense
	for _, userExpenses := range em.expenses {
		for _, expense := range userExpenses {
			for _, split := range expense.Splits {
				if split.User == userID {
					result = append(result, expense)
					break
				}
			}
		}
	}
	return result
}

// ======================= Demo / Sanity Test =======================
//
// WHY: Quick end-to-end run to validate flow and see JSON output.
// This section doubles as an example of how to use the API.
func main() {
	expenseManager := GetExpenseManager()
	userManager := GetUserManager()

	// Create users (IDs auto-assigned: "1", "2", "3")
	u1 := NewUser("animesh", "vamsi@gmail.com", "9999999999")
	u2 := NewUser("roy", "krishna@gmail.com", "8888888888")
	u3 := NewUser("somu", "jani@gmail.com", "7777777777")
	userManager.AddUser(u1)
	userManager.AddUser(u2)
	userManager.AddUser(u3)

	users := []string{u1.UserID, u2.UserID, u3.UserID}

	// Expense 1: Equal split of 100 paid by u1
	_ = expenseManager.AddExpense("Red Biryani", u1.UserID, EQUAL, users, 100, nil)

	// Expense 2: Percentage split 30/40/30 of 200 paid by u2
	_ = expenseManager.AddExpense("Groceries", u2.UserID, PERCENTAGE, users, 200, []float64{0.3, 0.4, 0.3})

	// Show all expenses where user "1" is involved
	result := expenseManager.ShowExpenses("1")
	b, _ := json.Marshal(result)
	fmt.Println(string(b))

	// Inspect user balances
	b1, _ := json.Marshal(u1)
	fmt.Println(string(b1))
	b2, _ := json.Marshal(u2)
	fmt.Println(string(b2))
	b3, _ := json.Marshal(u3)
	fmt.Println(string(b3))
}
