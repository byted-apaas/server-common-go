// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package exceptions

import (
	"fmt"
)

const (
	// Success
	SCFileDownload = ""
	SCSuccess      = "0"

	// User return error
	ErrCodeFunctionError   = "k_cf_ec_100001" // 函数错误
	ErrCodeValidationError = "k_cf_ec_100002" // 校验错误

	// User business error
	ErrCodeDeveloperError = "k_cf_ec_400001783"

	// System error
	ErrCodeInternalError = "k_cf_ec_200001" // 内部系统错误

	// For developer error code
	// 无权限操作流程
	ErrCodeFlowNoPermission = "k_wf_ec_200036"
	// 流程不存在
	ErrCodeFlowNotExist = "k_wf_ec_2001006"
	// 流程实例 ID 不存在
	ErrCodeFlowExecutionNotExist = "k_wf_ec_2001005"
	// 不支持调用此类型流程
	ErrCodeFlowNotSupportCallFlowType = "k_wf_ec_2001004"
	// 不支持撤销该流程 (无人工任务)
	ErrCodeFlowNotSupportRevokeFlow = "k_wf_ec_2001003"
	// 缺少流程的必填参数
	ErrCodeFlowNoReqInputParam = "k_wf_ec_2001002"
	// 流程参数中 APIName 无效
	ErrCodeFlowInvalidParam = "k_wf_ec_2001001"
	// 记录不存在
	ErrCodeDataRecordNotFound = "k_cf_ec_300004"

	// Deprecated
	ErrCodeRateLimitError = "k_cf_ec_200009" // 限流错误
	// Deprecated
	ECOpenAPIRateLimitError = "k_op_ec_20003"
	// Deprecated
	ECFaaSInfraRateLimitError = "k_fs_ec_000004"
	// Deprecated
	ErrCodeInvalidParams = "k_cf_ec_300001" // 请求参数错误
	// Deprecated
	ECUnknownError = "k_ec_000004"
	// Deprecated
	ECOpUnknownError = "k_op_ec_00001"
	// Deprecated
	ECSystemBusy = "k_op_ec_20001"
	// Deprecated
	ECSystemError = "k_op_ec_20002"
	// Deprecated
	FaaSInfraFailCodeMissingToken = "k_fs_ec_100001"
	// Deprecated
	ECTokenExpire = "k_ident_013000"
	// Deprecated
	ECIllegalToken = "k_ident_013001"
	// Deprecated
	ECMissingToken = "k_op_ec_10205"
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
	return fmt.Sprintf("%s [%s]", e.Message, e.Code)
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
		Code:    ErrCodeDeveloperError,
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

// Deprecated
func NewErrWithCode(code, format string, args ...interface{}) *BaseError {
	return &BaseError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		err:     fmt.Errorf(code+", "+format, args...),
		stack:   callers(4, 16),
	}
}

func NewErrWithCodeV2(code, msg, logid string) *BaseError {
	return &BaseError{
		Code:    code,
		Message: msg,
		err:     fmt.Errorf(fmt.Sprintf("%s [%s]（logid=%v）", msg, code, logid)),
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
