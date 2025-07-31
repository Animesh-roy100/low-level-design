package main

type ParkingLot struct {
	Name   string
	Floors []*ParkingFloor
}

func NewParkingLot(name string) *ParkingLot {
	return &ParkingLot{
		Name:   name,
		Floors: []*ParkingFloor{},
	}
}

func (pl *ParkingLot) AddFloor(floor *ParkingFloor) {
	pl.Floors = append(pl.Floors, floor)
}
