package product

import (
	"github.com/google/uuid"
)

type Product struct {
	ProductID uuid.UUID
	Name      string
	Price     float64
	Quantity  int
}

func NewProduct(name string, price float64, quantity int) *Product {
	return &Product{
		ProductID: uuid.New(),
		Name:      name,
		Price:     price,
		Quantity:  quantity,
	}
}

func (p *Product) UpdatePrice(newPrice float64) {
	p.Price = newPrice
}

func (p *Product) UpdateQuantity(newQuantity int) {
	p.Quantity = newQuantity
}
