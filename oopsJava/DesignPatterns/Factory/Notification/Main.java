package oopsJava.DesignPatterns.Factory.Notification;

public class Main {
    public static void main(String[] args) {
        Notification smsNotification = NotificationFactory.createNotification("sms");
        smsNotification.SendNotification("Hello SMS");

        Notification emailNotification = NotificationFactory.createNotification("email");
        emailNotification.SendNotification("Hello Email");

        Notification pushNotification = NotificationFactory.createNotification("push");
        pushNotification.SendNotification("Hello Push");
    }    
}
