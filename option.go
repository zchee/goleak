// SPDX-FileCopyrightText: Copyright 2020 The goleak Authors
// SPDX-License-Identifier: BSD-3-Clause

package goleak

import (
	"strings"
	"time"

	"github.com/zchee/goleak/internal/goroutine"
)

// Option lets users specify custom verifications.
type Option interface {
	apply(*opts)
}

// We retry up to 20 times if we can't find the goroutine that
// we are looking for. In between each attempt, we will sleep for
// a short while to let any running goroutines complete.
const _defaultRetries = 20

type opts struct {
	filters    []func(goroutine.Stack) bool
	maxRetries int
	maxSleep   time.Duration
}

// optionFunc lets us easily write options without a custom type.
type optionFunc func(*opts)

func (f optionFunc) apply(opts *opts) { f(opts) }

// IgnoreTopFunction ignores any goroutines where the specified function
// is at the top of the stack. The function name should be fully qualified,
// e.g., go.uber.org/goleak.IgnoreTopFunction
func IgnoreTopFunction(f string) Option {
	return addFilter(func(s goroutine.Stack) bool {
		return s.FirstFunction() == f
	})
}

// IgnoreCurrent records all current goroutines when the option is created, and ignores
// them in any future Find/Verify calls.
func IgnoreCurrent() Option {
	excludeIDSet := map[int64]bool{}
	for _, s := range goroutine.All() {
		excludeIDSet[s.ID()] = true
	}
	return addFilter(func(s goroutine.Stack) bool {
		return excludeIDSet[s.ID()]
	})
}

func maxSleep(d time.Duration) Option {
	return optionFunc(func(opts *opts) {
		opts.maxSleep = d
	})
}

func addFilter(f func(goroutine.Stack) bool) Option {
	return optionFunc(func(opts *opts) {
		opts.filters = append(opts.filters, f)
	})
}

func buildOpts(options ...Option) *opts {
	opts := &opts{
		maxRetries: _defaultRetries,
		maxSleep:   100 * time.Millisecond,
	}
	opts.filters = append(opts.filters,
		isTestStack,
		isSyscallStack,
		isStdLibStack,
		isTraceStack,
	)
	for _, option := range options {
		option.apply(opts)
	}
	return opts
}

func (vo *opts) filter(s goroutine.Stack) bool {
	for _, filter := range vo.filters {
		if filter(s) {
			return true
		}
	}
	return false
}

func (vo *opts) retry(i int) bool {
	if i >= vo.maxRetries {
		return false
	}

	d := time.Duration(int(time.Microsecond) << uint(i))
	if d > vo.maxSleep {
		d = vo.maxSleep
	}
	time.Sleep(d)
	return true
}

// isTestStack is a default filter installed to automatically skip goroutines
// that the testing package runs while the user's tests are running.
func isTestStack(s goroutine.Stack) bool {
	// Until go1.7, the main goroutine ran RunTests, which started
	// the test in a separate goroutine and waited for that test goroutine
	// to end by waiting on a channel.
	// Since go1.7, a separate goroutine is started to wait for signals.
	// T.Parallel is for parallel tests, which are blocked until all serial
	// tests have run with T.Parallel at the top of the stack.
	switch s.FirstFunction() {
	case "testing.RunTests", "testing.(*T).Run", "testing.(*T).Parallel":
		// In pre1.7 and post-1.7, background goroutines started by the testing
		// package are blocked waiting on a channel.
		return strings.HasPrefix(s.State(), "chan receive")
	}
	return false
}

func isSyscallStack(s goroutine.Stack) bool {
	// Typically runs in the background when code uses CGo:
	// https://github.com/golang/go/issues/16714
	return s.FirstFunction() == "runtime.goexit" && strings.HasPrefix(s.State(), "syscall")
}

func isStdLibStack(s goroutine.Stack) bool {
	// Importing os/signal starts a background goroutine.
	// The name of the function at the top has changed between versions.
	if f := s.FirstFunction(); f == "os/signal.signal_recv" || f == "os/signal.loop" {
		return true
	}

	// Using signal.Notify will start a runtime goroutine.
	return strings.Contains(s.Full(), "runtime.ensureSigM")
}

func isTraceStack(s goroutine.Stack) bool {
	if f := s.FirstFunction(); f != "runtime.goparkunlock" {
		return false
	}

	return strings.Contains(s.Full(), "runtime.ReadTrace")
}
