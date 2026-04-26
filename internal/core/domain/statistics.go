package domain

import "time"

type Statistics struct {
	TasksCreated               int
	TasksCompleted             int
	TasksCompletedRate         *float64
	TasksAverageCompletionTime *time.Duration
}

func NewStatistics(
	taskCreated int,
	tasksCompleted int,
	tasksCompletedRate *float64,
	tasksAverageCompletionTime *time.Duration,
) Statistics {
	return Statistics{
		TasksCreated:               taskCreated,
		TasksCompleted:             taskCreated,
		TasksCompletedRate:         tasksCompletedRate,
		TasksAverageCompletionTime: tasksAverageCompletionTime,
	}
}
