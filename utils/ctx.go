// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package utils

import (
	"context"
	"fmt"
	"strings"

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

func SetAppInfoToCtx(ctx context.Context, appInfo *structs.AppInfo) context.Context {
	return context.WithValue(ctx, constants.CtxKeyApp, appInfo)
}

func GetAppInfoFromCtx(ctx context.Context) (*structs.AppInfo, error) {
	cast, ok := ctx.Value(constants.CtxKeyApp).(*structs.AppInfo)
	if !ok {
		return nil, exp.InternalError("Can not find appInfo from ctx.")
	}

	return cast, nil
}

func SetEventInfoToCtx(ctx context.Context, appInfo *structs.EventInfo) context.Context {
	return context.WithValue(ctx, constants.CtxKeyEvent, appInfo)
}

func GetEventInfoFromCtx(ctx context.Context) (*structs.EventInfo, error) {
	cast, ok := ctx.Value(constants.CtxKeyEvent).(*structs.EventInfo)
	if !ok {
		return nil, exp.InternalError("Can not find eventInfo from ctx.")
	}

	return cast, nil
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

func SetAPaaSPersistFaaSMapToCtx(ctx context.Context, aPaaSPersistFaaSMap map[string]string) context.Context {
	return context.WithValue(ctx, constants.PersistFaaSKeySummarized, aPaaSPersistFaaSMap)
}

func GetAPaaSPersistFaaSMapFromCtx(ctx context.Context) (res map[string]string) {
	defer func() {
		if res == nil {
			res = make(map[string]string)
		}
	}()
	if ctx == nil {
		return
	}
	res, _ = ctx.Value(constants.PersistFaaSKeySummarized).(map[string]string)
	return res
}

func WithAPaaSPersistFaaSValue(ctx context.Context, key, value string) context.Context {
	m := GetAPaaSPersistFaaSMapFromCtx(ctx)
	key = strings.TrimPrefix(key, constants.APAAS_PERSIST_FAAS_PREFIX)
	m[constants.APAAS_PERSIST_FAAS_PREFIX+key] = value
	return SetAPaaSPersistFaaSMapToCtx(ctx, m)
}

func GetAPaaSPersistFaaSValueFromCtx(ctx context.Context, key string) string {
	m := GetAPaaSPersistFaaSMapFromCtx(ctx)
	if value, ok := m[key]; ok {
		return value
	}
	return ""
}

func GetAPaaSPersistFaaSMapStr(ctx context.Context) string {
	m := GetAPaaSPersistFaaSMapFromCtx(ctx)
	res, _ := json.Marshal(m)
	return string(res)
}

func SetAPaaSLaneToCtx(ctx context.Context, lane string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyAPaaSLane, lane)
}

func GetAPaaSLaneFromCtx(ctx context.Context) string {
	cast, _ := ctx.Value(constants.CtxKeyAPaaSLane).(string)
	return cast
}

func SetUserAndAuthTypeToCtx(ctx context.Context, authType *string) context.Context {
	return SetUserAndMixAuthTypeToCtx(ctx, authType, false)
}

func SetUserAndMixAuthTypeToCtx(ctx context.Context, authType *string, isMix bool) context.Context {
	userID := GetUserIDFromCtx(ctx)
	ctx = context.WithValue(ctx, constants.HttpHeaderKeyUser, fmt.Sprintf("%d", userID))
	if authType != nil {
		if userID == -1 || *authType == constants.AuthTypeSystem {
			ctx = context.WithValue(ctx, constants.AuthTypeKey, constants.AuthTypeSystem)
		} else if *authType == constants.AuthTypeUser {
			if isMix {
				ctx = context.WithValue(ctx, constants.AuthTypeKey, constants.AuthTypeMixUserSystem)
			} else {
				ctx = context.WithValue(ctx, constants.AuthTypeKey, constants.AuthTypeUser)
			}
		}
	}
	return ctx
}

func SetUserAndMixAuthTypeToHeaders(ctx context.Context, headers map[string][]string, isMix bool) map[string][]string {
	if headers == nil {
		headers = make(map[string][]string)
	}

	userID := GetUserIDFromCtx(ctx)
	headers[constants.HttpHeaderKeyUser] = []string{fmt.Sprintf("%d", userID)}
	if authType, ok := ctx.Value(constants.AuthTypeKey).(string); ok {
		if userID == -1 || authType == constants.AuthTypeSystem {
			headers[constants.HTTPHeaderKeyAuthType] = []string{constants.AuthTypeSystem}
		} else if authType == constants.AuthTypeUser {
			if isMix {
				headers[constants.HTTPHeaderKeyAuthType] = []string{constants.AuthTypeMixUserSystem}
			} else {
				headers[constants.HTTPHeaderKeyAuthType] = []string{constants.AuthTypeUser}
			}
		}
	}

	return headers
}

func SetFuncAPINameToCtx(ctx context.Context, funcAPIName string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyFunctionAPIName, funcAPIName)
}

func GetFuncAPINameFromCtx(ctx context.Context) string {
	cast, _ := ctx.Value(constants.CtxKeyFunctionAPIName).(string)
	return cast
}
func SetUserAndAuthTypeToHeaders(ctx context.Context, headers map[string][]string) map[string][]string {
	return SetUserAndMixAuthTypeToHeaders(ctx, headers, false)
}

func SetFunctionMetaConfToCtx(ctx context.Context, metaConf map[string]string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	apiNameToMetaConf := map[string]*structs.FunctionMeta{}
	for apiName, conf := range metaConf {
		functionMeta := structs.FunctionMeta{}
		err := json.Unmarshal([]byte(conf), &functionMeta)
		if err != nil {
			fmt.Printf("SetFunctionMetaConfToCtx failed, err: %+v", err)
			return nil
		}
		apiNameToMetaConf[strings.ToLower(apiName)] = &functionMeta
	}

	return context.WithValue(ctx, constants.CtxKeyFunctionMetaConf, apiNameToMetaConf)
}

func GetFunctionMetaConfFromCtx(ctx context.Context, apiName string) *structs.FunctionMeta {
	if ctx == nil {
		return nil
	}

	metaConfMap, _ := ctx.Value(constants.CtxKeyFunctionMetaConf).(map[string]*structs.FunctionMeta)
	conf, _ := metaConfMap[strings.ToLower(apiName)]
	return conf
}

func GetCurFunctionMetaConfFromCtx(ctx context.Context) *structs.FunctionMeta {
	if ctx == nil {
		return nil
	}

	metaConfMap, _ := ctx.Value(constants.CtxKeyFunctionMetaConf).(map[string]*structs.FunctionMeta)
	conf, _ := metaConfMap[strings.ToLower(GetFuncAPINameFromCtx(ctx))]
	return conf
}

// 参数对应的无权限字段
func GetParamUnauthFieldMapFromCtx(ctx context.Context) (keyToUnauthFields map[string]interface{}) {
	if ctx == nil {
		return nil
	}

	userContext := GetUserContext(ctx)
	return userContext.Permission.UnauthFields
}

func GetParamUnauthFieldByKey(ctx context.Context, key string) (unauthFields interface{}) {
	return GetParamUnauthFieldMapFromCtx(ctx)[key]
}

func GetParamUnauthFieldRecordByKey(ctx context.Context, key string) (unauthFields []string) {
	return ParseStrList(GetParamUnauthFieldMapFromCtx(ctx)[key])
}

func GetParamUnauthFieldRecordListByKey(ctx context.Context, key string) (unauthFieldsList [][]string) {
	return ParseStrsList(GetParamUnauthFieldMapFromCtx(ctx)[key])
}

// RecordUnauthField 记录对应的无权限字段
type RecordUnauthField = map[string]map[int64][]string

func SetRecordUnauthField(ctx context.Context, unauthFieldMap RecordUnauthField) context.Context {
	return context.WithValue(ctx, constants.CtxKeyUnauthFieldMap, unauthFieldMap)
}

func GetRecordUnauthField(ctx context.Context) (objToRecordIDToUnauthFields RecordUnauthField) {
	if ctx == nil {
		return nil
	}

	objToRecordIDToUnauthFields, _ = ctx.Value(constants.CtxKeyUnauthFieldMap).(RecordUnauthField)
	return objToRecordIDToUnauthFields
}

func GetRecordUnauthFieldByObject(ctx context.Context, objectAPIName string) (recordIDToUnauthFields map[int64][]string) {
	if ctx == nil {
		return nil
	}

	recordIDToUnauthFields, _ = GetRecordUnauthField(ctx)[objectAPIName]
	return recordIDToUnauthFields
}

func GetRecordUnauthFieldByObjectAndRecordID(ctx context.Context, objectAPIName string, recordID int64) (unauthFields []string) {
	if ctx == nil {
		return nil
	}
	unauthFields, _ = GetRecordUnauthFieldByObject(ctx, objectAPIName)[recordID]
	return unauthFields
}
