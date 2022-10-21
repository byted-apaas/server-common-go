// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package http

import (
	"context"

	exp "github.com/byted-apaas/server-common-go/exceptions"
	"github.com/tidwall/gjson"
)

const (
	SuccessCode                     = "0"
	FaaSInfraFailCodeInternalError  = "k_ec_000001"
	FaaSInfraFailCodeTokenExpire    = "k_ident_013000"
	FaaSInfraFailCodeIllegalToken   = "k_ident_013001"
	FaaSInfraFailCodeMissingToken   = "k_fs_ec_100001"
	FaaSInfraFailCodeRateLimitError = "k_fs_ec_000004"
)

func SendLog(ctx context.Context, data interface{}) error {
	body, _, err := GetFaaSInfraClient().PostJson(ctx, GetFaaSInfraPathSendLog(), map[string][]string{
		"Kldx-Version": {"4.0.0"}, // TODO FaaSInfra 后续下掉
	}, data, AppTokenMiddleware, TenantAndUserMiddleware, ServiceIDMiddleware)
	if err != nil {
		return exp.ErrWrap(err)
	}

	code := gjson.GetBytes(body, "code").String()
	msg := gjson.GetBytes(body, "msg").String()
	if code == SuccessCode {
		return nil
	}

	if IsSysError(code) {
		return exp.InternalError("Send log failed, ErrCode: %s, ErrMsg: %s", code, msg)
	}
	return exp.InvalidParamError("Send log failed, ErrCode: %s, ErrMsg: %s", code, msg)
}

func IsSysError(errCode string) bool {
	return errCode == FaaSInfraFailCodeInternalError ||
		errCode == FaaSInfraFailCodeTokenExpire ||
		errCode == FaaSInfraFailCodeIllegalToken ||
		errCode == FaaSInfraFailCodeMissingToken ||
		errCode == FaaSInfraFailCodeRateLimitError
}
