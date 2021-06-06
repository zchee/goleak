// SPDX-FileCopyrightText: Copyright 2020 The goleak Authors
// SPDX-License-Identifier: BSD-3-Clause

package goleak

import (
	"fmt"
	"testing"

	"github.com/zchee/goleak/internal/goroutine"
)

func filterStacks(stacks []goroutine.Stack, skipID int64, opts *opts) []goroutine.Stack {
	filtered := stacks[:0]

	for _, stack := range stacks {
		// Always skip the running goroutine.
		if stack.ID() == skipID {
			continue
		}
		// Run any default or user-specified filters.
		if opts.filter(stack) {
			continue
		}
		filtered = append(filtered, stack)
	}

	return filtered
}

// Find looks for extra goroutines, and returns a descriptive error if
// any are found.
func Find(options ...Option) error {
	cur := goroutine.Current().ID()

	opts := buildOpts(options...)
	var stacks []goroutine.Stack
	retry := true
	for i := 0; retry; i++ {
		stacks = filterStacks(goroutine.All(), cur, opts)

		if len(stacks) == 0 {
			return nil
		}
		retry = opts.retry(i)
	}

	return fmt.Errorf("found unexpected goroutines:\n%s", stacks)
}

// VerifyNone marks the given TestingT as failed if any extra goroutines are
// found by Find. This is a helper method to make it easier to integrate in
// tests by doing:
// 	defer VerifyNone(t)
func VerifyNone(tb testing.TB, options ...Option) {
	if err := Find(options...); err != nil {
		tb.Error(err)
	}
}
