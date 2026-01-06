from abc import ABC, abstractmethod 

# Step 1: Define an interface for payment strategy
class PaymentStrategy(ABC):
    @abstractmethod
    def pay(self, amount):
        pass

# Step 2: Implement specific payment strategies
class CreditCardPayment(PaymentStrategy):
    def pay(self, amount):
        print(f"Processing ${amount} payment using Credit Card")

class DebitCardPayment(PaymentStrategy):
    def pay(self, amount):
        print(f"Processing ${amount} payment using Debit Card")

class NetBankingPayment(PaymentStrategy):
    def pay(self, amount):
        print(f"Processing ${amount} payment using Net Banking")

class PaymentProcessor:
    def __init__(self, strategy: PaymentStrategy):
        self.strategy = strategy
    
    def set_strategy(self, strategy: PaymentStrategy):
        self.strategy = strategy
    
    def process_payment(self, amount):
        self.strategy.pay(amount)

processor = PaymentProcessor(CreditCardPayment())
processor.process_payment(100)

processor = PaymentProcessor(DebitCardPayment())
processor.process_payment(100)

processor = PaymentProcessor(NetBankingPayment())
processor.process_payment(100)