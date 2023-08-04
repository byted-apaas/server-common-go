// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package http

import (
	"strings"

	"github.com/byted-apaas/server-common-go/constants"
	"github.com/byted-apaas/server-common-go/utils"
)

const (
	OpenapiPathGetToken  = "/auth/v1/appToken"
	FaaSInfraPathSendLog = "/log/v1/namespaces/:namespace/logs/batchSend"
	InnerAPIGetFunction  = "/cloudfunction/v1/namespaces/:namespace/functions/detail"
)

func GetFaaSInfraPathSendLog() string {
	return strings.ReplaceAll(FaaSInfraPathSendLog, constants.ReplaceNamespace, utils.GetNamespace())
}

func GetInnerAPIPathGetFunction() string {
	return strings.ReplaceAll(InnerAPIGetFunction, constants.ReplaceNamespace, utils.GetNamespace())
}
