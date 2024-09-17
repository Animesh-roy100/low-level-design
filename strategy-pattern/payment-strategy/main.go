package main

import "fmt"

type PaymentMethod interface {
	Pay(amount float64) string
}

type CreditCard struct {
	Name       string
	CardNumber string
}

func (c *CreditCard) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f using Credit Card (%s)", amount, c.CardNumber)
}

type Cash struct{}

func (c *Cash) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f in cash", amount)
}

type UPI struct {
	UPIID string
}

func (u *UPI) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f using UPI ID (%s)", amount, u.UPIID)
}

type DebitCard struct {
	Name       string
	CardNumber string
}

func (d *DebitCard) Pay(amount float64) string {
	return fmt.Sprintf("Paid %.2f using Debit Card (%s)", amount, d.CardNumber)
}

type Item struct{}

type ShoppingCart struct {
	items         []Item
	paymentMethod PaymentMethod
}

func (s *ShoppingCart) SetPaymentMethod(paymentMethod PaymentMethod) {
	s.paymentMethod = paymentMethod
}

func (s *ShoppingCart) Checkout(amount float64) {
	if s.paymentMethod == nil {
		fmt.Println("No payment method set")
		return
	}
	fmt.Println(s.paymentMethod.Pay(amount))
}

func main() {
	// Create a shopping cart
	cart := &ShoppingCart{}

	creditCard := &CreditCard{Name: "John Doe", CardNumber: "1234-5678-9012-3456"}
	cart.SetPaymentMethod(creditCard)
	cart.Checkout(200.50)

	cash := &Cash{}
	cart.SetPaymentMethod(cash)
	cart.Checkout(100.00)

	upi := &UPI{UPIID: "john@upi"}
	cart.SetPaymentMethod(upi)
	cart.Checkout(50.25)

	debitCard := &DebitCard{Name: "John Doe", CardNumber: "6543-2109-8765-4321"}
	cart.SetPaymentMethod(debitCard)
	cart.Checkout(150.75)
}
