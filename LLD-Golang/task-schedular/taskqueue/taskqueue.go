package taskqueue

import "time"

type Task struct {
	ID         string
	ExecutedAt time.Time
	Interval   time.Duration // for recurring tasks
	Priority   int           // 0 = highest priority
	Job        func()
	MaxRetries int
	Status     string // PENDING, RUNNING, COMPLETED
}

// priority queue of tasks
type TaskQueue []*Task

func (tq TaskQueue) Len() int {
	return len(tq)
}

func (tq TaskQueue) Less(i, j int) bool {
	// Low priority value means high priority
	if tq[i].Priority != tq[j].Priority {
		return tq[i].Priority < tq[j].Priority
	}

	// for same priority check execution time
	return tq[i].ExecutedAt.Before(tq[j].ExecutedAt)
}

func (tq TaskQueue) Swap(i, j int) {
	tq[i], tq[j] = tq[j], tq[i]
}

func (tq *TaskQueue) Push(x any) {
	*tq = append(*tq, x.(*Task))
}

func (tq *TaskQueue) Pop() any {
	old := *tq
	n := len(old)
	task := old[n-1]
	old[n-1] = nil // avoid memory leak
	*tq = old[0 : n-1]
	return task
}
