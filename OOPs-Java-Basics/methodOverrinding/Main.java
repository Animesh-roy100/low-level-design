package oopsJava.methodOverrinding;

public class Main {
    public static void main(String[] args) {
        Animal animal = new Animal();
        animal.makeSound();

        Dog dog = new Dog();
        dog.makeSound();
    }    
}
