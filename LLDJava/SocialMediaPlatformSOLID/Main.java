package LLDJava.SocialMediaPlatformSOLID;
import java.util.*;

// Entity Classes ----------------------------------------- 
class User {
    int userId;
    String name;

    public User(int userId, String name) {
        this.userId = userId;
        this.name = name;
    }
}

class Post {
    int postId;
    String content;
    int userId;

    public Post(int postId, String content, int userId) {
        this.postId = postId;
        this.content = content;
        this.userId = userId;
    }
}

class Comment {
    int commentId;
    String comment;
    int postId;
    int userId;

    public Comment(int commentId, String comment, int postId, int userId) {
        this.commentId = commentId;
        this.comment = comment;
        this.postId = postId;
        this.userId = userId;
    }
}

// Repository classes ---------------------------------------
// Following Single Responsibility Principle (SRP)
class UserRepository {
    private Map<Integer, User> users = new HashMap<>();

    public void addUser(int userId, String name) {
        users.put(userId, new User(userId, name));
    }

    public User getUser(int userId) {
        return users.get(userId);
    }

    public Collection<User> getAllUsers() {
        return users.values();
    }
}

class PostRepository {
    private Map<Integer, Post> posts = new HashMap<>();

    public void addPost(int postId, String content, int userId) {
        posts.put(postId, new Post(postId, content, userId));
    }

    public Post getPost(int postId) {
        return posts.get(postId);
    }

    public List<Post> getPostsByUserId(int userId) {
        List<Post> userPosts = new ArrayList<>();
        for (Post post : posts.values()) {
            if (post.userId == userId) {
                userPosts.add(post);
            }
        }
        return userPosts;
    }

    public List<Post> getAllPosts() {
        return new ArrayList<>(posts.values());
    }
}

class CommentRepository {
    private Map<Integer, Comment> comments = new HashMap<>();

    public void addComment(int commentId, String comment, int postId, int userId) {
        comments.put(commentId, new Comment(commentId, comment, postId, userId));
    }

    public List<Comment> getCommentsByPostId(int postId) {
        List<Comment> postComments = new ArrayList<>();
        for (Comment comment : comments.values()) {
            if (comment.postId == postId) {
                postComments.add(comment);
            }
        }
        return postComments;
    }
}

// Service classes ---------------------------------------
class UserService {
    private UserRepository userRepository;

    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    public void addUser(int userId, String name) {
        System.out.println("Added new user " + userId);
        userRepository.addUser(userId, name);
    }

    public User getUser(int userId) {
        return userRepository.getUser(userId);
    }

    public void listUserNames() {
        for (User user : userRepository.getAllUsers()) {
            System.out.println(user.name);
        }
    }
}

class PostService {
    private PostRepository postRepository;
    private UserRepository userRepository;

    public PostService(PostRepository postRepository, UserRepository userRepository) {
        this.postRepository = postRepository;
        this.userRepository = userRepository;
    }

    public void addPost(int userId, String content, int postId) {
        if (userRepository.getUser(userId) != null) {
            System.out.println("Added new post " + userId);
            postRepository.addPost(postId, content, userId);
        } else {
            System.out.println("User not found");
        }
    }

    public List<Post> getPostsByUserId(int userId) {
        return postRepository.getPostsByUserId(userId);
    }

    public List<Post> getAllPosts() {
        return postRepository.getAllPosts();
    }
}

class CommentService {
    private CommentRepository commentRepository;
    private PostRepository postRepository;
    private UserRepository userRepository;

    public CommentService(CommentRepository commentRepository, PostRepository postRepository, UserRepository userRepository) {
        this.commentRepository = commentRepository;
        this.postRepository = postRepository;
        this.userRepository = userRepository;
    }

    public void addComment(int commenterId, String content, int postId, int commentId) {
        if (userRepository.getUser(commenterId) == null) {
            System.out.println("Commenter not found");
            return;
        }
        if (postRepository.getPost(postId) == null) {
            System.out.println("Post not found");
            return;
        }
        System.out.println("Added new comment " + commenterId);
        commentRepository.addComment(commentId, content, postId, commenterId);
    }

    public List<Comment> getCommentsByPostId(int postId) {
        return commentRepository.getCommentsByPostId(postId);
    }
}


class SocialMediaPlatform {
    private UserService userService;
    private PostService postService;
    private CommentService commentService;

    public SocialMediaPlatform(UserService userService, PostService postService, CommentService commentService) {
        this.userService = userService;
        this.postService = postService;
        this.commentService = commentService;
    }

    public void addUser(int userId, String name) {
        userService.addUser(userId, name);
    }

    public void getUsersName() {
        userService.listUserNames();
    }

    public User getUser(int userId) {
        return userService.getUser(userId);
    }

    public void addPost(int userId, String content, int postId) {
        postService.addPost(userId, content, postId);
    }

    public void addComment(int commenterId, String content, int postId, int commentId) {
        commentService.addComment(commenterId, content, postId, commentId);
    }

    public void getPosts() {
        for (Post post : postService.getAllPosts()) {
            User user = userService.getUser(post.userId);
            System.out.println("Post by user " + user.name + ": " + post.content + " PostID " + post.postId);
        }
    }

    public void showPosts() {
        for (Post post : postService.getAllPosts()) {
            User user = userService.getUser(post.userId);
            System.out.println("Post by " + user.name + ": " + post.content);
            for (Comment comment : commentService.getCommentsByPostId(post.postId)) {
                User commenter = userService.getUser(comment.userId);
                System.out.println("Commented by User " + commenter.name + ": " + comment.comment);
            }
        }
    }
}

public class Main {
    public static void main(String[] args) {
        System.out.println("Hello");

        UserRepository userRepository = new UserRepository();
        PostRepository postRepository = new PostRepository();
        CommentRepository commentRepository = new CommentRepository();

        UserService userService = new UserService(userRepository);
        PostService postService = new PostService(postRepository, userRepository);
        CommentService commentService = new CommentService(commentRepository, postRepository, userRepository);

        SocialMediaPlatform platform = new SocialMediaPlatform(userService, postService, commentService);

        platform.addUser(1, "Alice");
        platform.addUser(2, "Bob");

        platform.getUsersName();

        platform.addPost(1, "My First Post", 1);

        platform.getPosts();

        platform.addComment(2, "Welcome", 1, 1); // Added commentId for uniqueness

        platform.showPosts();
    }
}
