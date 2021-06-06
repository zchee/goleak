// SPDX-FileCopyrightText: Copyright 2020 The goleak Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build go1.11 && !go1.18
// +build go1.11,!go1.18

#include "textflag.h"

// func getg() *g
TEXT ·getg(SB), NOSPLIT, $0-8
	MOVD g, R0         // g
	MOVD R0, ret+0(FP)
	RET
