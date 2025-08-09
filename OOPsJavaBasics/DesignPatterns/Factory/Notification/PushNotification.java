package OOPsJavaBasics.DesignPatterns.Factory.Notification;

public class PushNotification implements Notification {
    public void SendNotification(String message) {
        System.out.println("Push Notification: " + message);
    }
}
