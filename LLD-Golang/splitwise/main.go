package splitwise

type User struct {
	ID     string
	Name   string
	Email  string
	Mobile string
}

type Expense struct {
	PayerID      string
	Amount       float64
	Participants []string
	SplitType    string // "EQUAL", "EXACT"
	Shares       []float64
}

type SplitStrategy interface {
	ComputeShares(amount float64, participants []string, shares []float64) ([]float64, error)
}

type EqualSplitStrategy struct{}

// func (e *EqualSplitStrategy) ComputeShares(amount float64, participants []string, shares []float64) ([]float64, error) {

// }

type ExactSplitStrategy struct{}

// func (e *ExactSplitStrategy) ComputeShares(amount float64, participants []string, shares []float64) ([]float64, error) {

// }

// type PercentageSplitStrategy struct{}

// func (p *PercentageSplitStrategy) ComputeShares(amount float64, participants []string, shares []float64) ([]float64, error) {
// }
