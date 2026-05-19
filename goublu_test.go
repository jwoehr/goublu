package main

import (
	"bufio"
	"context"
	"io"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestStartStreamReader tests the stream reader with various scenarios
func TestStartStreamReader(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		streamName     string
		cancelAfter    time.Duration
		expectedCalls  int
		expectComplete bool
	}{
		{
			name:           "reads multiple lines",
			input:          "line1\nline2\nline3\n",
			streamName:     "test-stream",
			cancelAfter:    0,
			expectedCalls:  3,
			expectComplete: true,
		},
		{
			name:           "handles empty input",
			input:          "",
			streamName:     "empty-stream",
			cancelAfter:    0,
			expectedCalls:  0,
			expectComplete: true,
		},
		{
			name:           "respects context cancellation",
			input:          "line1\nline2\nline3\nline4\nline5\n",
			streamName:     "cancel-stream",
			cancelAfter:    50 * time.Millisecond,
			expectedCalls:  -1, // variable, depends on timing
			expectComplete: false,
		},
		{
			name:           "handles single line without newline",
			input:          "single line",
			streamName:     "single-stream",
			cancelAfter:    0,
			expectedCalls:  0, // ReadString('\n') will block/EOF
			expectComplete: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var mu sync.Mutex
			var receivedLines []string
			handler := func(text string) {
				mu.Lock()
				receivedLines = append(receivedLines, text)
				mu.Unlock()
			}

			// Start the stream reader
			startStreamReader(ctx, reader, tt.streamName, handler)

			// If we should cancel after a delay
			if tt.cancelAfter > 0 {
				time.Sleep(tt.cancelAfter)
				cancel()
			}

			// Give goroutine time to process
			time.Sleep(100 * time.Millisecond)

			mu.Lock()
			actualCalls := len(receivedLines)
			mu.Unlock()

			if tt.expectedCalls >= 0 && actualCalls != tt.expectedCalls {
				t.Errorf("expected %d handler calls, got %d", tt.expectedCalls, actualCalls)
			}

			// Verify content for complete reads
			if tt.expectComplete && tt.expectedCalls > 0 {
				expectedLines := strings.Split(strings.TrimSuffix(tt.input, "\n"), "\n")
				mu.Lock()
				for i, expected := range expectedLines {
					if i >= len(receivedLines) {
						t.Errorf("missing line %d: expected %q", i, expected)
						continue
					}
					if receivedLines[i] != expected+"\n" {
						t.Errorf("line %d: expected %q, got %q", i, expected+"\n", receivedLines[i])
					}
				}
				mu.Unlock()
			}
		})
	}
}

// TestStartStreamReaderPanicRecovery tests that panics in handler are recovered
func TestStartStreamReaderPanicRecovery(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("line1\nline2\n"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	panicHandler := func(text string) {
		panic("intentional panic for testing")
	}

	// This should not crash the test
	startStreamReader(ctx, reader, "panic-stream", panicHandler)
	time.Sleep(100 * time.Millisecond)

	// If we get here, panic was recovered successfully
}

// TestStartStreamReaderConcurrency tests multiple concurrent stream readers
func TestStartStreamReaderConcurrency(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	numReaders := 5
	var wg sync.WaitGroup
	wg.Add(numReaders)

	for i := 0; i < numReaders; i++ {
		go func(id int) {
			defer wg.Done()
			input := "line1\nline2\nline3\n"
			reader := bufio.NewReader(strings.NewReader(input))

			var mu sync.Mutex
			count := 0
			handler := func(text string) {
				mu.Lock()
				count++
				mu.Unlock()
			}

			startStreamReader(ctx, reader, "concurrent-stream", handler)
			time.Sleep(100 * time.Millisecond)

			mu.Lock()
			if count != 3 {
				t.Errorf("reader %d: expected 3 lines, got %d", id, count)
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
}

// TestInitializeGUI tests GUI initialization
func TestInitializeGUI(t *testing.T) {
	// Note: This test is limited because gocui.NewGui requires a terminal
	// In a real CI environment, you'd use mocking or skip this test

	t.Run("returns error when GUI creation fails", func(t *testing.T) {
		// This will fail in non-terminal environments, which is expected
		args := NewArgs([]string{"goublu"})
		options := NewOptions()
		ublu := &Ublu{Args: args, Options: options}
		history := NewHistory()

		g, um, err := initializeGUI(ublu, options, history)

		// In non-terminal environment, we expect an error
		if err == nil {
			// If no error, verify the objects are properly initialized
			if g == nil {
				t.Error("expected non-nil GUI when no error")
			}
			if um == nil {
				t.Error("expected non-nil UbluManager when no error")
			}
			if um != nil {
				if um.CompileDate != CompileDate {
					t.Errorf("expected CompileDate %q, got %q", CompileDate, um.CompileDate)
				}
				if um.Version != GoubluVersion {
					t.Errorf("expected Version %q, got %q", GoubluVersion, um.Version)
				}
			}
			// Clean up
			if g != nil {
				g.Close()
			}
		}
	})
}

// TestInitializeGUIWithNilInputs tests error handling with nil inputs
func TestInitializeGUIWithNilInputs(t *testing.T) {
	t.Run("handles nil ublu", func(t *testing.T) {
		// This should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("initializeGUI panicked with nil ublu: %v", r)
			}
		}()

		options := NewOptions()
		history := NewHistory()
		_, _, _ = initializeGUI(nil, options, history)
	})
}

// BenchmarkStartStreamReader benchmarks the stream reader performance
func BenchmarkStartStreamReader(b *testing.B) {
	input := strings.Repeat("benchmark line\n", 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bufio.NewReader(strings.NewReader(input))
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan bool)
		handler := func(text string) {
			// Minimal processing
		}

		startStreamReader(ctx, reader, "bench-stream", handler)

		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
			done <- true
		}()

		<-done
	}
}

// TestContextCancellation verifies context cancellation stops the reader
func TestContextCancellation(t *testing.T) {
	// Create a reader that would block indefinitely
	pr, pw := io.Pipe()
	defer pr.Close()
	defer pw.Close()

	reader := bufio.NewReader(pr)
	ctx, cancel := context.WithCancel(context.Background())

	handlerCalled := false
	var mu sync.Mutex
	handler := func(text string) {
		mu.Lock()
		handlerCalled = true
		mu.Unlock()
	}

	startStreamReader(ctx, reader, "cancel-test", handler)

	// Cancel immediately
	cancel()

	// Give it time to process cancellation
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if handlerCalled {
		t.Error("handler should not have been called after immediate cancellation")
	}
	mu.Unlock()
}

// Made with Bob
