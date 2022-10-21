// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package exceptions

import "fmt"

func (e *BaseError) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			fmt.Fprintf(st, "[%s] %s\n%+v", e.Code, e.Message, e.stack)
		default:
			fmt.Fprintf(st, "[%s] %s", e.Code, e.Message)
		}
	case 's':
		fmt.Fprintf(st, "[%s] %s", e.Code, e.Message)
	}
}
