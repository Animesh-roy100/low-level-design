package main

import "fmt"

// observer interface ------------------------------------------------
type Observer interface {
	Update(temp float64)
}

// subject interface -------------------------------------------------
type Subject interface {
	RegisterObserver(observer Observer)
	RemoveObserver(observer Observer)
	NotifyObservers()
}

// concrete implementation of the Subject interface ------------------
type WeatherStation struct {
	observers   []Observer
	temperature float64
}

func (w *WeatherStation) RegisterObserver(observer Observer) {
	w.observers = append(w.observers, observer)
}

func (w *WeatherStation) RemoveObserver(observer Observer) {
	var index int
	for i, o := range w.observers {
		if o == observer {
			index = i
			break
		}
	}

	w.observers = append(w.observers[:index], w.observers[index+1:]...)
}

func (w *WeatherStation) NotifyObservers() {
	for _, observer := range w.observers {
		observer.Update(w.temperature)
	}
}

func (w *WeatherStation) SetTemperature(temp float64) {
	w.temperature = temp
	w.NotifyObservers()
}

// concrete implementations of the Observer interface ----------------
type PhoneDisplay struct {
	station *WeatherStation
}

func (p *PhoneDisplay) Update(temperature float64) {
	fmt.Printf("Phone display: Current temperature is %.2f°C\n", temperature)
}

type WindowDisplay struct {
	station *WeatherStation
}

func (w *WindowDisplay) Update(temperature float64) {
	fmt.Printf("Window display: Current temperature is %.2f°C\n", temperature)
}

func main() {
	// create weather station
	weatherStation := &WeatherStation{}

	// create observers
	phoneDisplay := &PhoneDisplay{station: weatherStation}
	windowDisplay := &WindowDisplay{station: weatherStation}

	// register observers with the weather station
	weatherStation.RegisterObserver(phoneDisplay)
	weatherStation.RegisterObserver(windowDisplay)

	weatherStation.SetTemperature(25.5)
	weatherStation.SetTemperature(30.0)

	// remove an observer
	weatherStation.RemoveObserver(phoneDisplay)

	weatherStation.SetTemperature(28.0)
}
