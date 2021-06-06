// SPDX-FileCopyrightText: Copyright 2020 The goleak Authors
// SPDX-License-Identifier: BSD-3-Clause

package goroutine

import (
	"bytes"
	"fmt"
)

// Stack represents a single Goroutine's stack.
type Stack struct {
	id            int64
	state         string
	firstFunction string
	fullStack     *bytes.Buffer
}

// ID returns the goroutine ID.
func (s Stack) ID() int64 {
	return s.id
}

// State returns the Goroutine's state.
func (s Stack) State() string {
	return s.state
}

// Full returns the full stack trace for this goroutine.
func (s Stack) Full() string {
	return s.fullStack.String()
}

// FirstFunction returns the name of the first function on the stack.
func (s Stack) FirstFunction() string {
	return s.firstFunction
}

func (s Stack) String() string {
	return fmt.Sprintf(
		"Goroutine %v in state %v, with %v on top of the stack:\n%s",
		s.id, s.state, s.firstFunction, s.Full())
}

// Current returns the stack for the current goroutine.
func Current() Stack {
	gp := getg()
	return Stack{
		id:            gp.goid,
		firstFunction: goFuncName(gp),
	}

}

// All returns the stacks for all running goroutines.
func All() (ss []Stack) {
	for _, gp := range allgs {
		ss = append(ss, Stack{
			id:            gp.goid,
			firstFunction: goFuncName(gp),
		})
	}
	return ss
}
