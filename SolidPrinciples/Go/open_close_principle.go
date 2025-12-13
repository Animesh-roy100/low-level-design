package main

/*
A class should be open for extension but closed for modification.
This means you should be able to add new functionality without changing existing code.

Open/Closed Principle means we should design our code so that new behavior can be added through extension
(interfaces, composition) rather than by modifying existing logic.

*/

/*

Violates Open/Closed Principle

type Shape struct {
	Type   string
	Width  float64
	Height float64
	Radius float64
}

func Area(s Shape) float64 {
	switch s.Type {
	case "rectangle":
		return s.Width * s.Height
	case "circle":
		return 3.14 * s.Radius * s.Radius
	default:
		return 0
	}
}

func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += Area(s)
	}
	return total
}

*/

// Shape is closed for modification
type Shape interface {
	Area() float64
}

// Rectangle
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Circle
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

// Triangle (Extension without modification)
type Triangle struct {
	Base, Height float64
}

func (t Triangle) Area() float64 {
	return 0.5 * t.Base * t.Height
}

// Closed for modification
func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}
