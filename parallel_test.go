package parallel

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"golang.org/x/net/context"
)

var errExpected error = errors.New("no good, sorry")

func TestSimple(t *testing.T) {
	err := Parallel(nil, nil, sleep, sleep)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWithChannelCancel(t *testing.T) {
	done := make(chan struct{})
	cancel := func() { close(done) }

	err := Parallel(done, cancel,
		func() error {
			time.Sleep(100 * time.Millisecond)
			cancel()
			time.Sleep(100 * time.Millisecond)
			return errors.New("shouldn't get here")
		},
	)

	if err != Canceled {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWithContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	err := Parallel(ctx.Done(), cancel,
		func() error {
			time.Sleep(100 * time.Millisecond)
			cancel()
			time.Sleep(100 * time.Millisecond)
			return errors.New("shouldn't get here")
		},
	)

	if err != Canceled {
		t.Fatalf("unexpected error: %v", err)
	}
}

func snooze(duration time.Duration) error {
	fmt.Printf("snoozing..\n")
	time.Sleep(duration)
	fmt.Printf("..snoozed!\n")
	return nil
}

func sleep() error {
	fmt.Printf("sleeping..\n")
	duration := time.Duration(rand.Int63n(1000)) * time.Millisecond
	time.Sleep(duration)
	fmt.Printf("..slept!\n")
	return nil
}

func bad() error {
	time.Sleep(800 * time.Millisecond)
	return errExpected
}
