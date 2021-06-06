// SPDX-FileCopyrightText: Copyright 2021 The goleak Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build go1.17 && !go1.18
// +build go1.17,!go1.18

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

	pcsp      uint32
	pcfile    uint32
	pcln      uint32
	npcdata   uint32
	cuOffset  uint32
	funcID    uint8 // funcID
	flag      uint8
	_         [1]byte
	nfuncdata uint8
}
