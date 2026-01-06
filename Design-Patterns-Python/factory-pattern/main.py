class Notification():
    def send_notification(self, message: str) -> None:
        raise NotImplementedError("This method should be overridden by subclasses")


class EmailNotification(Notification):
    def send_notification(self, message: str) -> None:
        print(f"Email notification sent: {message}")
        
    
class SMSNotification(Notification):
    def send_notification(self, message: str) -> None:
        print(f"SMS notification sent: {message}")

class PushNotification(Notification):
    def send_notification(self, message: str) -> None:
        print(f"Push notification sent: {message}")
    

def NotificationFactory(notification_type: str) -> Notification:
    if notification_type == "email":
        return EmailNotification()
    elif notification_type == "sms":
        return SMSNotification()
    elif notification_type == "push":
        return PushNotification();
    else:
        raise ValueError(f"Unknown notification type: {notification_type}")


def main():
    try:
        notification_type = "email"
        notification = NotificationFactory(notification_type)
        notification.send_notification("Order placed successfully")
    except ValueError as e:
        print(e)

if __name__ == "__main__":
    main()