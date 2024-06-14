#include<iostream>

class Hero {
    // properties
    public:
    char name[100];
    int health;
    char level;

    void print() {
        std::cout<<level<<std::endl;
    }
};

int main() {
    // std::cout<<"Hello World!";
    // creation of object/instance
    Hero h1;
    // std::cout <<"size : " << sizeof(h1) << std::endl;

    h1.health = 20;
    h1.level = 'A';

    std::cout<< "health is: " << h1.health << std::endl;
    std::cout<< "level is: " << h1.level << std::endl;

    return 0;
}