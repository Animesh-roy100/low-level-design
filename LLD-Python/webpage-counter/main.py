import threading

class WebPageVisitCounter:
    def __init__(self):
        self.total_pages = 0
        self.visit_counts = []
        self.lock = threading.Lock()

    def init(self, totalPages: int):
        # Initialize the visit counter.
        with self.lock:
            self.total_pages = totalPages
            self.visit_counts = [0] * totalPages
            print(f"Initialized visit counter with {totalPages} pages")

    def incrementVisitCount(self, pageIndex: int):
        # Increment visit count for a webpage.
        with self.lock:
            if pageIndex < 0 or pageIndex >= self.total_pages:
                print(f"Invalid page index: {pageIndex}")
                return

            self.visit_counts[pageIndex] += 1
            print(f"Page {pageIndex} visit incremented to {self.visit_counts[pageIndex]}")

    def getVisitCount(self, pageIndex: int) -> int:
        # Get total visit count for a webpage.
        with self.lock:
            if pageIndex < 0 or pageIndex >= self.total_pages:
                print(f"Invalid page index: {pageIndex}")
                return 0

            count = self.visit_counts[pageIndex]
            print(f"Page {pageIndex} visit count requested: {count}")
            return count


class Helper: 
    def log(self, message: str): 
        print(message) 

helper = Helper() 
counter = WebPageVisitCounter() 
counter.init(2) 
counter.incrementVisitCount(0) 
counter.incrementVisitCount(1) 
counter.incrementVisitCount(1) 
counter.incrementVisitCount(1) 
counter.incrementVisitCount(0) 
print(counter.getVisitCount(0)) # Output: 2 
print(counter.getVisitCount(1)) # Output: 3