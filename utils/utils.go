// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package utils

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/byted-apaas/server-common-go/constants"
	exp "github.com/byted-apaas/server-common-go/exceptions"
	"github.com/byted-apaas/server-common-go/structs"
)

func GetEnv() string {
	return os.Getenv(constants.EnvKEnvironment)
}

func GetTenantName() string {
	return os.Getenv(constants.EnvKTenantName)
}

func GetNamespace() string {
	return os.Getenv(constants.EnvKNamespace)
}

func GetServiceID() string {
	return os.Getenv(constants.EnvKSvcID)
}

func GetOpenAPIDomain(ctx context.Context) string {
	if ctx != nil {
		if domain, _ := ctx.Value(constants.CtxKeyOpenapiDomain).(string); domain != "" {
			return domain
		}
	}
	return os.Getenv(constants.EnvKOpenApiDomain)
}

func GetFaaSInfraDomain(ctx context.Context) string {
	if ctx != nil {
		if domain, _ := ctx.Value(constants.CtxKeyFaaSInfraDomain).(string); domain != "" {
			return domain
		}
	}
	return os.Getenv(constants.EnvKFaaSInfraDomain)
}

func GetAGWDomain(ctx context.Context) string {
	if ctx != nil {
		if domain, _ := ctx.Value(constants.CtxKeyAGWDomain).(string); domain != "" {
			return domain
		}
	}
	domain := os.Getenv(constants.EnvKInnerAPIDomain)
	if len(domain) > 0 {
		return domain
	}
	return GetAGWDomainByConf(ctx)
}

func GetEnvOrgID() string {
	return os.Getenv(constants.EnvKOrgID)
}

func GetAppIDAndSecret() (string, string, error) {
	tenantName := os.Getenv(constants.EnvKTenantName)
	namespace := os.Getenv(constants.EnvKNamespace)
	dClientID := os.Getenv(constants.EnvKClientID)
	dClientSecret := os.Getenv(constants.EnvKClientSecret)
	if tenantName == "" || namespace == "" || dClientID == "" || dClientSecret == "" {
		return "", "", exp.InternalError("Missing params in env.")
	}

	key := paddingN([]byte(tenantName+namespace), 32)
	clientID, err := AesDecryptText(0, key, dClientID)
	if err != nil {
		return "", "", exp.InternalError("Decrypt ClientID err: %v", err)
	}
	clientSecret, err := AesDecryptText(0, key, dClientSecret)
	if err != nil {
		return "", "", exp.InternalError("Decrypt ClientSecret err: %v", err)
	}
	return clientID, clientSecret, nil
}

func StrInStrs(strs []string, str string) bool {
	for _, v := range strs {
		if str == v {
			return true
		}
	}
	return false
}

func IntInInts(ns []int, n int) bool {
	for _, v := range ns {
		if n == v {
			return true
		}
	}
	return false
}

func Int64InInt64s(ns []int, n int) bool {
	for _, v := range ns {
		if n == v {
			return true
		}
	}
	return false
}

func IsLocalDebug(ctx context.Context) bool {
	return GetDebugTypeFromCtx(ctx) == constants.DebugTypeLocal
}

// IsDebug 是否调试
func IsDebug(ctx context.Context) bool {
	debugType := GetDebugTypeFromCtx(ctx)
	return debugType == constants.DebugTypeOnline || debugType == constants.DebugTypeLocal
}

func Int64Ptr(val int64) *int64 {
	return &val
}

func Int64ValueOfPtr(p *int64, defaultVal int64) int64 {
	if p == nil {
		return defaultVal
	}
	return *p
}

func IntPtr(val int) *int {
	return &val
}

func IntValueOfPtr(p *int, defaultVal int) int {
	if p == nil {
		return defaultVal
	}
	return *p
}

func BoolPtr(val bool) *bool {
	return &val
}

func PtrToInt(p *int, defaultVal int) int {
	if p == nil {
		return defaultVal
	}
	return *p
}

func BoolValueOfPtr(val *bool) bool {
	if val == nil {
		return false
	}
	return *val
}

func StringPtr(val string) *string {
	valPtr := new(string)
	*valPtr = val
	return valPtr
}

func StringValueOfPtr(p *string, defaultVal string) string {
	if p == nil {
		return defaultVal
	}
	return *p
}

// GetInnerAPIPSM
// open-sdk: from ctx
// faaS-sdk: from const by env
func GetInnerAPIPSM(ctx context.Context) string {
	psm := GetInnerAPIPSMFromCtx(ctx)
	if psm != "" {
		return psm
	}
	conf, ok := constants.EnvConfMap[GetEnv()+GetBoe(ctx)]
	if !ok {
		return ""
	}
	return conf.InnerAPIPSM
}

func GetAGWDomainByConf(ctx context.Context) string {
	conf, ok := constants.EnvConfMap[GetEnv()+GetBoe(ctx)]
	if !ok {
		return ""
	}
	return conf.InnerAPIDomain
}

func GetBoe(ctx context.Context) string {
	return GetEnvBoeFromCtx(ctx)
}

func GetLogIDFromExtra(extra map[string]interface{}) string {
	if logid, ok := extra[constants.HttpHeaderKeyLogID].(string); ok {
		return logid
	}
	return ""
}

func SetKEnvToCtxForRPC(ctx context.Context) context.Context {
	env := GetTTEnvFromCtx(ctx)
	if strings.HasPrefix(env, "ppe_") || strings.HasPrefix(env, "boe_") {
		ctx = SetKEnvToCtx(ctx, env)
	}
	return ctx
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

func IsMicroservice(ctx context.Context) bool {
	return os.Getenv(constants.EnvKFaaSScene) == "microservice"
}

func GetFaaSType(ctx context.Context) string {
	if IsMicroservice(ctx) {
		return constants.FaaSTypeMicroService
	}
	return constants.FaaSTypeFunction
}
func ToString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func GetRecordID(record interface{}) int64 {
	if record == nil {
		return 0
	}

	newRecord := structs.RecordOnlyID{}
	err := Decode(record, &newRecord)
	if err != nil {
		fmt.Printf("GetRecordID failed, err: %+v\n", err)
		return 0
	}

	return newRecord.GetID()
}

func ParseStrList(v interface{}) (strs []string) {
	strList, ok := v.([]interface{})
	if !ok {
		return nil
	}

	for _, str := range strList {
		if s, ok := str.(string); ok {
			strs = append(strs, s)
		}
	}

	return strs
}

func ParseStrsList(v interface{}) (strsList [][]string) {
	strList, ok := v.([]interface{})
	if !ok {
		return nil
	}

	for _, str := range strList {
		strsList = append(strsList, ParseStrList(str))
	}

	return strsList
}