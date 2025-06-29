// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tidwall/gjson"

	"github.com/byted-apaas/server-common-go/constants"
	exp "github.com/byted-apaas/server-common-go/exceptions"
	"github.com/byted-apaas/server-common-go/structs"
)

func GetEnv() string {
	return os.Getenv(constants.EnvKEnvironment)
}

func GetBizIDC() string {
	return os.Getenv(constants.EnvKBizIDC)
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

func GetOpenAPIPSMAndCluster(ctx context.Context) (psm string, cluster string) {
	return GetLGWPSMFromEnv(), GetLGWClusterFromEnv()
}

func GetOpenAPIDomain(ctx context.Context) string {
	// 取配置优先级：运行时环境变量 > 上下文 > SDK 配置文件
	openAPIDomain := ""
	if openAPIDomain = os.Getenv(constants.EnvKOpenApiDomain); openAPIDomain != "" {
		return openAPIDomain
	}

	if ctx != nil {
		if openAPIDomain, _ = ctx.Value(constants.CtxKeyOpenapiDomain).(string); openAPIDomain != "" {
			return openAPIDomain
		}
	}
	return GetOpenAPIDomainByConf(ctx)
}

func GetFaaSInfraDomain(ctx context.Context) string {
	domain := os.Getenv(constants.EnvKFaaSInfraDomain)
	if domain != "" {
		return domain
	}

	if ctx != nil {
		if domain, _ = ctx.Value(constants.CtxKeyFaaSInfraDomain).(string); domain != "" {
			return domain
		}
	}
	return ""
}

func GetAGWDomain(ctx context.Context) string {
	domain := os.Getenv(constants.EnvKInnerAPIDomain)
	if domain != "" {
		return domain
	}

	if ctx != nil {
		if domain, _ = ctx.Value(constants.CtxKeyAGWDomain).(string); domain != "" {
			return domain
		}
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
// Deprecated
func GetInnerAPIPSM(ctx context.Context) string {
	psm := GetInnerAPIPSMFromCtx(ctx)
	if psm != "" {
		return psm
	}
	conf, ok := constants.EnvConfMap[GetEnv()+GetBizIDC()+GetBoe(ctx)]
	if !ok {
		return ""
	}
	return conf.InnerAPIPSM
}

func GetOpenAPIDomainByConf(ctx context.Context) string {
	conf, ok := constants.EnvConfMap[GetEnv()+GetBizIDC()+GetBoe(ctx)]
	if !ok {
		return ""
	}
	return conf.OpenAPIDomain
}

func GetAGWDomainByConf(ctx context.Context) string {
	conf, ok := constants.EnvConfMap[GetEnv()+GetBizIDC()+GetBoe(ctx)]
	if !ok {
		return ""
	}
	return conf.InnerAPIDomain
}

// Deprecated
func GetFaaSInfraPSM(ctx context.Context) string {
	conf, ok := constants.EnvConfMap[GetEnv()+GetBizIDC()+GetBoe(ctx)]
	if !ok {
		return ""
	}
	return conf.FaaSInfraPSM
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
	data, err := JsonMarshalBytes(v)
	if err != nil {
		return ""
	}
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

func ErrorWrapper(body []byte, extra map[string]interface{}, err error) ([]byte, error) {
	if err != nil {
		return nil, exp.ErrWrap(err)
	}

	code := gjson.GetBytes(body, "code").String()
	msg := gjson.GetBytes(body, "msg").String()
	switch code {
	case exp.SCFileDownload:
		return body, nil
	case exp.SCSuccess:
		data := gjson.GetBytes(body, "data")
		if data.Type == gjson.String {
			return []byte(data.Str), nil
		}
		return []byte(data.Raw), nil
	default:
		return nil, exp.NewErrWithCodeV2(code, msg, GetLogIDFromExtra(extra))
	}
}

// OpenMesh 是否开启 Mesh，有开关可以关闭 Mesh，有些场景不允许走 Mesh
func OpenMesh(ctx context.Context) bool {
	if IsCloseMesh(ctx) || IsDebug(ctx) || IsExternalFaaS() {
		return false
	}

	return EnableMesh()
}

// EnableMesh 是否支持 Mesh，通过 FaaS 的环境变量来判断
func EnableMesh() bool {
	return IsTrueString(os.Getenv(constants.EnvKMeshHttp)) && IsTrueString(os.Getenv(constants.EnvKMeshUDS)) && GetSocketAddr() != ""
}

func GetSocketAddr() string {
	return strings.TrimSpace(os.Getenv(constants.EnvKSocketAddr))
}

func IsTrueString(str string) bool {
	return strings.ToLower(str) == "true"
}

func IsExternalFaaS() bool {
	return GetFaaSPlatform() != "3"
}

// RequestToCurl 将http.Client和*http.Request转换为等效的curl命令字符串，同时不影响原请求
func RequestToCurl(client *http.Client, req *http.Request) (string, error) {
	var cmd []string
	cmd = append(cmd, "curl")

	// 添加verbose选项
	if client != nil && client.Transport != nil {
		cmd = append(cmd, "-v")
	}

	// 添加请求方法
	method := req.Method
	if method == "" {
		method = http.MethodGet
	}
	cmd = append(cmd, fmt.Sprintf("-X %s", method))

	// 添加请求头
	for key, values := range req.Header {
		for _, value := range values {
			escapedValue := strings.ReplaceAll(value, `"`, `\"`)
			cmd = append(cmd, fmt.Sprintf("-H \"%s: %s\"", key, escapedValue))
		}
	}

	// 处理请求体，确保不影响原请求
	var bodyBytes []byte
	var err error

	// 如果请求体存在且可以读取
	if req.Body != nil {
		// 克隆请求体
		bodyBytes, err = cloneRequestBody(req)
		if err != nil {
			return "", fmt.Errorf("克隆请求体失败: %v", err)
		}

		if len(bodyBytes) > 0 {
			contentType := req.Header.Get("Content-Type")

			if strings.Contains(contentType, "application/json") {
				escapedBody := strings.ReplaceAll(string(bodyBytes), `"`, `\"`)
				cmd = append(cmd, fmt.Sprintf("-d \"%s\"", escapedBody))
			} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
				formData, err := url.ParseQuery(string(bodyBytes))
				if err == nil {
					for key, values := range formData {
						for _, value := range values {
							cmd = append(cmd, fmt.Sprintf("-d \"%s=%s\"", key, url.QueryEscape(value)))
						}
					}
				} else {
					escapedBody := strings.ReplaceAll(string(bodyBytes), `"`, `\"`)
					cmd = append(cmd, fmt.Sprintf("-d \"%s\"", escapedBody))
				}
			} else {
				escapedBody := strings.ReplaceAll(string(bodyBytes), `"`, `\"`)
				cmd = append(cmd, fmt.Sprintf("-d \"%s\"", escapedBody))
			}
		}
	}

	// 添加请求URL
	escapedURL := req.URL.String()
	if strings.ContainsAny(escapedURL, " \"'") {
		escapedURL = fmt.Sprintf("'%s'", escapedURL)
	}
	cmd = append(cmd, escapedURL)

	// 添加超时设置
	if client != nil && client.Timeout > 0 {
		timeoutSeconds := int(client.Timeout / time.Second)
		if timeoutSeconds == 0 {
			timeoutSeconds = 1
		}
		cmd = append(cmd, fmt.Sprintf("--max-time %d", timeoutSeconds))
	}

	return strings.Join(cmd, " "), nil
}

// cloneRequestBody 克隆请求体，不影响原请求
func cloneRequestBody(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return nil, nil
	}

	// 读取原始请求体
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// 恢复原始请求体
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes, nil
}
