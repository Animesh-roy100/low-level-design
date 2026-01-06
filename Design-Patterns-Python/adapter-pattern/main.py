from abc import ABC, abstractmethod

# The Adapter Pattern allows incompatible interfaces to work together by wrapping an existing class with a new interface.
# Convert the interface of a class into another interface the client expects.

class PaymentProcessor(ABC):
    @abstractmethod
    def pay(self, amount: float) -> None:
        pass

# Incompatible Class
class PayPalService:
    def make_payment(self, total_amount: float) -> None:
        print(f"Paid {total_amount} using PayPal")
    
# Adapter
class PayPalAdapter(PaymentProcessor):
    def __init__(self, paypal_service: PayPalService):
        self.paypal_service = paypal_service
    
    def pay(self, amount: float) -> None:
        self.paypal_service.make_payment(amount)


def checkout(processor: PaymentProcessor, amount: float):
    processor.pay(amount)

paypal = PayPalService()
adapter = PayPalAdapter(paypal)

checkout(adapter, 500)
    