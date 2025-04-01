package schedular

import (
	"container/heap"
	"fmt"
	"task-schedular/taskqueue"
	"time"
)

type Schedular interface {
	AddTask(task *taskqueue.Task)
	Start()
	Stop()
	GetTaskStatus(id string) string
}

type TaskSchedular struct {
	queue     taskqueue.TaskQueue
	stopChan  chan struct{}
	running   bool
	taskMap   map[string]*taskqueue.Task // for status tracking
	doneChans map[string]chan bool       // signal task completion
}

func NewTaskScheduler() *TaskSchedular {
	return &TaskSchedular{
		queue:     make(taskqueue.TaskQueue, 0),
		stopChan:  make(chan struct{}),
		running:   false,
		taskMap:   make(map[string]*taskqueue.Task),
		doneChans: make(map[string]chan bool),
	}
}

func (s *TaskSchedular) AddTask(task *taskqueue.Task) {
	task.Status = "pending"
	s.taskMap[task.ID] = task
	s.doneChans[task.ID] = make(chan bool)
	heap.Push(&s.queue, task)
}

func (s *TaskSchedular) Start() {
	if s.running {
		return
	}

	s.running = true

	go func() {
		for s.queue.Len() > 0 {
			nextTask := s.queue[0]
			delay := time.Until(nextTask.ExecutedAt)

			select {
			case <-time.After(delay):
				task := heap.Pop(&s.queue).(*taskqueue.Task)
				s.executeTask(task)
			case <-s.stopChan:
				fmt.Println("Scheduler stopped")
				s.running = false
				return
			}
		}
		s.running = false
	}()
}

func (s *TaskSchedular) executeTask(task *taskqueue.Task) {
	task.Status = "running"
	fmt.Printf("Executing task: %s (Priority: %d) at %v\n", task.ID, task.Priority, time.Now())

	go func(t *taskqueue.Task) {
		attempts := 0
		for attempts <= t.MaxRetries {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Task %s panicked: %v\n", t.ID, r)
				}
			}()

			t.Job()
			t.Status = "completed"
			s.doneChans[t.ID] <- true
			break // Success, exit retry loop
		}

		// Handle recurrence
		if t.Interval > 0 && t.Status == "completed" {
			t.ExecutedAt = time.Now().Add(t.Interval)
			t.Status = "pending"
			heap.Push(&s.queue, t)
		} else if t.Status != "completed" {
			fmt.Printf("Task %s failed after %d retries\n", t.ID, t.MaxRetries)
			s.doneChans[t.ID] <- false
		}
	}(task)
}

func (s *TaskSchedular) Stop() {
	if s.running {
		close(s.stopChan)
	}
}

func (s *TaskSchedular) GetTaskStatus(id string) string {
	if task, exists := s.taskMap[id]; exists {
		return task.Status
	}

	return "unknown"
}
