// SPDX-FileCopyrightText: Copyright 2020 The goleak Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build go1.13 && !go1.18
// +build go1.13,!go1.18

package goroutine

import (
	_ "runtime" // for go:linkname
	_ "unsafe"  // for go:linkname
)

// gobuf is the offsets of sp, pc, and g are known to (hard-coded in) libmach.
type gobuf struct {
	sp   uintptr
	pc   uintptr
	g    uintptr // actual type is guintptr
	ctxt uintptr // actual type is unsafe.Pointer
	ret  uint64  // actual type is sys.Uintreg
	lr   uintptr
	bp   uintptr
}

// stack describes a Go execution stack.
type stack struct {
	lo uintptr
	hi uintptr
}

// ancestorInfo records details of where a goroutine was started.
type ancestorInfo struct {
	pcs  []uintptr
	goid int64
	gopc uintptr
}

// g is the stack parameters.
type g struct {
	stack       stack
	stackguard0 uintptr
	stackguard1 uintptr

	_panic       uintptr // actual type is *_panic.
	_defer       uintptr // actual type is *_defer.
	m            uintptr // actual type is *m.
	sched        gobuf
	syscallsp    uintptr
	syscallpc    uintptr
	stktopsp     uintptr
	param        uintptr // actual type is unsafe.Pointer
	atomicstatus uint32
	stackLock    uint32
	goid         int64
	schedlink    uintptr // actual type is guintptr
	waitsince    int64
	waitreason   uint8 // actual type is waitReason

	preempt       bool
	preemptStop   bool
	preemptShrink bool

	asyncSafePoint bool

	paniconfault bool
	gcscandone   bool
	throwsplit   bool

	activeStackChans bool

	raceignore     int8
	sysblocktraced bool
	sysexitticks   int64
	traceseq       uint64
	tracelastp     uintptr // actual type is puintptr. last P emitted an event for this goroutine
	lockedm        uintptr // actual type is muintptr
	sig            uint32
	writebuf       []byte
	sigcode0       uintptr
	sigcode1       uintptr
	sigpc          uintptr
	gopc           uintptr
	ancestors      *[]ancestorInfo
	startpc        uintptr

	// only use goid and startpc.
	// the fields before it are only listed to calculate the struct offset.
}

// getg returns the pointer to the current g.
//
// The compiler rewrites calls to this function into instructions
// that fetch the g directly from TLS or from the dedicated register.
func getg() *g

// goid returns the ID of the current goroutine.
func goid() int64 {
	return getg().goid
}

// goFuncName returns the function name of the gp goroutine.
func goFuncName(gp *g) string {
	return funcname(findfunc(gp.startpc))
}

//go:linkname allgs runtime.allgs
var allgs allgps

type allgps []*g

func goFuncNames() []string {
	var gps []string
	for _, gp := range allgs {
		gps = append(gps, goFuncName(gp))
	}
	return gps
}
