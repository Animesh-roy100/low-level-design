package main

import "fmt"

type Observer interface {
	Update(news string)
}

type Subject interface {
	Subscribe(observer Observer)
	Unsubscribe(observer Observer)
	Notify()
}

type NewsSubject struct {
	Observers []Observer
	News      string
}

func (n *NewsSubject) Subscribe(observer Observer) {
	n.Observers = append(n.Observers, observer)
}

func (n *NewsSubject) Unsubscribe(observer Observer) {
	index := 0
	for i, n := range n.Observers {
		if n == observer {
			index = i
			break
		}
	}

	n.Observers = append(n.Observers[:index], n.Observers[index+1:]...)
}

func (n *NewsSubject) Notify() {
	for _, observer := range n.Observers {
		observer.Update(n.News)
	}
}

func (n *NewsSubject) SetNews(news string) {
	n.News = news
	n.Notify()
}

type NewsObserver struct {
	Name string
}

func (n *NewsObserver) Update(news string) {
	fmt.Printf("%s received news: %s\n", n.Name, news)
}

func main() {
	newsSubject := &NewsSubject{}

	observer1 := &NewsObserver{Name: "Animesh"}
	observer2 := &NewsObserver{Name: "Bikram"}

	newsSubject.Subscribe(observer1)
	newsSubject.Subscribe(observer2)

	newsSubject.SetNews("Breaking News: Important Announcement!")
	newsSubject.SetNews("Another News Update!")

	newsSubject.Unsubscribe(observer1)

	newsSubject.SetNews("New News update!")
}
