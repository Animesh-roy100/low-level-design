package oopsJava.Inheritance;

class Engine {
    int horsePower;

    public Engine(int horsePower) {
        this.horsePower = horsePower;
    }

    public void start() {
        System.out.println("Engine started with horsepower: " + horsePower);
    }
}

class Car extends Engine {  // Car "is-a" Engine
    String brandName;

    public Car(int horsePower, String brandName) {
        super(horsePower);
        this.brandName = brandName;
    }

    public void showDetails() {
        System.out.println("Car brand: " + brandName);
        start();  // Inherited start method from Engine
    }
}

public class Main {
    public static void main(String[] args) {
        Car myCar = new Car(300, "Tesla");
        myCar.showDetails();  // Accessing inherited methods and fields
    }
}
