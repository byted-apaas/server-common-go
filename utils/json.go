// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package utils

import (
	rawJson "encoding/json"
	"errors"
	"fmt"

	"github.com/json-iterator/go"
)

var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonUseNumber = jsoniter.Config{EscapeHTML: true, SortMapKeys: true, ValidateJsonRawMessage: true, UseNumber: true}.Froze()
)

func JsonMarshalBytes(val interface{}) ([]byte, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func JsonUnmarshalBytes(val []byte, addr interface{}) error {
	if err := jsonUseNumber.Unmarshal(val, addr); err != nil {
		e := rawJson.Unmarshal(val, &addr)
		return errors.New(fmt.Sprintf("JsonUnmarshalBytes failed, err: %v", e))
	}
	return nil
}

func Decode(input interface{}, output interface{}) error {
	bytes, err := JsonMarshalBytes(input)
	if err != nil {
		return err
	}

	if e := JsonUnmarshalBytes(bytes, output); e != nil {
		return e
	}
	return nil
}
