package main

import (
	"fmt"
	"task-schedular/schedular"
	"task-schedular/taskqueue"
	"time"
)

func main() {
	fmt.Println("Task Schedular")

	// Create a scheduler
	scheduler := schedular.NewTaskScheduler()

	// Define some tasks
	task1 := &taskqueue.Task{
		ID:         "1",
		ExecutedAt: time.Now().Add(2 * time.Second),
		Interval:   0, // One-time task
		Priority:   1,
		Job: func() {
			fmt.Println("Task 1: High priority, one-time")
		},
		MaxRetries: 2,
		Status:     "pending",
	}

	task2 := &taskqueue.Task{
		ID:         "2",
		ExecutedAt: time.Now().Add(1 * time.Second),
		Interval:   3 * time.Second, // Recurring every 3s
		Priority:   0,               // Highest priority
		Job: func() {
			fmt.Println("Task 2: Highest priority, recurring")
		},
		MaxRetries: 1,
		Status:     "pending",
	}

	task3 := &taskqueue.Task{
		ID:         "3",
		ExecutedAt: time.Now().Add(3 * time.Second),
		Interval:   0,
		Priority:   2,
		Job: func() {
			fmt.Println("Task 3: Low priority, one-time")
		},
		MaxRetries: 0,
		Status:     "pending",
	}

	// Add tasks to the scheduler
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)

	// Start the scheduler
	scheduler.Start()

	// Monitor status and wait for some tasks
	time.Sleep(5 * time.Second)
	fmt.Printf("Task 1 status: %s\n", scheduler.GetTaskStatus("1"))
	fmt.Printf("Task 2 status: %s\n", scheduler.GetTaskStatus("2"))
	fmt.Printf("Task 3 status: %s\n", scheduler.GetTaskStatus("3"))

	// Wait a bit longer to see recurrence
	time.Sleep(5 * time.Second)
	scheduler.Stop()
	fmt.Println("Scheduler stopped manually")
}
