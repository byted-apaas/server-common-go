// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package utils

import (
	"context"

	"github.com/byted-apaas/server-common-go/constants"
	exp "github.com/byted-apaas/server-common-go/exceptions"
	"github.com/byted-apaas/server-common-go/structs"
)

func SetTenantToCtx(ctx context.Context, tenant *structs.Tenant) context.Context {
	return context.WithValue(ctx, constants.CtxKeyTenant, tenant)
}

func GetTenantFromCtx(ctx context.Context) (*structs.Tenant, error) {
	cast, ok := ctx.Value(constants.CtxKeyTenant).(*structs.Tenant)
	if !ok {
		return nil, exp.InternalError("Can not find tenant from ctx.")
	}

	return cast, nil
}

func GetTenantIDFromCtx(ctx context.Context) int64 {
	cast, ok := ctx.Value(constants.CtxKeyTenant).(*structs.Tenant)
	if !ok {
		return 0
	}

	return cast.ID
}

func GetNamespaceFromCtx(ctx context.Context) string {
	cast, ok := ctx.Value(constants.CtxKeyTenant).(*structs.Tenant)
	if !ok {
		return ""
	}

	return cast.Namespace
}

func SetUserIDToCtx(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, constants.CtxKeyUser, userID)
}

func GetUserIDFromCtx(ctx context.Context) int64 {
	cast, ok := ctx.Value(constants.CtxKeyUser).(int64)
	if !ok {
		return -1
	}
	return cast
}

func SetTTEnvToCtx(ctx context.Context, ttEnv string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyTTEnv, ttEnv)
}

func GetTTEnvFromCtx(ctx context.Context) string {
	cast, _ := ctx.Value(constants.CtxKeyTTEnv).(string)

	return cast
}

func SetKEnvToCtx(ctx context.Context, kEnv string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyKEnv, kEnv)
}

func SetLogIDToCtx(ctx context.Context, logID string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyLogID, logID)
}

func GetLogIDFromCtx(ctx context.Context) string {
	cast, _ := ctx.Value(constants.CtxKeyLogID).(string)

	return cast
}

func SetSourceTypeToCtx(ctx context.Context, sourceType int) context.Context {
	return context.WithValue(ctx, constants.CtxKeySourceType, sourceType)
}

func GetSourceTypeFromCtx(ctx context.Context) int {
	cast, _ := ctx.Value(constants.CtxKeySourceType).(int)

	return cast
}

func SetTriggerTypeToCtx(ctx context.Context, triggerType string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyTriggerType, triggerType)
}

func GetTriggerTypeFromCtx(ctx context.Context) string {
	cast, _ := ctx.Value(constants.CtxKeyTriggerType).(string)

	return cast
}

func SetDebugTypeToCtx(ctx context.Context, debugType int) context.Context {
	return context.WithValue(ctx, constants.CtxKeyDebugType, debugType)
}

func GetDebugTypeFromCtx(ctx context.Context) int {
	cast, _ := ctx.Value(constants.CtxKeyDebugType).(int)

	return cast
}

func SetFunctionNameToCtx(ctx context.Context, functionName string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyFunctionName, functionName)
}

func GetFunctionNameFromCtx(ctx context.Context) string {
	cast, _ := ctx.Value(constants.CtxKeyFunctionName).(string)

	return cast
}

func SetTriggerTaskIDToCtx(ctx context.Context, taskID int64) context.Context {
	return context.WithValue(ctx, constants.CtxKeyTriggerTaskID, taskID)
}

func GetTriggerTaskIDFromCtx(ctx context.Context) int64 {
	cast, _ := ctx.Value(constants.CtxKeyTriggerTaskID).(int64)
	return cast
}

func SetDistributedMaskToCtx(ctx context.Context, mask string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyDistributedMask, mask)
}

func GetDistributedMaskFromCtx(ctx context.Context) string {
	cast, _ := ctx.Value(constants.CtxKeyDistributedMask).(string)
	return cast
}

func SetLoopMaskToCtx(ctx context.Context, mask []string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyLoopMask, mask)
}

func GetLoopMaskFromCtx(ctx context.Context) []string {
	cast, _ := ctx.Value(constants.CtxKeyLoopMask).([]string)
	return cast
}

func SetApiTimeoutToCtx(ctx context.Context, timeout map[string]int64) context.Context {
	return context.WithValue(ctx, constants.CtxKeyAPITimeoutMap, timeout)
}

func GetApiTimeoutFromCtx(ctx context.Context) map[string]int64 {
	cast, _ := ctx.Value(constants.CtxKeyAPITimeoutMap).(map[string]int64)
	return cast
}

func SetInnerAPIPSMToCtx(ctx context.Context, psm string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyInnerAPIPSM, psm)
}

func GetInnerAPIPSMFromCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	domain, _ := ctx.Value(constants.CtxKeyInnerAPIPSM).(string)
	return domain
}

func SetEnvBoeToCtx(ctx context.Context, boe string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyEnvBoe, boe)
}

func GetEnvBoeFromCtx(ctx context.Context) string {
	cast, _ := ctx.Value(constants.CtxKeyEnvBoe).(string)
	return cast
}

func SetApiTimeoutMethodToCtx(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyAPITimeoutMethod, method)
}

func SetUserContext(ctx context.Context, userCtx structs.UserContext) context.Context {
	return context.WithValue(ctx, constants.CtxUserContext, userCtx)
}

func GetUserContext(ctx context.Context) (res structs.UserContext) {
	if ctx == nil {
		return
	}
	res, _ = ctx.Value(constants.CtxUserContext).(structs.UserContext)
	return res
}
