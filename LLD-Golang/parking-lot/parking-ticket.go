package main

import (
	"time"
)

const baseCharge = 100.00

type ParkingTicket struct {
	EntryTime   time.Time
	ExitTime    time.Time
	Vehicle     VehicleInterface
	Spot        *ParkingSpot
	TotalCharge float64
}

func NewParkingTicket(vehicle VehicleInterface, spot *ParkingSpot) *ParkingTicket {
	return &ParkingTicket{EntryTime: time.Now(), ExitTime: time.Time{}, Vehicle: vehicle, Spot: spot, TotalCharge: 0.00}
}

func (p *ParkingTicket) SetExitTime(exitTime time.Time) {
	p.ExitTime = exitTime
}

func (p *ParkingTicket) CalculateTotalCharge() float64 {
	if p.ExitTime == (time.Time{}) {
		p.TotalCharge = baseCharge
		return p.TotalCharge
	}
	duration := p.ExitTime.Sub(p.EntryTime)
	hours := duration.Hours()
	additionalCharge := hours * p.Vehicle.GetVehicleCost()
	p.TotalCharge = baseCharge + additionalCharge
	return p.TotalCharge
}
