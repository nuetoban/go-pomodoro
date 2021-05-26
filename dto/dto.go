package dto

import "time"

type State int

const (
	STARTED State = iota + 1
	DONE
	CANCELED
)

type Task struct {
	ID       int
	Name     string
	Interval time.Duration
	State    State

	CreatedAt  time.Time
	FinishedAt time.Time
}
