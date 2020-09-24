// Copyright 2020 The goleak Authors.
// SPDX-License-Identifier: BSD-3-Clause

package goroutine

import (
	_ "runtime" // for go:linkname
	_ "unsafe"  // for go:linkname
)

type funcInfo struct {
	*_func
	datap uintptr // *moduledata
}

//go:linkname findfunc runtime.findfunc
func findfunc(pc uintptr) funcInfo

//go:linkname funcname runtime.funcname
func funcname(f funcInfo) string
