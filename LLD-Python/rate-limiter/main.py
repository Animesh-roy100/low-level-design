from abc import ABC, abstractmethod
import collections
import threading
import time

class RateLimiter(ABC):
    @abstractmethod
    def allow_requests(self, user_id: str) -> bool:
        pass


class SlidingWindowRateLimiter(RateLimiter):
    def __init__(self, max_requests: int, time_window: float):
        self.max_requests = max_requests
        self.time_window = time_window
        self.user_requests = {} # user_id -> deque of timestamps
        self.lock = threading.Lock()
    
    def allow_requests(self, user_id: str) -> bool:
       current_time = time.time()
       with self.lock:
            if user_id not in self.user_requests:
               self.user_requests[user_id] = collections.deque()
            
            requests = self.user_requests[user_id]
            while requests and requests[0] < current_time - self.time_window:
                requests.popleft()
            
            if len(requests) < self.max_requests:
                requests.append(current_time)
                # print("Sliding Window: True")
                return True
            # print("Sliding Window: False")
            return False
        

class TokenBucketRateLimiter(RateLimiter):
    def __init__(self, capacity: int, refill_rate: float):
        self.capacity = capacity
        self.refill_rate = refill_rate
        self.buckets = {} # user_id -> {'tokens': float, 'last_refill': float}
        self.lock = threading.Lock()
    
    def allow_requests(self, user_id: str):
        now = time.time()
        with self.lock:
            if user_id not in self.buckets:
                self.buckets[user_id] = {'tokens': self.capacity, 'last_refill': now}
            
            bucket = self.buckets[user_id]

            # Refill tokens
            elapsed = now - bucket['last_refill']
            refill_tokens = elapsed * self.refill_rate

            bucket['tokens'] = min(self.capacity, bucket['tokens'] + refill_tokens)
            bucket['last_refill'] = now

            if bucket['tokens'] >= 1:
                bucket['tokens'] -= 1
                # print("Token Bucket: True")
                return True
            # print("Token Bucket: False")
            return False

# class RateLimiterFactory:
#     @staticmethod
#     def create(rate_limiter_type: str) -> RateLimiter:
#         if rate_limiter_type == 'token_bucket':
#             return TokenBucketRateLimiter(100, 100/60)
#         elif rate_limiter_type == 'sliding_window':
#             return SlidingWindowRateLimiter(100, 60.0)
#         else:
#             raise ValueError(f"Unknown Limiter Type: {rate_limiter_type}")

def RateLimiterFactory(rate_limiter_type: str) -> RateLimiter:
    if rate_limiter_type == 'token_bucket':
        return TokenBucketRateLimiter(100, 100/60)
    elif rate_limiter_type == 'sliding_window':
        return SlidingWindowRateLimiter(100, 60.0)
    else:
        raise ValueError(f"Unknown Limiter Type: {rate_limiter_type}")
        
if __name__ == "__main__":
    # sliding_window = RateLimiterFactory.create('sliding_window')
    # token_bucket = RateLimiterFactory.create('token_bucket')
    sliding_window = RateLimiterFactory('sliding_window')
    token_bucket = RateLimiterFactory('token_bucket')
    # test = RateLimiterFactory('test')

    user = "1"

    print("Sliding Window Test: ")
    for i in range(105):
        if sliding_window.allow_requests(user):
            pass
        else:
            print(f"Request {i+1} denied")
            break
    
    print("Token Bucket Test: ")
    allowed_count = 0
    for i in range(105):
        if token_bucket.allow_requests(user):
            allowed_count += 1
    print(f"Allowed: {allowed_count}/105")

    time.sleep(1) # Refill
    if token_bucket.allow_requests(user):
        print("post wait request allowed")
