// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package http

import (
	"context"
	"github.com/byted-apaas/server-common-go/utils"

	exp "github.com/byted-apaas/server-common-go/exceptions"
	"github.com/tidwall/gjson"
)

func SendLog(ctx context.Context, data interface{}) error {
	body, extra, err := GetFaaSInfraClient().PostJson(ctx, GetFaaSInfraPathSendLog(), map[string][]string{
		"Kldx-Version": {"4.0.0"}, // TODO FaaSInfra 后续下掉
	}, data, AppTokenMiddleware, TenantAndUserMiddleware, ServiceIDMiddleware)
	if err != nil {
		return exp.ErrWrap(err)
	}

	code := gjson.GetBytes(body, "code").String()
	msg := gjson.GetBytes(body, "msg").String()
	if code == exp.SCSuccess {
		return nil
	}
	return exp.NewErrWithCode(code, msg, utils.GetLogIDFromExtra(extra))
}
