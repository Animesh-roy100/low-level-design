#include<iostream>
#include<thread>

void threadFunction() {
    std::cout << "Hello from the thread!" << std::endl;
}

int main() {
    std::thread t(threadFunction);
    t.join(); // Wait for the thread to finish
    return 0;
}