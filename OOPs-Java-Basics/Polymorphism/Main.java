package oopsJava.Polymorphism;

public class Main {
    public static void main(String[] arg) {
        Shape[] shapes = new Shape[2];
        shapes[0] = new Circle(5);
        shapes[1] = new Rectangle(3, 4);

        for(Shape shape: shapes) {
            System.out.println("Area of " + shape.getClass().getSimpleName() + " is " + shape.area());
        }
    }
}
