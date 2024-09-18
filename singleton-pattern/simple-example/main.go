package main

import (
	"fmt"
	"sync"
)

var once sync.Once

type single struct{}

var singleInstance *single

// getInstance provides global access to the Singleton instance
func getInstance() *single {
	// Ensuring that the instance is created only once
	once.Do(
		func() {
			fmt.Println("Creating single instance now.")
			singleInstance = &single{}
		},
	)

	return singleInstance
}

func main() {
	var wg sync.WaitGroup

	// Simulate 30 goroutines accessing the Singleton instance concurrently
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Get the singleton instance
			instance := getInstance()

			fmt.Printf("Instance address: %p\n", instance)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All goroutines have completed.")
}
