package LLDJava.NotificationService;

import java.time.LocalDateTime;
import java.util.*;

enum Channel {
    EMAIL, SMS, PUSH
}

// Notification interface ---------------------------------------------

interface Notification {
    Channel getChannel();
    String getRecipient();
    String getContent();
}

final class Email implements Notification {
    private final String toEmail;
    private final String body;
    private final String subject;

    public Email(String toEmail, String body, String subject) {
        this.toEmail = toEmail;
        this.body = body;
        this.subject = subject;
    }

    @Override
    public Channel getChannel() { return Channel.EMAIL; }
    @Override
    public String getRecipient() { return toEmail; }
    @Override
    public String getContent() { return body; }

    public String getSubject() { return subject; }
}

class Push implements Notification {
    private final String toDeviceId;
    private final String title;
    private final String payload;

    public Push(String toDeviceId, String title, String payload) {
        this.toDeviceId = toDeviceId;
        this.title = title;
        this.payload = payload;
    }

    @Override
    public Channel getChannel() { return Channel.PUSH; }
    @Override
    public String getRecipient() { return toDeviceId; }
    @Override
    public String getContent() { return payload; }

    public String getTitle() { return title; }
}

class SMS implements Notification {
    private final String toPhoneNumber;
    private final String message;

    public SMS(String toPhoneNumber, String message) {
        this.toPhoneNumber = toPhoneNumber;
        this.message = message;
    }

    @Override
    public Channel getChannel() { return Channel.SMS; }
    @Override
    public String getRecipient() { return toPhoneNumber; }
    @Override
    public String getContent() { return message; }
}

// Notification Sender Interface ---------------------------------

interface NotificationSender {
    void send(Notification notification);
}

interface ScheduleNotificationSender extends NotificationSender {
    void schedule(Notification notification, LocalDateTime dateTime);
}

class EmailNotificationSender implements ScheduleNotificationSender {
    @Override
    public void send(Notification notification) {
       System.out.println("Sending Email to " + notification.getRecipient() + " with content: " + notification.getContent());
    }
    @Override
    public void schedule(Notification notification, LocalDateTime dateTime) {
        System.out.println("Sending Email to " + notification.getRecipient() + " with content: " + notification.getContent() + " at " + dateTime);
    }
}

class PushNotificationSender implements ScheduleNotificationSender {
    @Override
    public void send(Notification notification) {
        System.out.println("Sending Push to " + notification.getRecipient() + " with content: " + notification.getContent());
    }
    @Override
    public void schedule(Notification notification, LocalDateTime dateTime) {
        System.out.println("Sending Push to " + notification.getRecipient() + " with content: " + notification.getContent() + " at " + dateTime);
    }
}

class SMSNotificationSender implements ScheduleNotificationSender {
    @Override
    public void send(Notification notification) {
        System.out.println("Sending SMS to " + notification.getRecipient() + " with content: " + notification.getContent());
    }
    @Override
    public void schedule(Notification notification, LocalDateTime dateTime) {
        System.out.println("Sending SMS to " + notification.getRecipient() + " with content: " + notification.getContent() + " at " + dateTime);
    }
}

// Notification Sender Factory ---------------------------------

interface NotificationSenderFactory {
    Optional<NotificationSender> getSender(Channel channel);
    Optional<ScheduleNotificationSender> getSenderWithSchedule(Channel channel);
}

class DefaultNotificationSenderFactory implements NotificationSenderFactory {
    private final Map<Channel, NotificationSender> senderMap;
    private final Map<Channel, ScheduleNotificationSender> scheduleSenderMap;

    public DefaultNotificationSenderFactory() {
        this.senderMap = new HashMap<>();
        this.scheduleSenderMap = new HashMap<>();

        // Register senders externally
        EmailNotificationSender emailSender = new EmailNotificationSender();
        SMSNotificationSender smsSender = new SMSNotificationSender();
        PushNotificationSender pushSender = new PushNotificationSender();

        senderMap.put(Channel.EMAIL, emailSender);
        senderMap.put(Channel.SMS, smsSender);
        senderMap.put(Channel.PUSH, pushSender);

        scheduleSenderMap.put(Channel.EMAIL, emailSender);
        scheduleSenderMap.put(Channel.SMS, smsSender);
        scheduleSenderMap.put(Channel.PUSH, pushSender);
    }

    @Override
    public Optional<NotificationSender> getSender(Channel channel) {
        return Optional.ofNullable(senderMap.get(channel));
    }
    @Override
    public Optional<ScheduleNotificationSender> getSenderWithSchedule(Channel channel) {
        return Optional.ofNullable(scheduleSenderMap.get(channel));
    }
}

// Notification Dispatcher ---------------------------------

class NotificationDispatcher {
    private final NotificationSenderFactory senderFactory;

    public NotificationDispatcher(NotificationSenderFactory senderFactory) {
        this.senderFactory = senderFactory;
    }

    public void dispatch(Notification notification) {
        Optional<NotificationSender> senderOpt = senderFactory.getSender(notification.getChannel());
        if (senderOpt.isPresent()) {
            senderOpt.get().send(notification);
        } else {
            System.out.println("No sender found for channel: " + notification.getChannel());
        }
    }

    public void scheduleDispatch(Notification notification, LocalDateTime dateTime) {
        Optional<ScheduleNotificationSender> senderOpt = senderFactory.getSenderWithSchedule(notification.getChannel());
        if (senderOpt.isPresent()) {
            senderOpt.get().schedule(notification, dateTime);
        } else {
            System.out.println("No scheduled sender found for channel: " + notification.getChannel());
        }
    }
}


public class Main {
    public static void main(String[] args) {
        NotificationSenderFactory senderFactory = new DefaultNotificationSenderFactory();
        NotificationDispatcher dispatcher = new NotificationDispatcher(senderFactory);

        Notification email = new Email("animesh@email.com", "Hello Animesh!", "Welcome");
        dispatcher.dispatch(email);
        dispatcher.scheduleDispatch(email, LocalDateTime.now().plusHours(2));

        Notification sms = new SMS("123456789", "Hello Animesh!");
        dispatcher.dispatch(sms);
        dispatcher.scheduleDispatch(sms, LocalDateTime.now().plusMinutes(30));

        Notification push = new Push("123456789", "New Message", "You have a new message.");
        dispatcher.dispatch(push);
        dispatcher.scheduleDispatch(push, LocalDateTime.now().plusMinutes(30));
    }
}
