// SPDX-FileCopyrightText: Copyright 2020 The goleak Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build go1.16 && !go1.17
// +build go1.16,!go1.17

package goroutine

// _func is layout of in-memory per-function information prepared by linker
// See https://golang.org/s/go12symtab.
//
// Keep in sync with linker (../cmd/link/internal/ld/pcln.go:/pclntab)
// and with package debug/gosym and with symtab.go in package runtime.
type _func struct {
	entry   uintptr
	nameoff int32

	args        int32
	deferreturn uint32

	pcsp      int32
	pcfile    int32
	pcln      int32
	npcdata   int32
	funcID    uint8 // funcID
	nfuncdata uint8
}
