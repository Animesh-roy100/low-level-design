package LLDJava.ParkingLot;

public class Main {
    public static void main(String[] args) {
        ParkingLot parkingLot = new ParkingLot(2, 10);
        Vehicle moto1 = new Vehicle("MOTO1", VehicleType.MOTORCYCLE);
        Vehicle car1 = new Vehicle("CAR1", VehicleType.CAR);
        Vehicle truck1 = new Vehicle("TRUCK1", VehicleType.TRUCK);
        parkingLot.parkVehicle(moto1);
        parkingLot.parkVehicle(car1);
        parkingLot.parkVehicle(truck1);
        parkingLot.displayAvailableSpots();
        parkingLot.removeVehicle("CAR1");
        parkingLot.displayAvailableSpots();
    }
}
