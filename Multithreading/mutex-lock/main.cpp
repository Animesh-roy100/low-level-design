#include<iostream>
#include<thread>
#include<mutex>

using namespace std;

int counter=0;
mutex mtx; // Mutex to protect the critical section

void increment(int id) {
    for(int i=0; i<100000; i++) {
        lock_guard<mutex> lock(mtx); // Lock the mutex to protect the counter
        ++counter;
    }
    cout << "Thread " << id << " finished" << endl;
}

int main() {
    thread t1(increment, 1);
    thread t2(increment, 2);

    t1.join();
    t2.join();

    cout << "Final counter value: " << counter << endl;
    return 0;
}