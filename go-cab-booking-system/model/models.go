package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// UserID         uint
	Name           string
	Phone          string `gorm:"uniqueIndex"`
	Email          string `gorm:"uniqueIndex"`
	Password       string
	SavedLocations []Location `gorm:"many2many:user_saved_locations;"`
	RideHistory    []Ride     `gorm:"foreignKey:UserID;"`
}

type Location struct {
	gorm.Model
	// UserID    uint
	Latitude  float64
	Longitude float64
	Address   string
}

type Driver struct {
	gorm.Model
	// DriverID        uint
	Name            string
	Phone           string  `gorm:"uniqueIndex"`
	Email           string  `gorm:"uniqueIndex"`
	VehicleDetails  Vehicle `gorm:"foreignKey:DriverID;"`
	DriverStatus    string
	CurrentLocation string `gorm:"foreignKey:CurrentLocationID;"`
}

type Vehicle struct {
	gorm.Model
	// VehicleID     uint
	DriverID      uint
	VehicleNumber string `gorm:"uniqueIndex"`
	VehicleType   string
	Capacity      int
}

type Ride struct {
	gorm.Model
	// RideID         uint
	Driver         uint
	User           uint
	PickupLocation string `gorm:"foreignKey:PickupLocationID"`
	DropLocation   string `gorm:"foreignKey:DropLocationID"`
	RideStatus     string
	Fare           float64
	StartTime      time.Time
	EndTime        time.Time
	Payment        Payment `gorm:"foreignKey:RideID"`
}

type Payment struct {
	gorm.Model
	// PaymentID     uint
	RideID        uint
	UserID        uint
	DriverID      uint
	PaymentAmount float64
	PaymentMethod string
	PaymentStatus string
}
