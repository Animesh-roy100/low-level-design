package OOPsJavaBasics.DesignPatterns.Factory.Notification;

public class SMSNotification implements Notification {
    public void SendNotification(String message) {
        System.out.println("SMS Notification: " + message);
    }
}