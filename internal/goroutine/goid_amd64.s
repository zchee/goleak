// SPDX-FileCopyrightText: Copyright 2020 The goleak Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build go1.11 && !go1.18
// +build go1.11,!go1.18

#include "textflag.h"

// func getg() *g
TEXT ·getg(SB), NOSPLIT, $0-8
	MOVQ (TLS), R14
	MOVQ R14, ret+0(FP)
	RET
