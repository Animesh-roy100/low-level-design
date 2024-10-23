#include<iostream>
#include<unordered_map>
#include<mutex>

using namespace std;

class LRUCache {
private:
    class Node {
        public:
            int key;
            int value;
            Node *prev;
            Node *next;
            Node(int _key, int _value) {
                key =  _key;
                value = _value;
            }
    };

    int cap;
    unordered_map<int, Node*> m;
    Node *head;
    Node *tail;
    mutex mtx;

    void deleteNode(Node *node) {
        Node *prevNode = node->prev;
        Node *nextNode = node->next;
        prevNode->next = nextNode;
        nextNode->prev = prevNode; 
    }

    void addNode(Node *node) {
        Node *temp = head->next;
        node->prev = head;
        node->next = temp;
        head->next = node;
        temp->prev = node;
    }

public:
    LRUCache(int capacity) {
        cap = capacity;
        head = new Node(-1, -1);
        tail = new Node(-1, -1);
        head->next = tail;
        tail->prev = head;
    }

    int get(int key) {
        lock_guard<mutex> lock(mtx);
        if(m.find(key) == m.end()) return -1;

        Node *resNode = m[key];
        int res = resNode->value;
        m.erase(key);
        deleteNode(resNode);
        addNode(resNode);
        m[key] = head->next;

        return res;
    }

    void put(int key, int value) {
        lock_guard<mutex> lock(mtx);
        if(m.find(key) != m.end()) {
            Node *currentNode = m[key];
            m.erase(key);
            deleteNode(currentNode);
        }

        if(m.size() == cap) {
            m.erase(tail->prev->key);
            deleteNode(tail->prev);
        }

        addNode(new Node(key, value));
        m[key] = head->next;
    }
};

int main() {
    LRUCache cache(2);

    cache.put(1, 1);
    cache.put(2, 2);

    cout << cache.get(2) << endl;

    cache.put(3, 3);

    cout<< cache.get(1) << endl;
}


