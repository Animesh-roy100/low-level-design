package LLDJava.ParkingLot;

import java.util.ArrayList;
import java.util.List;

public class Level {   
    private int levelNumber;
    private List<ParkingSpot> spots;

    public Level(int levelNumber, int numSpots) {
        this.levelNumber = levelNumber;
        this.spots = new ArrayList<>();
        int third = numSpots / 3;
        for (int i = 0; i < numSpots; i++) {
            SpotType type;
            if (i < third) {
                type = SpotType.MOTORCYCLE;
            } else if (i < 2 * third) {
                type = SpotType.COMPACT;
            } else {
                type = SpotType.LARGE;
            }
            spots.add(new ParkingSpot(i + 1, type));
        }
    }

    public int getLevelNumber() {
        return levelNumber;
    }

    public boolean parkVehicle(Vehicle vehicle) {
        for (ParkingSpot spot : spots) {
            if (!spot.isOccupied() && spot.canPark(vehicle.getType())) {
                spot.park(vehicle);
                System.out.println("Parked " + vehicle.getLicensePlate() + " at Level " + levelNumber + ", Spot " + spot.getSpotNumber());
                return true;
            }
        }
        return false;
    }

    public boolean removeVehicle(String licensePlate) {
        for (ParkingSpot spot : spots) {
            if (spot.isOccupied() && spot.getVehicle() != null && spot.getVehicle().getLicensePlate().equals(licensePlate)) {
                spot.removeVehicle();
                System.out.println("Removed " + licensePlate + " from Level " + levelNumber + ", Spot " + spot.getSpotNumber());
                return true;
            }
        }
        return false;
    }

    public int getAvailableSpots() {
        int count = 0;
        for (ParkingSpot spot : spots) {
            if (!spot.isOccupied()) {
                count++;
            }
        }
        return count;
    }
}
