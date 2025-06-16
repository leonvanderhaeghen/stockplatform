package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

// State represents the circuit breaker state
type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	mu                sync.RWMutex
	state             State
	failureCount      int
	successCount      int
	lastFailureTime   time.Time
	maxFailures       int
	timeout           time.Duration
	maxRequests       int
	onStateChange     func(name string, from State, to State)
	name              string
}

// Config holds circuit breaker configuration
type Config struct {
	Name         string
	MaxFailures  int
	Timeout      time.Duration
	MaxRequests  int
	OnStateChange func(name string, from State, to State)
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config Config) *CircuitBreaker {
	if config.MaxFailures == 0 {
		config.MaxFailures = 5
	}
	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second
	}
	if config.MaxRequests == 0 {
		config.MaxRequests = 1
	}

	return &CircuitBreaker{
		state:         StateClosed,
		maxFailures:   config.MaxFailures,
		timeout:       config.Timeout,
		maxRequests:   config.MaxRequests,
		onStateChange: config.OnStateChange,
		name:          config.Name,
	}
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	if !cb.canExecute() {
		return errors.New("circuit breaker is open")
	}

	err := fn()
	cb.recordResult(err == nil)
	return err
}

// canExecute checks if the circuit breaker allows execution
func (cb *CircuitBreaker) canExecute() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		return time.Since(cb.lastFailureTime) > cb.timeout
	case StateHalfOpen:
		return cb.successCount < cb.maxRequests
	default:
		return false
	}
}

// recordResult records the result of an execution
func (cb *CircuitBreaker) recordResult(success bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if success {
		cb.onSuccess()
	} else {
		cb.onFailure()
	}
}

// onSuccess handles successful execution
func (cb *CircuitBreaker) onSuccess() {
	cb.failureCount = 0

	switch cb.state {
	case StateHalfOpen:
		cb.successCount++
		if cb.successCount >= cb.maxRequests {
			cb.setState(StateClosed)
		}
	}
}

// onFailure handles failed execution
func (cb *CircuitBreaker) onFailure() {
	cb.failureCount++
	cb.lastFailureTime = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failureCount >= cb.maxFailures {
			cb.setState(StateOpen)
		}
	case StateHalfOpen:
		cb.setState(StateOpen)
	}
}

// setState changes the circuit breaker state
func (cb *CircuitBreaker) setState(state State) {
	if cb.state == state {
		return
	}

	prevState := cb.state
	cb.state = state

	switch state {
	case StateClosed:
		cb.failureCount = 0
		cb.successCount = 0
	case StateOpen:
		cb.successCount = 0
	case StateHalfOpen:
		cb.successCount = 0
	}

	if cb.onStateChange != nil {
		cb.onStateChange(cb.name, prevState, state)
	}
}

// State returns the current state
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Counts returns failure and success counts
func (cb *CircuitBreaker) Counts() (failures, successes int) {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failureCount, cb.successCount
}

// String returns string representation of the state
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateHalfOpen:
		return "half-open"
	case StateOpen:
		return "open"
	default:
		return "unknown"
	}
}
