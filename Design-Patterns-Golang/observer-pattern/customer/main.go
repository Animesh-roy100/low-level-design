package main

import "fmt"

type Subject interface {
	Subscribe(observer Observer)
	Unsubscribe(observer Observer)
	NotifyAll()
}

type Item struct {
	ObserverList []Observer
	Name         string
	InStock      bool
}

func NewItem(name string) *Item {
	return &Item{
		Name: name,
	}
}

func (i *Item) Subscribe(observer Observer) {
	i.ObserverList = append(i.ObserverList, observer)
}

func (i *Item) Unsubscribe(observer Observer) {
	index := 0
	for i, o := range i.ObserverList {
		if o == observer {
			index = i
			break
		}
	}

	i.ObserverList = append(i.ObserverList[:index], i.ObserverList[index+1:]...)
}

func (i *Item) NotifyAll() {
	for _, o := range i.ObserverList {
		o.Update(i.Name)
	}
}

func (i *Item) UpdateAvailability() {
	fmt.Printf("Item %s is now in stock\n", i.Name)
	i.InStock = true
	i.NotifyAll()
}

type Observer interface {
	Update(name string)
	GetID() string
}

type Customer struct {
	id string
}

func (c *Customer) Update(itemName string) {
	fmt.Printf("Sending email to customer %s for item %s\n", c.id, itemName)
}

func (c *Customer) GetID() string {
	return c.id
}

func main() {
	item := NewItem("Shirt")

	customer1 := &Customer{id: "1"}
	customer2 := &Customer{id: "2"}

	item.Subscribe(customer1)
	item.Subscribe(customer2)

	item.UpdateAvailability()
}
