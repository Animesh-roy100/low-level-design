#include<iostream>

class Hero {
    // properties
    char name[100];
    int health;
    char level;
};

int main() {
    // std::cout<<"Hello World!";
    // creation of object/instance
    Hero h1;
    std::cout <<"size : " << sizeof(h1) << std::endl;
    return 0;
}