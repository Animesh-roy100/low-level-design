package scheduler

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
	taskChan  chan *taskqueue.Task       // channel for worker pool
	workers   int
}

func NewTaskScheduler(numWorkers int) *TaskSchedular {
	s := &TaskSchedular{
		queue:     make(taskqueue.TaskQueue, 0),
		stopChan:  make(chan struct{}),
		running:   false,
		taskMap:   make(map[string]*taskqueue.Task),
		doneChans: make(map[string]chan bool),
		taskChan:  make(chan *taskqueue.Task, numWorkers),
		workers:   numWorkers,
	}
	// start workers for task execution
	s.startWorkers()
	return s
}

func (s *TaskSchedular) startWorkers() {
	for workerID := range make([]int, s.workers) {
		go func() {
			fmt.Printf("Worker %d started\n", workerID)
			for {
				select {
				case task := <-s.taskChan:
					s.executeTask(task)
				case <-s.stopChan:
					fmt.Printf("Worker %d stopped\n", workerID)
					return
				}
			}
		}()
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
				task.Status = "running"
				fmt.Printf("Task %s assigned to worker pool (Priority: %d) at %v\n", task.ID, task.Priority, time.Now())
				s.taskChan <- task
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
	attempts := 0
	for attempts <= task.MaxRetries {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Task %s panicked: %v\n", task.ID, r)
			}
		}()

		task.Job()
		task.Status = "completed"
		s.doneChans[task.ID] <- true
		break // Success
	}

	if task.Interval > 0 && task.Status == "completed" {
		task.ExecutedAt = time.Now().Add(task.Interval)
		task.Status = "pending"
		heap.Push(&s.queue, task)
	} else if task.Status != "completed" {
		fmt.Printf("Task %s failed after %d retries\n", task.ID, task.MaxRetries)
		s.doneChans[task.ID] <- false
	}
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
