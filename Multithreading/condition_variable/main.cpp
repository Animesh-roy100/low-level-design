#include <iostream>
#include <thread>
#include <queue>
#include <mutex>
#include <condition_variable>

std::queue<int> dataQueue;
std::mutex mtx;
std::condition_variable cv;
bool done = false;  // Indicates if the producer is done

void producer(int count) {
    for (int i = 0; i < count; ++i) {
        std::this_thread::sleep_for(std::chrono::milliseconds(100));  // Simulate work
        std::unique_lock<std::mutex> lock(mtx);
        dataQueue.push(i);
        std::cout << "Produced: " << i << std::endl;
        cv.notify_one();  // Notify one waiting consumer
    }
    
    // Producer finished
    std::unique_lock<std::mutex> lock(mtx);
    done = true;
    cv.notify_all();  // Notify all consumers that production is done
}

void consumer() {
    while (true) {
        std::unique_lock<std::mutex> lock(mtx);
        cv.wait(lock, [] { return !dataQueue.empty() || done; });  // Wait until queue is not empty or producer is done
        
        if (!dataQueue.empty()) {
            int value = dataQueue.front();
            dataQueue.pop();
            lock.unlock();  // Unlock the mutex early to allow other threads to acquire it
            std::cout << "Consumed: " << value << std::endl;
        } else if (done) {
            break;  // Exit if producer has finished
        }
    }
}

int main() {
    std::thread producerThread(producer, 10);
    std::thread consumerThread(consumer);
    
    producerThread.join();
    consumerThread.join();

    return 0;
}
