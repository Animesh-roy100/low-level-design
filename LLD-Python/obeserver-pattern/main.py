from abc import ABC, abstractmethod
from typing import List


# ================= Observer Interface =================
class Observer(ABC):
    @abstractmethod
    def update(self, news: str) -> None:
        pass


# ================= Subject Interface =================
class Subject(ABC):
    @abstractmethod
    def subscribe(self, observer: Observer) -> None:
        pass
    
    @abstractmethod
    def unsubscribe(self, observer: Observer) -> None:
        pass

    @abstractmethod
    def notify(self) -> None:
        pass


# ================= Concrete Subject =================
class NewsSubject(Subject):
    def __init__(self):
        self.observers: List[Observer] = []   
        self.news: str = ""

    def subscribe(self, observer: Observer) -> None:
        self.observers.append(observer)
    
    def unsubscribe(self, observer: Observer) -> None:
        self.observers = [
            o for o in self.observers if o != observer
        ]
    
    def notify(self) -> None:
        for observer in self.observers:
            observer.update(self.news)
    
    def set_news(self, news: str) -> None:
        self.news = news
        self.notify()

# ================= Concrete Observer =================
class NewsObserver(Observer):
    def __init__(self, name: str):
        self.name = name
    
    def update(self, news: str) -> None:
        print(f"{self.name} received news: {news}")


# ================= Main =================
def main():
    news_subject = NewsSubject()

    observer1 = NewsObserver("Animesh")
    observer2 = NewsObserver("Bikram")

    news_subject.subscribe(observer1)
    news_subject.subscribe(observer2)

    news_subject.set_news("Breaking News: Important Announcement!")
    news_subject.set_news("Another News Update!")

    news_subject.unsubscribe(observer1)

    news_subject.set_news("New News update!")

if __name__ == "__main__":
    main()
    
