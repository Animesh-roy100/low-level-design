package LLDJava.SocialMediaPlatform;
import java.util.*;

class User {
    int userId;
    String name;
    List<Post> posts;
    
    public User(int id, String name){
        this.userId = id;
        this.name = name;
        this.posts = new ArrayList<>();
    }

    public void addPost(String content, int postId){
        posts.add(new Post(content, postId));
    }

    public List<Post> getPosts(){
        return posts;
    }
}

class Post {
    int postId;
    String content;
    List<Comment> comments;

    public Post(String content, int postId) {
        this.content = content;
        this.postId = postId;
        this.comments = new ArrayList<>();
    }

    public void addComment(String comment, int userId) {
        comments.add(new Comment(comment, userId));
    }

    public List<Comment> getComments() {
        return comments;
    }

}

class Comment {
    int userId;
    String comment;

    public Comment(String comment, int userId) {
        this.comment = comment;
        this.userId = userId;
    }
}

class SocialMediaPlatform {
    // List<User> users;
    HashMap<Integer, User> users;

    public SocialMediaPlatform(){
        this.users = new HashMap<>();
    }

    public void addUser(int userId, String name){
        System.out.println("added new user" + userId);
        users.put(userId, new User(userId, name));
    }

    public void getUsersName() {
        for(User user: users.values()) {
            System.out.println(user.name);
        }
    }

    public User getUser(int userId) {
        return users.get(userId);
    }

    public void addPost(int userId, String content, int postId) {
        System.out.println("added new post" + userId);
        User user = users.get(userId);

        if(user != null) {
            user.addPost(content, postId);
        }
    }

    public void addComment(int commenterId, String content, int postId){
        System.out.println("added new comment" + commenterId);
        User commenter = getUser(commenterId);

        System.out.println("commenter ");
        // Got error in this
        // commenter.getPosts().get(postId).addComment(content, commenterId);
        List<Post> posts = commenter.getPosts();
        for (Post post: posts) {
            if(post.postId == postId) {
                post.addComment(content, commenterId);
            }
        }
    }

    public void getPosts() {
        for (User user : users.values()) {
            List<Post> posts = user.getPosts();

            for(Post post : posts) {
                System.out.println("Post by user " + user.name + ": " + post.content + " PostID " + post.postId );
            }
        }
    }

    public void showPosts() {
        for (User user: users.values()){
            for(Post post: user.getPosts()) {
                System.out.println("Post by " + user.name + ": " + post.content);
                for (Comment comment: post.comments) {
                    User commenter = getUser(comment.userId);
                    System.out.println("Commented by User " + commenter.name + ": " + comment.comment);
                }
            }
        }
    }
}

public class Main {
    public static void main(String[] args) {
        System.out.println("Hello");
        SocialMediaPlatform platform = new SocialMediaPlatform();

        platform.addUser(1, "Alice");
        platform.addUser(2, "Bob");

        platform.getUsersName();

        platform.addPost(1, "My First Post", 1);

        platform.getPosts();

        platform.addComment(2, "Welcome", 1);

        platform.showPosts();
    }
}