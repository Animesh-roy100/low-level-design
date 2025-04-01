package oopsJava.DesignPatterns.Factory.Notification;

public class EmailNotification implements Notification {
    public void SendNotification(String message) {
        System.out.println("Email Notification: " + message);
    }
}
