// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package exceptions

import (
	"fmt"
)

const (
	ErrCodeFunctionError   = "k_cf_ec_100001" // 函数错误
	ErrCodeValidationError = "k_cf_ec_100002" // 校验错误

	ErrCodeInternalError  = "k_cf_ec_200001" // 内部系统错误
	ErrCodeRateLimitError = "k_cf_ec_200009" // 限流错误

	ErrCodeInvalidParams  = "k_cf_ec_300001" // 请求参数错误
	ErrCodeDeveloperError = "k_cf_ec_300002" // 开发者错误
)

type BaseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	err     error
	stack   *stack
}

func (e *BaseError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func InternalError(format string, args ...interface{}) *BaseError {
	return &BaseError{
		Code:    ErrCodeInternalError,
		Message: fmt.Sprintf(format, args...),
		err:     fmt.Errorf(format, args...),
		stack:   callers(4, 16),
	}
}

func InvalidParamError(format string, args ...interface{}) *BaseError {
	return &BaseError{
		Code:    ErrCodeInvalidParams,
		Message: fmt.Sprintf(format, args...),
		err:     fmt.Errorf(format, args...),
		stack:   callers(4, 16),
	}
}

func DeveloperError(format string, args ...interface{}) *BaseError {
	return &BaseError{
		Code:    ErrCodeDeveloperError,
		Message: fmt.Sprintf(format, args...),
		err:     fmt.Errorf(format, args...),
		stack:   callers(4, 16),
	}
}

func NewErrWithCode(code, format string, args ...interface{}) *BaseError {
	return &BaseError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		err:     fmt.Errorf(code+", "+format, args...),
		stack:   callers(4, 16),
	}
}

func ErrWrap(err error) *BaseError {
	if err == nil {
		return nil
	}

	baseErr, ok := err.(*BaseError)
	if ok {
		return baseErr
	}

	return &BaseError{
		Code:    ErrCodeInternalError,
		Message: err.Error(),
		err:     fmt.Errorf(err.Error()),
		stack:   callers(4, 16),
	}
}

func ParseErrForUser(err error) *BaseError {
	if err == nil {
		return nil
	}

	baseErr, ok := err.(*BaseError)
	if ok {
		return baseErr
	}

	return &BaseError{
		Code:    ErrCodeDeveloperError,
		Message: err.Error(),
		err:     fmt.Errorf(err.Error()),
		stack:   callers(4, 16),
	}
}
