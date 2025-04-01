package oopsJava.DesignPatterns.Factory.Notification;

public class NotificationFactory {
    public static Notification createNotification(String notificationType) {
        switch (notificationType) {
            case "email":
                return new EmailNotification();
            case "sms": 
                return new SMSNotification();
            case "push":
                return new PushNotification();
            default:
                throw new IllegalArgumentException("Unknown notification type: " + notificationType);
        }
    }
}
