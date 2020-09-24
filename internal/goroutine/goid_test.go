// Copyright 2020 The goleak Authors.
// SPDX-License-Identifier: BSD-3-Clause

// +build go1.11
// +build !go1.17

package goroutine

import (
	"runtime"
	"sync"
	"testing"
)

func TestInitialGoID(t *testing.T) {
	const max = 10000
	if id := goid(); id < 0 || id > max {
		t.Errorf("got goid = %d, want 0 < goid <= %d", id, max)
	}
}

// TestGoIDSquence verifies that goid returns values which could plausibly be
// goroutine IDs. If this test breaks or becomes flaky, the structs in
// goid_unsafe.go may need to be updated.
func TestGoIDSquence(t *testing.T) {
	// Goroutine IDs are cached by each P.
	runtime.GOMAXPROCS(1)

	// Fill any holes in lower range.
	for i := 0; i < 50; i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			wg.Done()

			// Leak the goroutine to prevent the ID from being
			// reused.
			select {}
		}()
		wg.Wait()
	}

	id := goid()
	for i := 0; i < 100; i++ {
		var (
			newID int64
			wg    sync.WaitGroup
		)
		wg.Add(1)
		go func() {
			newID = goid()
			wg.Done()

			// Leak the goroutine to prevent the ID from being
			// reused.
			select {}
		}()
		wg.Wait()
		if max := id + 100; newID <= id || newID > max {
			t.Errorf("unexpected goroutine ID pattern, got goid = %d, want %d < goid <= %d (previous = %d)", newID, id, max, id)
		}
		id = newID
	}
}

func TestGoFuncNameFromG(t *testing.T) {
	for _, name := range goFuncNames() {
		t.Log(name)
	}
}
