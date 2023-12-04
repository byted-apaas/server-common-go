// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package http

import (
	"context"

	"github.com/tidwall/gjson"

	exp "github.com/byted-apaas/server-common-go/exceptions"
	"github.com/byted-apaas/server-common-go/utils"
)

func SendLog(ctx context.Context, data interface{}) error {
	body, extra, err := GetFaaSInfraClient(ctx).PostJson(ctx, GetFaaSInfraPathSendLog(), map[string][]string{
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
	return exp.NewErrWithCodeV2(code, msg, utils.GetLogIDFromExtra(extra))
}
