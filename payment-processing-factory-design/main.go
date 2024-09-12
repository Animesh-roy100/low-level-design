package main

import "fmt"

// Factory Design Pattern

// The PaymentProcessor interface defines a common method for processing payments.
type PaymentProcessor interface {
	ProcessPayment(amount float64)
}

// Concrete classes (CreditCardProcessor, NetBankingProcessor, UPIProcessor, DebitCardProcessor)
// implement specific payment processing logic.

type CreditCardProcessor struct{}

func (c *CreditCardProcessor) ProcessPayment(amount float64) {
	fmt.Println("Payment processed using Credit Card: ", amount)
}

type NetBankingProcessor struct{}

func (n *NetBankingProcessor) ProcessPayment(amount float64) {
	fmt.Println("Payment processed using Net Banking: ", amount)
}

type UPIProcessor struct{}

func (u *UPIProcessor) ProcessPayment(amount float64) {
	fmt.Println("Payment processed using UPI: ", amount)
}

type DebitCardProcessor struct{}

func (d *DebitCardProcessor) ProcessPayment(amount float64) {
	fmt.Println("Payment processed using Debit Card: ", amount)
}

// PaymentProcessorFactory is responsible for creating instances
// of the appropriate payment processor based on the type provided.
func PaymentProcessFactory(paymentType string) PaymentProcessor {
	switch paymentType {
	case "CreditCard":
		return &CreditCardProcessor{}
	case "NetBanking":
		return &NetBankingProcessor{}
	case "UPI":
		return &UPIProcessor{}
	case "DebitCard":
		return &DebitCardProcessor{}
	default:
		return nil
	}
}

func main() {
	fmt.Println("Payment Processing Factory Design Pattern")

	processor := PaymentProcessFactory("CreditCard")
	processor.ProcessPayment(1000.0)
	processor = PaymentProcessFactory("NetBanking")
	processor.ProcessPayment(1000.0)
	processor = PaymentProcessFactory("UPI")
	processor.ProcessPayment(1000.0)
	processor = PaymentProcessFactory("DebitCard")
	processor.ProcessPayment(1000.0)
}
