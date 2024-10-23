#include <iostream>
#include <unordered_map>

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
            key = _key;
            value = _value;
        }
    };
    
    int cap;
    unordered_map<int, Node*> m;
    Node *head;
    Node *tail;
    
    // Function to delete a node from the doubly linked list
    void deleteNode(Node *node) {
        Node *prevNode = node->prev;
        Node *nextNode = node->next;
        prevNode->next = nextNode;
        nextNode->prev = prevNode;
    }
    
    // Function to add a node just after the head
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
        if(m.find(key) != m.end()) {
            Node *existingNode = m[key];
            m.erase(key);
            deleteNode(existingNode);
        }
        
        if(m.size() == cap) {
            m.erase(tail->prev->key);
            deleteNode(tail->prev);
        }
        
        addNode(new Node(key, value));
        m[key] = head->next;
    }
};

// Example function to test the LRUCache
int main() {
    LRUCache cache(2);  // Cache capacity of 2

    // Adding key-value pairs
    cache.put(1, 1);
    cache.put(2, 2);

    cout << "Get 1: " << cache.get(1) << endl;  // Returns 1
    cache.put(3, 3);    // Removes key 2 and adds key 3
    cout << "Get 2: " << cache.get(2) << endl;  // Returns -1 (not found)

    cache.put(4, 4);    // Removes key 1 and adds key 4
    cout << "Get 1: " << cache.get(1) << endl;  // Returns -1 (not found)
    cout << "Get 3: " << cache.get(3) << endl;  // Returns 3
    cout << "Get 4: " << cache.get(4) << endl;  // Returns 4

    return 0;
}
