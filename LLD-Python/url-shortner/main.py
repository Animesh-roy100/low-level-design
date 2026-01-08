from abc import ABC, abstractmethod
import time

# User Data 
class UserData:
    def __init__(self, user_id: str, short_url: str, long_url: str, current_time: float, expiration_minutes: int):
        self.user_id = user_id
        self.short_url = short_url
        self.long_url = long_url
        self.current_time = current_time
        self.expiration = expiration_minutes * 60 # seconds


class HashGenerator(ABC):
    @abstractmethod
    def hash_generator(self, user_id: str, url: str) -> str:
        pass
        
class BasicHashFunction(HashGenerator):
    def hash_generator(self, user_id: str, url: str):
        hash_processor = user_id + url
        return hex(hash(hash_processor))[2:]

class StorageService:
    def __init__(self):
        self.storage = {} # hashKey -> UserData
    
    def store_data(self, hash_key: str, user_data: UserData):
        self.storage[hash_key] = user_data

    def get_data(self, hash_key: str):
        return self.storage.get(hash_key)

    def clean_urls(self):
        current_time = time.time()
        expired_keys = [
            key for key, value in self.storage.items()
            if current_time > value.current_time + value.expiration
        ]
        for key in expired_keys:
            del self.storage[key]


class UrlService:
    def __init__(self, storage_service: StorageService, hash_generator: HashGenerator):
        self.storage_service = storage_service
        self.hash_generator = hash_generator
        self.base_url = "https://tiny.url/"

    def create_url(self, user_id: str, url: str, minutes: int) -> str:
        hash_key = self.hash_generator.hash_generator(user_id, url)

        existing = self.storage_service.get_data(hash_key)
        if existing:
            return existing.short_url

        short_url = self.base_url + hash_key
        user_data = UserData(
            user_id=user_id,
            short_url=short_url,
            long_url=url,
            current_time=time.time(),
            expiration_minutes=minutes
        )
        self.storage_service.store_data(hash_key, user_data)
        return short_url

    def view_url(self, short_url: str):
        hash_key = short_url.split("/")[-1]
        user_data = self.storage_service.get_data(hash_key)

        if not user_data:
            return None

        if time.time() > user_data.current_time + user_data.expiration:
            return None

        return user_data.long_url


if __name__ == "__main__":
    storage_service = StorageService()
    hash_generator = BasicHashFunction()
    url_service = UrlService(storage_service, hash_generator)

    # Akshay
    short_url = url_service.create_url(
        "youtube",
        "https://www.youtube.com/",
        1
    )
    print(short_url)
    print(url_service.view_url(short_url))

    # Akshay1
    short_url = url_service.create_url(
        "youtube",
        "https://www.youtube.com/",
        1
    )
    print(short_url)
    print(url_service.view_url(short_url))

    # Wait until expiration
    time.sleep(62)

    if url_service.view_url(short_url) is None:
        print("Expired")
