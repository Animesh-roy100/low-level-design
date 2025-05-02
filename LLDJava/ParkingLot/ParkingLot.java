package LLDJava.ParkingLot;

import java.util.ArrayList;
import java.util.List;

public class ParkingLot {
    private List<Level> levels;

    public ParkingLot(int numLevels, int spotsPerLevel) {
        levels = new ArrayList<>();
        for (int i = 0; i < numLevels; i++) {
            levels.add(new Level(i + 1, spotsPerLevel));
        }
    }

    public boolean parkVehicle(Vehicle vehicle) {
        for (Level level : levels) {
            if (level.parkVehicle(vehicle)) {
                return true;
            }
        }
        System.out.println("No available spot for " + vehicle.getLicensePlate());
        return false;
    }

    public boolean removeVehicle(String licensePlate) {
        for (Level level : levels) {
            if (level.removeVehicle(licensePlate)) {
                return true;
            }
        }
        System.out.println("Vehicle " + licensePlate + " not found");
        return false;
    }

    public void displayAvailableSpots() {
        for (Level level : levels) {
            System.out.println("Level " + level.getLevelNumber() + ": " + level.getAvailableSpots() + " spots available");
        }
    }
}
