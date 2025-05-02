#include <iostream>
#include <vector>
#include <string>
#include <algorithm>

struct Restaurant {
    std::string name;
    double avgPrice;      // Average price per customer in rupees
    double deliveryTime;  // Estimated delivery time in minutes
    double rating;        // Rating out of 5
    double distance;      // Distance to user in kilometers
};

double calculateRank(const Restaurant& restaurant, double a, double b, double c) {
    // Normalize values: lower price and delivery time are better, higher rating is better
    double normalizedPrice = 1.0 / restaurant.avgPrice;
    double normalizedDeliveryTime = 1.0 / restaurant.deliveryTime;
    double normalizedRating = restaurant.rating / 5.0;

    // Weighted rank formula
    return a * normalizedPrice + b * normalizedDeliveryTime + c * normalizedRating;
}

bool compareRestaurants(const Restaurant& r1, const Restaurant& r2, double a, double b, double c) {
    return calculateRank(r1, a, b, c) > calculateRank(r2, a, b, c);
}

std::vector<Restaurant> recommendRestaurants(const std::vector<Restaurant>& restaurants, double a, double b, double c) {
    std::vector<Restaurant> sortedRestaurants = restaurants;
    std::sort(sortedRestaurants.begin(), sortedRestaurants.end(),
              [&a, &b, &c](const Restaurant& r1, const Restaurant& r2) {
                  return compareRestaurants(r1, r2, a, b, c);
              });
    return sortedRestaurants;
}

int main() {
    // Sample data: {name, avgPrice, deliveryTime, rating, distance}
    std::vector<Restaurant> restaurants = {
        {"Tasty Bites", 15.0, 30.0, 4.5, 2.0},
        {"Quick Eats", 20.0, 25.0, 4.0, 1.5},
        {"Budget Diner", 10.0, 40.0, 3.5, 3.0}
    };

    // Weights: adjust based on priority (sum need not be 1)
    double a = 0.3; // Weight for price
    double b = 0.4; // Weight for delivery time
    double c = 0.3; // Weight for rating

    std::vector<Restaurant> recommended = recommendRestaurants(restaurants, a, b, c);

    std::cout << "Recommended Restaurants:\n";
    for (const auto& restaurant : recommended) {
        std::cout << restaurant.name << " - Rank: " << calculateRank(restaurant, a, b, c)
                  << " (Price: $" << restaurant.avgPrice << ", Delivery: " << restaurant.deliveryTime
                  << " min, Rating: " << restaurant.rating << "/5)\n";
    }

    return 0;
}