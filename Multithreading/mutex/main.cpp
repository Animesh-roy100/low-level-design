#include<iostream>
#include<thread>
#include<mutex>

std::mutex mtx; // Mutex to protect the critical section

void sharedResourceAccess() {
    mtx.lock();
    // Critical section
    mtx.unlock();
}


int main() {
    std::thread t1(sharedResourceAccess);
    std::thread t2(sharedResourceAccess);
    t1.join();
    t2.join();
    return 0;
}