package LLDJava.WhatsAppGroup;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;

class User {
    private String userId;
    private String name;
    private boolean onlineStatus;
    private boolean showOnlineStatus;

    public User(String userId, String name) {
        this.userId = userId;
        this.name = name;
        this.onlineStatus = false;
        this.showOnlineStatus = true;
    }

    public String getUserId() {
        return userId;
    }

    public String getName() {
        return name;
    }

    public boolean isOnlineStatus() {
        return onlineStatus;
    }

    public void setOnlineStatus(boolean onlineStatus) {
        this.onlineStatus = onlineStatus;
    }

    public boolean isShowOnlineStatus() {
        return showOnlineStatus;
    }

    public void setShowOnlineStatus(boolean showOnlineStatus) {
        this.showOnlineStatus = showOnlineStatus;
    }
}

class Message {
    private String messageId;
    private User sender;
    private String content;
    private Date timestamp;

    public Message(String messageId, User sender, String content) {
        this.messageId = messageId;
        this.sender = sender;
        this.content = content;
        this.timestamp = new Date();
    }

    public String getMessageId() {
        return messageId;
    }

    public User getSender() {
        return sender;
    }

    public String getContent() {
        return content;
    }

    public Date getTimestamp() {
        return timestamp;
    }
}

class Group {
    private String groupId;
    private String name;
    private List<User> members;
    private List<User> admins;
    private List<Message> messages;

    public Group(String groupId, String name) {
        this.groupId = groupId;
        this.name = name;
        this.members = new ArrayList<>();
        this.admins = new ArrayList<>();
        this.messages = new ArrayList<>();
    }

    public void addMember(User user) {
        if (!members.contains(user)) {
            members.add(user);
        }
    }

    public void removeMember(User user) {
        members.remove(user);
        admins.remove(user);
    }

    public void addAdmin(User user) {
        if (members.contains(user) && !admins.contains(user)) {
            admins.add(user);
        }
    }

    public void removeAdmin(User user) {
        admins.remove(user);
    }

    public void sendMessage(Message message) {
        if (members.contains(message.getSender())) {
            messages.add(message);
        }
    }

    public List<User> getActiveUsers() {
        List<User> activeUsers = new ArrayList<>();
        for (User user : members) {
            if (user.isOnlineStatus() && user.isShowOnlineStatus()) {
                activeUsers.add(user);
            }
        }
        return activeUsers;
    }

    public String getGroupId() {
        return groupId;
    }

    public String getName() {
        return name;
    }

    public List<User> getMembers() {
        return members;
    }

    public List<User> getAdmins() {
        return admins;
    }

    public List<Message> getMessages() {
        return messages;
    }
}

public class WhatsAppGroup {
    public static void main(String[] args) {
        User user1 = new User("u1", "Alice");
        user1.setOnlineStatus(true);
        User user2 = new User("u2", "Bob");
        user2.setOnlineStatus(true);
        user2.setShowOnlineStatus(false);

        Group group = new Group("g1", "Friends");
        group.addMember(user1);
        group.addMember(user2);
        group.addAdmin(user1);

        Message msg = new Message("m1", user1, "Hello, group!");
        group.sendMessage(msg);

        List<User> activeUsers = group.getActiveUsers();
        System.out.println("Active users:");
        for (User user : activeUsers) {
            System.out.println(user.getName());
        }
    }
}