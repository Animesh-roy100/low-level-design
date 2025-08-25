package LLDJava.CarRentalSystem;

import java.time.LocalDateTime;
import java.util.*;
import java.util.concurrent.atomic.AtomicInteger;

enum ReservationStatus {
    INITIATED, SCHEDULED, INPROGRESS, COMPLETED, CANCELLED;
}

enum PaymentMode {
    ONLINE, CASH;
}

enum CarType {
    MINIVAN, SUV, SEDAN, SPORT, HATCHBACK;
}

enum PricingType {
    HOURLY, DAILY, WEEKLY;
}

// Observer Pattern Implementation
interface ReservationObserver {
    void update(Reservation reservation, String message);
}

class ReservationNotifier {
    private List<ReservationObserver> observers = new ArrayList<>();
    
    public void addObserver(ReservationObserver observer) {
        observers.add(observer);
    }
    
    public void removeObserver(ReservationObserver observer) {
        observers.remove(observer);
    }
    
    public void notifyObservers(Reservation reservation, String message) {
        for (ReservationObserver observer : observers) {
            observer.update(reservation, message);
        }
    }
}

// Notification Services 
class EmailNotificationService implements ReservationObserver {
    public void update(Reservation reservation, String message) {
        User user = reservation.getUser();
        System.out.println("Email sent to " + user.getEmail() + ": " + message);
    }
}

class SMSNotificationService implements ReservationObserver {
    public void update(Reservation reservation, String message) {
        User user = reservation.getUser();
        System.out.println("SMS sent to " + user.getPhone() + ": " + message);
    }
}

class PushNotificationService implements ReservationObserver {
    public void update(Reservation reservation, String message) {
        User user = reservation.getUser();
        System.out.println("Push notification sent to user " + user.getName() + ": " + message);
    }
}

// classes -----------------------------------------------
class Location {
    private String address;
    private int zipCode;
    private String city;
    private String state;
    private String country;

    public Location(String address, int zipCode, String city, String state, String country) {
        this.address = address;
        this.zipCode = zipCode;
        this.city = city;
        this.state = state;
        this.country = country;
    }

    // Getters and setters
    public String getAddress() { return address; }
    public void setAddress(String address) { this.address = address; }
    
    public int getZipCode() { return zipCode; }
    public void setZipCode(int zipCode) { this.zipCode = zipCode; }
    
    public String getCity() { return city; }
    public void setCity(String city) { this.city = city; }
    
    public String getState() { return state; }
    public void setState(String state) { this.state = state; }
    
    public String getCountry() { return country; }
    public void setCountry(String country) { this.country = country; }
}

class User {
    private String id;
    private String name;
    private String email;
    private String phone;
    private String password;

    public User(String id, String name, String email, String phone, String password) {
        this.id = id;
        this.name = name;
        this.email = email;
        this.phone = phone;
        this.password = password;
    }

    // Getters and setters
    public String getId() { return id; }
    public String getEmail() { return email; }
    public String getPhone() { return phone; }
    public String getName() { return name; }
}

class Vehicle {
    private int id;
    private String make;
    private String model;
    private int year;
    private String numberPlate;
    private double pricePerHour;

    public Vehicle(int id, String make, String model, int year, String numberPlate, double pricePerHour) {
        this.id = id;
        this.make = make;
        this.model = model;
        this.year = year;
        this.numberPlate = numberPlate;
        this.pricePerHour = pricePerHour;
    }

    // Getters and setters
    public int getId() { return id; }
    public void setId(int id) { this.id = id; }

    public double getPricePerHour() { return pricePerHour; }
    public void setPricePerHour(double pricePerHour) { this.pricePerHour = pricePerHour; }

    public String toString() {
        return year + " " + make + " " + model + " (" + numberPlate + ")";
    }
}

class Car extends Vehicle {
    private CarType type;

    public Car(int id, String make, String model, int year, String numberPlate, double pricePerHour, CarType type) {
        super(id, make, model, year, numberPlate, pricePerHour);
        this.type = type;
    }

    // Getters and setters
    public CarType getType() { return type; }
    public void setType(CarType type) { this.type = type; }
}

class Reservation {
    private int reservationId;
    private User user;
    private Vehicle vehicle;
    private LocalDateTime startDateTime;
    private LocalDateTime endDateTime;
    private ReservationStatus reservationStatus;
    private Location location;
    private PricingType pricingType;

    public Reservation(int reservationId, User user, Vehicle vehicle, LocalDateTime startDateTime, LocalDateTime endDateTime, ReservationStatus reservationStatus, Location location, PricingType pricingType) {
        this.reservationId = reservationId;
        this.user = user;
        this.vehicle = vehicle;
        this.startDateTime = startDateTime;
        this.endDateTime = endDateTime;
        this.reservationStatus = reservationStatus;
        this.location = location;
        this.pricingType = pricingType;
    }

    // Getters and setters
    public int getReservationId() { return reservationId; }
    public void setReservationId(int reservationId) { this.reservationId = reservationId;}

    public ReservationStatus getReservationStatus() { return reservationStatus; }
    public void setReservationStatus(ReservationStatus reservationStatus) { this.reservationStatus = reservationStatus; }

    public Vehicle getVehicle() { return vehicle; }
    public void setVehicle(Vehicle vehicle) { this.vehicle = vehicle; }

    public User getUser() { return user; }
    public void setUser(User user) { this.user = user; }

    public long getDurationInHours() {
        return java.time.Duration.between(startDateTime, endDateTime).toHours();
    }

    public PricingType getPricingType() { return pricingType; }
    public void setPricingType(PricingType pricingType) { this.pricingType = pricingType; }
}

class PaymentDetails {
    private int paymentId;
    private String paymentReference;
    private LocalDateTime dateOfPayment;
    private boolean refundMode;
    private PaymentMode paymentMode;
    private double amount;

    public PaymentDetails(int paymentId, String paymentReference, LocalDateTime dateOfPayment, boolean refundMode, PaymentMode paymentMode, double amount) {
        this.paymentId = paymentId;
        this.paymentReference = paymentReference;
        this.dateOfPayment = dateOfPayment;
        this.refundMode = refundMode;
        this.paymentMode = paymentMode;
        this.amount = amount;
    }

    // Getters and setters
    public int getPaymentId() { return paymentId; }
    public void setPaymentId(int paymentId) { this.paymentId = paymentId; }
}

// Vehicle Inventory Management 
class VehicleInventoryManagement {
    private List<Vehicle> vehicles = new ArrayList<>();

    public List<Vehicle> getVehicles() {
        return new ArrayList<>(vehicles); // Return a copy to avoid external modification
    }
    
    public void setVehicles(List<Vehicle> vehicles) { this.vehicles = new ArrayList<>(vehicles);}
    public void addVehicle(Vehicle vehicle) {vehicles.add(vehicle);}
    public void removeVehicle(int vehicleId) { vehicles.removeIf(v -> v.getId() == vehicleId);}

    public Vehicle getVehicleById(int vehicleId) {
        for (Vehicle v : vehicles) {
            if (v.getId() == vehicleId) {
                return v;
            }
        }
        return null;
    }

    public List<Vehicle> getVehiclesByType(CarType type) {
        List<Vehicle> result = new ArrayList<>();
        for (Vehicle vehicle : vehicles) {
            if (vehicle instanceof Car) {
                Car car = (Car) vehicle;
                if (car.getType() == type) {
                    result.add(vehicle);
                }
            }
        }
        return result;
    }
}

// Reservation Manager
class ReservationManager {
    private List<Reservation> reservations = new ArrayList<>();
    private AtomicInteger nextId = new AtomicInteger(1);
    private ReservationNotifier notifier = new ReservationNotifier();
    
    public void addObserver(ReservationObserver observer) {
        notifier.addObserver(observer);
    }

    public Reservation createReservation(Reservation reservation) {
        reservation.setReservationId(nextId.getAndIncrement());
        reservations.add(reservation);
        notifier.notifyObservers(reservation, "Reservation created with ID: " + reservation.getReservationId());
        return reservation;
    }

    public Reservation completeReservation(int reservationId) {
        for (Reservation r: reservations) {
            if (r.getReservationId() == reservationId) {
                r.setReservationStatus(ReservationStatus.COMPLETED);
                notifier.notifyObservers(r, "Reservation completed: " + reservationId);
                return r;
            }
        }
        return null;
    }

    public Reservation changeStatusToScheduled(int reservationId) {
        for (Reservation reservation : reservations) {
            if (reservation.getReservationId() == reservationId) {
                reservation.setReservationStatus(ReservationStatus.SCHEDULED);
                notifier.notifyObservers(reservation, "Reservation scheduled: " + reservationId);
                return reservation;
            }
        }
        return null;
    }

    public Reservation changeStatusToProgress(int reservationId) {
        for (Reservation reservation : reservations) {
            if (reservation.getReservationId() == reservationId) {
                reservation.setReservationStatus(ReservationStatus.INPROGRESS);
                notifier.notifyObservers(reservation, "Reservation in progress: " + reservationId);
                return reservation;
            }
        }
        return null;
    }

    public Reservation cancelReservation(int reservationId) {
        for (Reservation reservation : reservations) {
            if (reservation.getReservationId() == reservationId) {
                reservation.setReservationStatus(ReservationStatus.CANCELLED);
                notifier.notifyObservers(reservation, "Reservation cancelled: " + reservationId);
                return reservation;
            }
        }
        return null;
    }

    public Reservation getReservationById(int reservationId) {
        for (Reservation reservation : reservations) {
            if (reservation.getReservationId() == reservationId) {
                return reservation;
            }
        }
        return null;
    }
}

class Store {
    private int storeId;
    private VehicleInventoryManagement inventoryManagement;
    private ReservationManager reservationManager;
    private Location location;

    public Store(int storeId, Location location) {
        this.storeId = storeId;
        this.location = location;
        this.inventoryManagement = new VehicleInventoryManagement();
        this.reservationManager = new ReservationManager();
    }
    
   public List<Vehicle> getVehiclesByType(CarType type) {
        return inventoryManagement.getVehiclesByType(type);
    }

    public List<Vehicle> getVehicles() {
        return inventoryManagement.getVehicles();
    }

    public void addVehicle(Vehicle vehicle) {
        inventoryManagement.addVehicle(vehicle);
    }
    
    public void addReservationObserver(ReservationObserver observer) {
        reservationManager.addObserver(observer);
    }

    public Reservation updateOrCreateReservation(Reservation reservation, ReservationStatus status) {
        reservation.setReservationStatus(status);
        if (reservation.getReservationId() == 0) {
            return reservationManager.createReservation(reservation);
        }

        // Update existing reservation
        Reservation existingReservation = reservationManager.getReservationById(reservation.getReservationId());
        if (existingReservation != null) {
            existingReservation.setReservationStatus(status);
            String statusMessage = "Reservation status updated to: " + status;
            return existingReservation;
        }
        return null;
    } 
 

    // Getters and setters
    public int getStoreId() { return storeId; }
    public void setStoreId(int storeId) { this.storeId = storeId; }
    
    public VehicleInventoryManagement getInventoryManagement() { return inventoryManagement; }
    public void setInventoryManagement(VehicleInventoryManagement inventoryManagement) { 
        this.inventoryManagement = inventoryManagement; 
    }
    
    public Location getLocation() { return location; }
    public void setLocation(Location location) { this.location = location; }
    
    public ReservationManager getReservationManager() { return reservationManager; }
    public void setReservationManager(ReservationManager reservationManager) { 
        this.reservationManager = reservationManager; 
    }
}

class VehicleRentalSystem {
    private List<Store> stores = new ArrayList<>();
    private List<User> users = new ArrayList<>();
    
    public Store getStoreByLocation(Location location) {
        for (Store store : stores) {
            Location storeLocation = store.getLocation();
            if (storeLocation.getAddress().equals(location.getAddress()) && 
                storeLocation.getCity().equals(location.getCity())) {
                return store;
            }
        }
        return null;
    }

    // Getters and setters
    public void addStore(Store store) {stores.add(store);}
    public List<Store> getStores() { return stores; }
    public void setStores(List<Store> stores) { this.stores = stores; }
    
    public List<User> getUsers() { return users; }
    public void setUsers(List<User> users) { this.users = users; }

    public void addUser(User user) {users.add(user);}
    public User getUserById(String userId) {
        for (User user : users) {
            if (user.getId().equals(userId)) {
                return user;
            }
        }
        return null;
    }
}

interface PricingStrategy {
    double computePrice(Reservation reservation);
}

class HourlyPricing implements PricingStrategy {
    public double computePrice(Reservation reservation) {
        long hours = reservation.getDurationInHours();
        return hours * reservation.getVehicle().getPricePerHour();
    }
}
class DailyPricing implements PricingStrategy {
    private static final double DAILY_RATE_MULTIPLIER = 0.8; // 20% discount for daily rate
    
    public double computePrice(Reservation reservation) {
        long hours = reservation.getDurationInHours();
        double days = Math.ceil(hours / 24.0);
        return days * 24 * reservation.getVehicle().getPricePerHour() * DAILY_RATE_MULTIPLIER;
    }
}

class WeeklyPricing implements PricingStrategy {
    private static final double WEEKLY_RATE_MULTIPLIER = 0.7; // 30% discount for weekly rate

    public double computePrice(Reservation reservation) {
        long hours = reservation.getDurationInHours();
        double weeks = Math.ceil(hours / (24.0 * 7));
        return weeks * 24 * 7 * reservation.getVehicle().getPricePerHour() * WEEKLY_RATE_MULTIPLIER;
    }
}


interface Payment {
    boolean processPayment(double amount);
}

class OnlinePayment implements Payment {
    public boolean processPayment(double amount) {
        System.out.println("Processing online payment of $" + amount);
        // Actual payment processing logic would go here
        return true;
    }
}

class CashPayment implements Payment {
    @Override
    public boolean processPayment(double amount) {
        System.out.println("Processing cash payment of $" + amount);
        // Actual payment processing logic would go here
        return true;
    }
}

class PaymentService {
    public PaymentDetails processPayment(Bill bill, PaymentMode paymentMode) {
        Payment payment;

        switch (paymentMode) {
            case ONLINE:
                payment = new OnlinePayment();
                break;
            case CASH:
                payment = new CashPayment();
                break;
            default:
                throw new IllegalArgumentException("Unsupported payment mode: " + paymentMode);
        }

        boolean success = payment.processPayment(bill.computeBillAmount());
        if (success) {
            bill.setPaid(true);
            return new PaymentDetails(new Random().nextInt(100000), "REF" + new Random().nextInt(1000), LocalDateTime.now(), false, paymentMode, bill.getAmount());
        }

        return null;
    }
}


class Bill {
    private int billId;
    private Reservation reservation;
    private double amount;
    private LocalDateTime billingDate;
    private boolean isPaid;
    private PricingStrategy pricingStrategy;

    public Bill(int billId, Reservation reservation, PricingStrategy pricingStrategy) {
        this.billId = billId;
        this.reservation = reservation;
        this.pricingStrategy = pricingStrategy;
        this.billingDate = LocalDateTime.now();
        this.amount = computeBillAmount();
        this.isPaid = false;
    }

    public double computeBillAmount() {
        return pricingStrategy.computePrice(reservation);
    }

    // Getters and setters
    public boolean isPaid() { return isPaid; }
    public void setPaid(boolean paid) { isPaid = paid; }

    public double getAmount() { return amount; }
    public void setAmount(double amount) { this.amount = amount; }
}

public class Main {
    public static void main(String[] args) {
        VehicleRentalSystem rentalSystem = new VehicleRentalSystem();

        // Create locations
        Location location1 = new Location("123 Main St", 12345, "New York", "NY", "USA");
        Location location2 = new Location("456 Oak Ave", 67890, "Los Angeles", "CA", "USA");

        // Create stores ---------------------------------------
        Store store1 = new Store(1, location1);
        Store store2 = new Store(2, location2);

        rentalSystem.addStore(store1);
        rentalSystem.addStore(store2);

        // Create users -----------------------------------------
        User user1 = new User("U1", "John Doe", "john@example.com", "123-456-7890", "password123");
        User user2 = new User("U2", "Jane Smith", "jane@example.com", "098-765-4321", "password456");

        rentalSystem.addUser(user1);
        rentalSystem.addUser(user2);

        // Create vehicles ------------------------------------------
        Car car1 = new Car(1, "Toyota", "Camry", 2022, "ABC123", 25.0, CarType.SEDAN);
        Car car2 = new Car(2, "Honda", "CR-V", 2021, "XYZ789", 30.0, CarType.SUV);
        Car car3 = new Car(3, "Ford", "Mustang", 2023, "SPT123", 50.0, CarType.SPORT);

        // Add vehicles to stores
        store1.addVehicle(car1);
        store1.addVehicle(car2);
        store2.addVehicle(car3);

        // Set up notification services
        EmailNotificationService emailService = new EmailNotificationService();
        SMSNotificationService smsService = new SMSNotificationService();
        PushNotificationService pushService = new PushNotificationService();
        
        store1.addReservationObserver(emailService);
        store1.addReservationObserver(smsService);
        store1.addReservationObserver(pushService);

        // Create a reservation
        LocalDateTime startTime = LocalDateTime.now();
        LocalDateTime endTime = startTime.plusHours(48); // 2 days

        Reservation reservation = new Reservation(
            0, // ID will be set by reservation manager
            user1,
            car1,
            startTime,
            endTime,
            ReservationStatus.INITIATED,
            location1,
            PricingType.DAILY
        );

        // Add reservation through store
        store1.updateOrCreateReservation(reservation, ReservationStatus.SCHEDULED);

        // Create a bill for the reservation
        PricingStrategy pricingStrategy;
        switch (reservation.getPricingType()) {
            case HOURLY:
                pricingStrategy = new HourlyPricing();
                break;
            case DAILY:
                pricingStrategy = new DailyPricing();
                break;
            case WEEKLY:
                pricingStrategy = new WeeklyPricing();
                break;
            default:
                throw new IllegalArgumentException("Unknown pricing type: " + reservation.getPricingType());
        }

        Bill bill = new Bill(1, reservation, pricingStrategy);
        System.out.println("Bill amount: $" + bill.getAmount());

        // Process payment
        PaymentService paymentService = new PaymentService();
        PaymentDetails paymentDetails = paymentService.processPayment(bill, PaymentMode.ONLINE);

        if (paymentDetails != null) {
            System.out.println("Payment processed successfully. Payment ID: " + paymentDetails.getPaymentId());
        }

        // Complete reservation
        store1.getReservationManager().completeReservation(reservation.getReservationId());

        // Demonstrate getting available vehicles
        System.out.println("\nAvailable vehicles at Store 1:");
        for (Vehicle vehicle : store1.getVehicles()) {
            System.out.println("- " + vehicle.toString());
        }
    }
}