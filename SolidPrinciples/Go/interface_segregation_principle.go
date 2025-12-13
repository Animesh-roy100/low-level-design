package main

/*

No client should be forced to depend on methods it does not use.
It is better to have many small, client-specific interfaces than one large, general-purpose interface.

Interface Segregation Principle states that no client should be forced to depend on methods it does not use,
so interfaces should be small and focused.

Before - Violates ISP

type Worker interface {
	Work()
	Eat()
	Sleep()
}

type Human struct{}

func (h Human) Work()  { fmt.Println("Human working") }
func (h Human) Eat()   { fmt.Println("Human eating") }
func (h Human) Sleep() { fmt.Println("Human sleeping") }

type Robot struct{}

func (r Robot) Work()  { fmt.Println("Robot working") }
func (r Robot) Eat()   ----> meaningless for Robot
func (r Robot) Sleep()  ----> meaningless for Robot

*/

type Workable interface {
	Work()
}

type Eatable interface {
	Eat()
}

type Sleepable interface {
	Sleep()
}

type Human struct{}

func (h Human) Work()  { println("Human working") }
func (h Human) Eat()   { println("Human eating") }
func (h Human) Sleep() { println("Human sleeping") }

type Robot struct{}

func (r Robot) Work() { println("Robot working") }

// func main() {
// 	var w1 Workable = Human{}
// 	var e1 Eatable = Human{}
// 	var s1 Sleepable = Human{}

// 	w1.Work()
// 	e1.Eat()
// 	s1.Sleep()

// 	var w2 Workable = Robot{}
// 	w2.Work()
// }
