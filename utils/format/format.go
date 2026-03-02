// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package format

import (
	"fmt"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Any retrieves a value from the given map and tries to return a string representation.
func Any(v interface{}) string {
	switch t := v.(type) {
	case bool:
		return strconv.FormatBool(t)

	case string:
		return t

	case *string:
		if t == nil {
			return ""
		}
		return *t

	case []byte:
		return string(t)

	case int:
		return strconv.Itoa(t)

	case int8:
		const base = 10
		return strconv.FormatInt(int64(t), base)

	case int16:
		const base = 10
		return strconv.FormatInt(int64(t), base)

	case int32:
		const base = 10
		return strconv.FormatInt(int64(t), base)

	case int64:
		const base = 10
		return strconv.FormatInt(t, base)

	case float32:
		const fmt = 'f'
		const prec = -1
		const bizSize = 32
		return strconv.FormatFloat(float64(t), fmt, prec, bizSize)

	case float64:
		const fmt = 'f'
		const prec = -1
		const bizSize = 64
		return strconv.FormatFloat(t, fmt, prec, bizSize)

	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%#v", v)
		}

		return string(b)
	}
}
