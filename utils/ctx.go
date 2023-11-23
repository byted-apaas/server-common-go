// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/byted-apaas/server-common-go/constants"
	exp "github.com/byted-apaas/server-common-go/exceptions"
	"github.com/byted-apaas/server-common-go/structs"
)

func SetTenantToCtx(ctx context.Context, tenant *structs.Tenant) context.Context {
	return context.WithValue(ctx, constants.CtxKeyTenant, tenant)
}

func GetTenantFromCtx(ctx context.Context) (*structs.Tenant, error) {
	tenant := structs.Tenant{}
	err := Decode(ctx.Value(constants.CtxKeyTenant), &tenant)
	if err != nil {
		return nil, exp.InternalError("decode tenant failed, err: %+v", err)
	}

	return &tenant, nil
}

func GetTenantIDFromCtx(ctx context.Context) int64 {
	tenant, err := GetTenantFromCtx(ctx)
	if err != nil {
		return 0
	}

	return tenant.ID
}

func SetFaaSLaneIDCtx(ctx context.Context, laneID string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyLaneID, laneID)
}

func GetFaaSLaneIDFromCtx(ctx context.Context) string {
	cast, ok := ctx.Value(constants.CtxKeyLaneID).(string)
	if ok {
		return cast
	}
	return ""
}

func SetAppInfoToCtx(ctx context.Context, appInfo *structs.AppInfo) context.Context {
	return context.WithValue(ctx, constants.CtxKeyApp, appInfo)
}

func GetAppInfoFromCtx(ctx context.Context) (*structs.AppInfo, error) {
	appInfo := structs.AppInfo{}
	err := Decode(ctx.Value(constants.CtxKeyApp), &appInfo)
	if err != nil {
		return nil, exp.InternalError("Decode appInfo failed, err: %+v", err)
	}

	return &appInfo, nil
}

func SetEventInfoToCtx(ctx context.Context, appInfo *structs.EventInfo) context.Context {
	return context.WithValue(ctx, constants.CtxKeyEvent, appInfo)
}

func GetEventInfoFromCtx(ctx context.Context) (*structs.EventInfo, error) {
	eventInfo := structs.EventInfo{}
	err := Decode(ctx.Value(constants.CtxKeyEvent), &eventInfo)
	if err != nil {
		return nil, exp.InternalError("Decode eventInfo failed, err: %+v", err)
	}

	return &eventInfo, nil
}

func GetNamespaceFromCtx(ctx context.Context) string {
	tenant, err := GetTenantFromCtx(ctx)
	if err != nil {
		return ""
	}

	return tenant.Namespace
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

func SetFuncAPINameToCtx(ctx context.Context, funcAPIName string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyFunctionAPIName, funcAPIName)
}

func GetFuncAPINameFromCtx(ctx context.Context) string {
	cast, _ := ctx.Value(constants.CtxKeyFunctionAPIName).(string)
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

	userCtx := structs.UserContext{}
	err := Decode(ctx.Value(constants.CtxUserContext), &userCtx)
	if err != nil {
		fmt.Printf("Decode userCtx failed, err: %+v", err)
	}
	return userCtx
}

func SetUserContextMap(ctx context.Context, userCtxMap map[string]interface{}) context.Context {
	return context.WithValue(ctx, constants.CtxUserContextMap, userCtxMap)
}

func GetUserContextMap(ctx context.Context) (res map[string]interface{}) {
	defer func() {
		if res == nil {
			res = make(map[string]interface{})
		}
	}()
	if ctx == nil {
		return
	}
	res, _ = ctx.Value(constants.CtxUserContextMap).(map[string]interface{})
	return res
}

func SetAPaaSLaneToCtx(ctx context.Context, lane string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyAPaaSLane, lane)
}

func GetAPaaSLaneFromCtx(ctx context.Context) string {
	cast, _ := ctx.Value(constants.CtxKeyAPaaSLane).(string)
	return cast
}

func SetAPaaSPersistFaaSMapToCtx(ctx context.Context, aPaaSPersistFaaSMap map[string]string) context.Context {
	return context.WithValue(ctx, constants.PersistFaaSKeySummarized, aPaaSPersistFaaSMap)
}

var aPaaSPersistMutex sync.Mutex

func GetAPaaSPersistFaaSMapFromCtx(ctx context.Context) (copyRes map[string]string) {
	copyRes = map[string]string{}
	if ctx == nil {
		return
	}

	aPaaSPersistMutex.Lock()
	defer aPaaSPersistMutex.Unlock()

	res, _ := ctx.Value(constants.PersistFaaSKeySummarized).(map[string]string)
	// copy 是为了避免出现 map 并发读写的问题
	for k, v := range res {
		copyRes[k] = v
	}
	return
}

func WithAPaaSPersistFaaSValue(ctx context.Context, key, value string) context.Context {
	m := GetAPaaSPersistFaaSMapFromCtx(ctx)

	aPaaSPersistMutex.Lock()
	defer aPaaSPersistMutex.Unlock()

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

func SetUserAndAuthTypeToCtx(ctx context.Context, authType *string) context.Context {
	return SetUserAndMixAuthTypeToCtx(ctx, authType, false)
}

// SetUserAndMixAuthTypeToCtx 设置鉴权方式
// - 接口级配置优先级高于函数级
// - 函数级配置
// - oql 接口使用 system 和 mix_user_system
// - 除 oql 之外的接口使用 system 和 user
func SetUserAndMixAuthTypeToCtx(ctx context.Context, authType *string, isMix bool) context.Context {
	if authType == nil || *authType == "" {
		authType = GetGlobalAuthType()
	}

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
		} else if *authType == constants.AuthTypeMixUserSystem {
			ctx = context.WithValue(ctx, constants.AuthTypeKey, constants.AuthTypeMixUserSystem)
		}
	}
	return ctx
}

func SetAPaaSPersistHeader(ctx context.Context, header http.Header) {
	if ctx == nil {
		return
	}
	if persistHeader, ok := ctx.Value(constants.PersistAPaaSKeySummarized).(map[string]string); ok {
		for key, value := range persistHeader {
			header.Add(key, value)
		}
	}
}

func SetUserAndMixAuthTypeToHeaders(ctx context.Context, headers map[string][]string, isMix bool) map[string][]string {
	if headers == nil {
		headers = make(map[string][]string)
	}

	authType, ok := ctx.Value(constants.AuthTypeKey).(string)
	if !ok || authType == "" {
		if v := GetGlobalAuthType(); v != nil && *v != "" {
			authType = *v
		}
	}

	userID := GetUserIDFromCtx(ctx)
	headers[constants.HttpHeaderKeyUser] = []string{fmt.Sprintf("%d", userID)}
	if authType != "" {
		if userID == -1 || authType == constants.AuthTypeSystem {
			headers[constants.HTTPHeaderKeyAuthType] = []string{constants.AuthTypeSystem}
		} else if authType == constants.AuthTypeUser {
			if isMix {
				headers[constants.HTTPHeaderKeyAuthType] = []string{constants.AuthTypeMixUserSystem}
			} else {
				headers[constants.HTTPHeaderKeyAuthType] = []string{constants.AuthTypeUser}
			}
		} else if authType == constants.AuthTypeMixUserSystem {
			headers[constants.HTTPHeaderKeyAuthType] = []string{constants.AuthTypeMixUserSystem}
		}
	}

	return headers
}

func SetUserAndAuthTypeToHeaders(ctx context.Context, headers map[string][]string) map[string][]string {
	return SetUserAndMixAuthTypeToHeaders(ctx, headers, false)
}

// SetFunctionMetaConfToCtx 提供给框架层使用
func SetFunctionMetaConfToCtx(ctx context.Context, metaConf map[string]string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	apiNameToMetaConf := map[string]*structs.FunctionMeta{}
	for apiName, conf := range metaConf {
		functionMeta := structs.FunctionMeta{}
		err := json.Unmarshal([]byte(conf), &functionMeta)
		if err != nil {
			fmt.Printf("SetFunctionMetaConfToCtx failed, err: %+v\n", err)
			return ctx
		}
		apiNameToMetaConf[strings.ToLower(apiName)] = &functionMeta
	}

	return context.WithValue(ctx, constants.CtxKeyFunctionMetaConf, apiNameToMetaConf)
}

func GetFunctionMetaConfFromCtx(ctx context.Context, apiName string) *structs.FunctionMeta {
	if ctx == nil {
		return nil
	}

	apiNameToFuncMeta := map[string]*structs.FunctionMeta{}
	err := Decode(ctx.Value(constants.CtxKeyFunctionMetaConf), &apiNameToFuncMeta)
	if err != nil {
		fmt.Printf("Decode metaConf failed, err: %+v", err)
		return nil
	}

	conf, _ := apiNameToFuncMeta[strings.ToLower(apiName)]
	return conf
}

func GetCurFunctionMetaConfFromCtx(ctx context.Context) *structs.FunctionMeta {
	if ctx == nil {
		return nil
	}

	return GetFunctionMetaConfFromCtx(ctx, strings.ToLower(GetFuncAPINameFromCtx(ctx)))
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

	unauthField := RecordUnauthField{}
	err := Decode(ctx.Value(constants.CtxKeyUnauthFieldMap), &unauthField)
	if err != nil {
		fmt.Printf("Decode unauthField failed, err: %+v", err)
		return unauthField
	}

	return unauthField
}

func GetRecordUnauthFieldByObject(ctx context.Context, objectAPIName string) (recordIDToUnauthFields map[int64][]string) {
	if ctx == nil {
		return nil
	}

	// 加锁，避免并发读写问题
	unauthFieldMapMutex.Lock()
	defer unauthFieldMapMutex.Unlock()
	recordIDToUnauthFields, _ = GetRecordUnauthField(ctx)[objectAPIName]
	if recordIDToUnauthFields == nil {
		return nil
	}

	// 深 copy 避免后续 map 的并发读写问题
	result := map[int64][]string{}
	for recordID, unauthFields := range recordIDToUnauthFields {
		result[recordID] = append([]string{}, unauthFields...)
	}
	return result
}

func GetRecordUnauthFieldByObjectAndRecordID(ctx context.Context, objectAPIName string, recordID int64) (unauthFields []string) {
	if ctx == nil {
		return nil
	}
	unauthFields, _ = GetRecordUnauthFieldByObject(ctx, objectAPIName)[recordID]
	return unauthFields
}

var unauthFieldMapMutex sync.Mutex

func WriteUnauthFieldMapWithLock(ctx context.Context, objectAPIName string, recordID int64, unauthFields []string) {
	unauthFieldMapMutex.Lock()
	defer unauthFieldMapMutex.Unlock()

	unauthFieldMap := GetRecordUnauthField(ctx)
	if unauthFieldMap == nil {
		return
	}

	if v, ok := unauthFieldMap[objectAPIName]; !ok || v == nil {
		unauthFieldMap[objectAPIName] = map[int64][]string{}
	}

	unauthFieldMap[objectAPIName][recordID] = unauthFields
}

func GetSDKConf(ctx context.Context) *structs.SDKConf {
	if ctx == nil {
		return nil
	}

	sdkConfStr, ok := ctx.Value(constants.CtxKeySDKConf).(string)
	if !ok || sdkConfStr == "" {
		return nil
	}

	sdkConf := structs.SDKConf{}
	err := JsonUnmarshalBytes([]byte(sdkConfStr), &sdkConf)
	if err != nil {
		return nil
	}
	return &sdkConf
}

func GetSDKTransientConf(ctx context.Context) *structs.SDKTransientConf {
	sdkConf := GetSDKConf(ctx)
	if sdkConf == nil {
		return nil
	}
	return sdkConf.TransientConf
}

func IsCloseMesh(ctx context.Context) bool {
	transientConf := GetSDKTransientConf(ctx)
	if transientConf == nil {
		return false
	}

	return transientConf.IsCloseMesh
}

func GetTraceHeader(ctx context.Context) map[string]string {
	traceHeader := map[string]string{}
	if ctx == nil {
		return traceHeader
	}

	if traceParent, ok := ctx.Value("traceparent").(string); ok && traceParent != "" {
		traceHeader["traceparent"] = traceParent
	}

	if traceState, ok := ctx.Value("tracestate").(string); ok && traceState != "" {
		traceHeader["tracestate"] = traceState
	}

	return traceHeader
}

type RuntimeType string

const (
	RuntimeTypeRuntime    RuntimeType = "0" // 运行态
	RuntimeTypeCloudDebug RuntimeType = "1" // 云端调试
	RuntimeTypeLocalDebug RuntimeType = "2" // 本地调试
)

func SetRuntimeType(ctx context.Context, runtimeType string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, constants.CtxKeyRuntimeType, runtimeType)
}

func GetRuntimeType(ctx context.Context) RuntimeType {
	if ctx == nil {
		ctx = context.Background()
	}

	if runtimeType, ok := ctx.Value(constants.CtxKeyRuntimeType).(string); ok {
		return RuntimeType(runtimeType)
	}
	return RuntimeTypeRuntime
}

func IsRuntime(ctx context.Context) bool {
	if ctx == nil {
		ctx = context.Background()
	}

	return GetRuntimeType(ctx) == RuntimeTypeRuntime
}
