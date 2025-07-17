#include<iostream>
#include<thread>
#include<mutex>
#include<condition_variable>
#include<queue>

std::mutex mtx;
std::condition_variable cv;
std::queue<int> buffer;

void producer() {
    // Produce data
    {
        std::lock_guard<std::mutex> lock(mtx);
        buffer.push(42);
    }
    cv.notify_one();  // Notify one waiting consumer
}

void consumer() {
    // Consume data
    int data;
    {
        std::unique_lock<std::mutex> lock(mtx);
        cv.wait(lock, [] { return !buffer.empty(); });
        data = buffer.front();
        buffer.pop();
    }
    // Process data
}

int main() {
    std::thread producerThread(producer);
    std::thread consumerThread(consumer);

    producerThread.join();
    consumerThread.join();

    return 0;
}