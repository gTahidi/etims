package testmanager

import (
	"fmt"
	"time"
)

// ProcessState represents the current state of a business process
type ProcessState string

const (
	StateInitial     ProcessState = "INITIAL"
	StateInProgress  ProcessState = "IN_PROGRESS"
	StateCompleted   ProcessState = "COMPLETED"
	StateFailed      ProcessState = "FAILED"
	StateCompensated ProcessState = "COMPENSATED"
)

// ProcessConfig holds configuration for business processes
type ProcessConfig struct {
	RetryAttempts   int
	RetryDelay      time.Duration
	TimeoutDuration time.Duration
}

// ProcessContext maintains the state during process execution
type ProcessContext struct {
	ProcessID   string
	CurrentStep string
	StartTime   time.Time
	Data        map[string]interface{}
	Errors      []error
	State       ProcessState
	Metadata    map[string]interface{}
}

// NewProcessContext creates a new process context
func NewProcessContext() *ProcessContext {
	return &ProcessContext{
		ProcessID: generateProcessID(),
		StartTime: time.Now(),
		Data:      make(map[string]interface{}),
		Errors:    make([]error, 0),
		State:     StateInitial,
		Metadata:  make(map[string]interface{}),
	}
}

// ProcessStep defines a single step in a business process
type ProcessStep struct {
	Name       string
	Execute    func(ctx *ProcessContext) error
	Validate   func(ctx *ProcessContext) error
	Compensate func(ctx *ProcessContext) error
}

// generateProcessID creates a unique process ID
func generateProcessID() string {
	return fmt.Sprintf("PROC_%d", time.Now().UnixNano())
}
