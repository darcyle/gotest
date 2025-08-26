package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Retry executes a function with exponential backoff and jitter
// This is a generic function that demonstrates Go generics usage
func Retry[T any](ctx context.Context, attempts int, base time.Duration, fn func(context.Context) (T, error)) (T, error) {
	// TODO: Implement retry logic with exponential backoff
	// TODO: Add jitter to prevent thundering herd
	// TODO: Respect context cancellation/timeout
	// TODO: Wrap errors with attempt information
	
	var zero T
	
	if attempts <= 0 {
		return zero, fmt.Errorf("attempts must be > 0, got %d", attempts)
	}
	
	var lastErr error
	
	for attempt := 1; attempt <= attempts; attempt++ {
		// Check context before attempting
		select {
		case <-ctx.Done():
			return zero, fmt.Errorf("context cancelled after %d attempts: %w", attempt-1, ctx.Err())
		default:
		}
		
		result, err := fn(ctx)
		if err == nil {
			return result, nil
		}
		
		lastErr = err
		
		// Don't sleep after the last attempt
		if attempt == attempts {
			break
		}
		
		// Don't retry on context cancellation
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			break
		}
		
		// TODO: Calculate backoff duration with jitter
		backoff := calculateBackoffWithJitter(base, attempt)
		
		// Sleep with context cancellation check
		timer := time.NewTimer(backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return zero, fmt.Errorf("context cancelled during backoff after attempt %d: %w", attempt, ctx.Err())
		case <-timer.C:
			// Continue to next attempt
		}
	}
	
	return zero, fmt.Errorf("all %d attempts failed, last error: %w", attempts, lastErr)
}

// calculateBackoffWithJitter implements exponential backoff with jitter
func calculateBackoffWithJitter(base time.Duration, attempt int) time.Duration {
	// TODO: Exponential backoff: base * 2^(attempt-1)
	// TODO: Add jitter to prevent thundering herd
	
	// Exponential backoff
	exponential := base * time.Duration(1<<uint(attempt-1))
	
	// Add jitter (±25%)
	jitter := time.Duration(rand.Int63n(int64(exponential/2))) // 0 to 50%
	jitter = jitter - time.Duration(rand.Int63n(int64(exponential/4))) // -25% to +25%
	
	result := exponential + jitter
	
	// Cap at reasonable maximum (e.g., 30 seconds)
	maxBackoff := 30 * time.Second
	if result > maxBackoff {
		result = maxBackoff
	}
	
	return result
}

// RetryableError wraps an error to indicate it should be retried
type RetryableError struct {
	Err error
}

func (e RetryableError) Error() string {
	return e.Err.Error()
}

func (e RetryableError) Unwrap() error {
	return e.Err
}

// IsRetryable checks if an error should be retried
func IsRetryable(err error) bool {
	// TODO: Define which errors are retryable
	// Database connection errors, timeouts, etc. should be retryable
	// Application logic errors should not be retryable
	
	var retryableErr RetryableError
	if errors.As(err, &retryableErr) {
		return true
	}
	
	// Context errors are not retryable
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	
	// Add more specific error type checks here
	// For example, database connection errors, network timeouts, etc.
	
	return false
}

// Example usage of Retry function
func ExampleRetryUsage(ctx context.Context, store PurchaseStore) error {
	// TODO: Example of using Retry with database operations
	
	result, err := Retry(ctx, 3, 100*time.Millisecond, func(ctx context.Context) ([]Purchase, error) {
		return store.ClaimBatchForEnrichment(ctx, 10)
	})
	
	if err != nil {
		return fmt.Errorf("failed to claim purchases after retries: %w", err)
	}
	
	fmt.Printf("Successfully claimed %d purchases\n", len(result))
	return nil
}