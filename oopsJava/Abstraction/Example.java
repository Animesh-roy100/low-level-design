package oopsJava.Abstraction;

// Abstract class representing a generic Car
abstract class Car {
    // Abstract method to start the car
    abstract void start();
}

// Concrete class Audi extending Car
class Audi extends Car {
    // Implementation of start method for Audi
    void start() {
        System.out.println("Audi car starts");
    }
}

// Concrete class BMW extending Car
class BMW extends Car {
    // Implementation of start method for BMW
    void start() {
        System.out.println("BMW car starts");
    }
}

// Drive class that can drive any Car
class Drive {
    // Method to drive a car by calling its start method
    void drive(Car c) {
        c.start();
    }
}

// Main class to execute the program
public class Example {
    public static void main(String[] args) {
        Drive d = new Drive();          // Create a Drive instance
        d.drive(new Audi());            // Drive an Audi
        d.drive(new BMW());             // Drive a BMW
    }
}
