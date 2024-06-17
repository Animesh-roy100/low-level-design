#include<iostream>
using namespace std;

class Hero {
    // properties
    private: 
    int health;

    public:
    // char name[100];
    char level;

    void print() {
        std::cout<<level<<std::endl;
    }

    int getHealth() {
        return health;
    }

    int getLevel() {
        return level;
    }

    void setHealth(int h) {
        health = h;
    }

    void setLevel(char ch) {
        level = ch;
    }
};

int main() {
    // std::cout<<"Hello World!";
    // creation of object/instance

    // Static Allocation
    Hero h1;
    h1.setHealth(80);
    h1.setLevel('B');
    cout << "level is " << h1.level << endl;
    cout << "health is " << h1.getHealth() << endl;
    std::cout <<"size : " << sizeof(h1) << std::endl;

    // h1.health = 20;
    // h1.setHealth(20);
    // h1.level = 'A';

    // std::cout<< "health is: " << h1.health << std::endl;
    // std::cout<< "health is: " << h1.getHealth() << std::endl;
    // std::cout<< "level is: " << h1.level << std::endl;

    // Dynamic allocation
    Hero *h2 = new Hero;
    h2->setHealth(70);
    h2->setLevel('A');
    cout << "level is " << (*h2).level << endl;
    cout << "health is " << h2->getHealth() << endl;




    return 0;
}