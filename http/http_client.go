// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/byted-apaas/server-common-go/constants"
	exp "github.com/byted-apaas/server-common-go/exceptions"
	"github.com/byted-apaas/server-common-go/utils"
	"github.com/byted-apaas/server-common-go/utils/format"
	"github.com/byted-apaas/server-common-go/version"
)

type ClientType int

const (
	OpenAPIClient ClientType = iota + 1
	FaaSInfraClient
)

type HttpClient struct {
	Type ClientType
	http.Client
	MeshClient        *http.Client
	FromSDK           version.ISDKInfo
	rateLimitLogCount int64
}

var (
	openapiClientOnce sync.Once
	openapiClient     *HttpClient

	fsInfraClientOnce sync.Once
	fsInfraClient     *HttpClient
)

func GetOpenapiClient() *HttpClient {
	openapiClientOnce.Do(func() {
		openapiClient = &HttpClient{
			Type: OpenAPIClient,
			Client: http.Client{
				Transport: &http.Transport{
					DialContext:         TimeoutDialer(constants.HttpClientDialTimeoutDefault, 0),
					TLSHandshakeTimeout: constants.HttpClientTLSTimeoutDefault,
					MaxIdleConns:        1000,
					MaxIdleConnsPerHost: 10,
					IdleConnTimeout:     60 * time.Second,
				},
			},
			FromSDK: version.GetCommonSDKInfo(),
		}

		if utils.EnableMesh() {
			openapiClient.MeshClient = &http.Client{
				Transport: &http.Transport{
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						unixAddr, err := net.ResolveUnixAddr("unix", utils.GetSocketAddr())
						if err != nil {
							return nil, err
						}
						return net.DialUnix("unix", nil, unixAddr)
					},
					TLSHandshakeTimeout: constants.HttpClientTLSTimeoutDefault,
					MaxIdleConns:        1000,
					MaxIdleConnsPerHost: 10,
					IdleConnTimeout:     60 * time.Second,
				},
			}
		}
	})
	return openapiClient
}

func GetFaaSInfraClient(ctx context.Context) *HttpClient {
	fsInfraClientOnce.Do(func() {
		fsInfraClient = &HttpClient{
			Type: FaaSInfraClient,
			Client: http.Client{
				Transport: &http.Transport{
					DialContext:         TimeoutDialer(constants.HttpClientDialTimeoutDefault, 0),
					TLSHandshakeTimeout: constants.HttpClientTLSTimeoutDefault,
					MaxIdleConns:        1000,
					MaxIdleConnsPerHost: 10,
					IdleConnTimeout:     60 * time.Second,
				},
			},
			FromSDK: version.GetCommonSDKInfo(),
		}
	})

	if utils.EnableMesh() {
		fsInfraClient.MeshClient = &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					unixAddr, err := net.ResolveUnixAddr("unix", utils.GetSocketAddr())
					if err != nil {
						return nil, err
					}
					return net.DialUnix("unix", nil, unixAddr)
				},
				TLSHandshakeTimeout: constants.HttpClientTLSTimeoutDefault,
				MaxIdleConns:        1000,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     60 * time.Second,
			},
		}
	}
	return fsInfraClient
}

func (c *HttpClient) getActualDomain(ctx context.Context) string {
	switch c.Type {
	case OpenAPIClient:
		return utils.GetOpenAPIDomain(ctx)
	case FaaSInfraClient:
		return utils.GetFaaSInfraDomain(ctx)
	default:
		return ""
	}
}

func (c *HttpClient) doRequest(ctx context.Context, req *http.Request, headers map[string][]string, reqBody []byte, midList []ReqMiddleWare) ([]byte, map[string]interface{}, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// 限流控制
	err := c.checkPodRateLimit(ctx)
	if err != nil {
		return nil, nil, err
	}

	// 反压降速控制
	checkPressureAndDecelerate(ctx)

	// 执行中间件
	for _, mid := range midList {
		err = mid(ctx, req)
		if err != nil {
			return nil, nil, err
		}
	}

	// 设置 header 与 context
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	ctx = c.appendContextAndHeaders(ctx, req)

	// 超时控制
	var cancel context.CancelFunc
	ctx, cancel = GetTimeoutCtx(ctx)
	defer cancel()

	psm, cluster := utils.GetOpenAPIPSMAndCluster(ctx)
	if c.Type == FaaSInfraClient {
		psm, cluster = utils.GetFaaSInfraPSMFromEnv()
	}

	isUseMesh := utils.OpenMesh(ctx) && psm != "" && cluster != "" && c.MeshClient != nil

	if isUseMesh {
		req, err = c.transferToMeshReq(ctx, req, psm, cluster)
		if err != nil {
			return nil, nil, err
		}
	}

	start := time.Now()
	var resp *http.Response
	var respBody []byte
	defer func() {
		c.logRequest(ctx, req, resp, err, reqBody, respBody, start)
	}()

	// 执行请求
	_ = utils.InvokeFuncWithRetry(2, 5*time.Millisecond, func() error {
		if isUseMesh {
			resp, err = c.MeshClient.Do(req.WithContext(ctx))
		} else {
			resp, err = c.Do(req.WithContext(ctx)) // 走 dns
		}

		// 重试逻辑：只在 dial 超时错误时重试
		var opErr *net.OpError
		if errors.As(err, &opErr) && opErr.Op == "dial" && opErr.Timeout() {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, exp.InternalError("doRequest failed, err: %v, logid: %v", err, utils.GetLogIDFromCtx(ctx))
	}
	if resp == nil {
		return nil, nil, exp.InternalError("doRequest failed, resp is nil, logid: %v", utils.GetLogIDFromCtx(ctx))
	}

	extra, ctx := c.extractResponseInfo(ctx, resp)

	if resp.Body != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, extra, exp.InternalError("doRequest readBody failed, err: %v, logid: %v", err, utils.GetLogIDFromCtx(ctx))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, extra, exp.InternalError("doRequest failed, statusCode is %d, logid: %v, respBody: %s", resp.StatusCode, utils.GetLogIDFromCtx(ctx), string(respBody))
	}

	return respBody, extra, nil
}

func (c *HttpClient) Get(ctx context.Context, path string, headers map[string][]string, midList ...ReqMiddleWare) ([]byte, map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, c.getActualDomain(ctx)+path, nil)
	if err != nil {
		return nil, nil, exp.InternalError("HttpClient.Get failed, err: %v", err)
	}

	return c.doRequest(ctx, req, headers, nil, midList)
}

func (c *HttpClient) PostJson(ctx context.Context, path string, headers map[string][]string, data interface{}, midList ...ReqMiddleWare) ([]byte, map[string]interface{}, error) {
	body, err := utils.JsonMarshalBytes(data)
	if err != nil {
		return nil, nil, exp.InternalError("HttpClient.PostJson failed, err: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.getActualDomain(ctx)+path, bytes.NewReader(body))
	if err != nil {
		return nil, nil, exp.InternalError("HttpClient.PostJson failed, err: %v", err)
	}

	if headers == nil {
		headers = map[string][]string{}
	}
	headers[constants.HttpHeaderKeyContentType] = []string{constants.HttpHeaderValueJson}
	return c.doRequest(ctx, req, headers, body, midList)
}

func (c *HttpClient) PostBson(ctx context.Context, path string, headers map[string][]string, data interface{}, midList ...ReqMiddleWare) ([]byte, map[string]interface{}, error) {
	body, err := bson.Marshal(data)
	if err != nil {
		return nil, nil, exp.InternalError("HttpClient.PostBson failed, err: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.getActualDomain(ctx)+path, bytes.NewReader(body))
	if err != nil {
		return nil, nil, exp.InternalError("HttpClient.PostBson failed, err: %v", err)
	}

	if headers == nil {
		headers = map[string][]string{}
	}
	headers[constants.HttpHeaderKeyContentType] = []string{constants.HttpHeaderValueBson}
	return c.doRequest(ctx, req, headers, body, midList)
}

func (c *HttpClient) PostFormData(ctx context.Context, path string, headers map[string][]string, body *bytes.Buffer, midList ...ReqMiddleWare) ([]byte, map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodPost, c.getActualDomain(ctx)+path, body)
	if err != nil {
		return nil, nil, exp.InternalError("HttpClient.PostFormData failed, err: %v", err)
	}
	return c.doRequest(ctx, req, headers, body.Bytes(), midList)
}

func (c *HttpClient) PatchJson(ctx context.Context, path string, headers map[string][]string, data interface{}, midList ...ReqMiddleWare) ([]byte, map[string]interface{}, error) {
	body, err := utils.JsonMarshalBytes(data)
	if err != nil {
		return nil, nil, exp.InternalError("HttpClient.PatchJson failed, err: %v", err)
	}

	req, err := http.NewRequest(http.MethodPatch, c.getActualDomain(ctx)+path, bytes.NewReader(body))
	if err != nil {
		return nil, nil, exp.InternalError("HttpClient.PatchJson failed, err: %v", err)
	}

	if headers == nil {
		headers = map[string][]string{}
	}
	headers[constants.HttpHeaderKeyContentType] = []string{constants.HttpHeaderValueJson}
	return c.doRequest(ctx, req, headers, body, midList)
}

func (c *HttpClient) DeleteJson(ctx context.Context, path string, headers map[string][]string, data interface{}, midList ...ReqMiddleWare) ([]byte, map[string]interface{}, error) {
	body, err := utils.JsonMarshalBytes(data)
	if err != nil {
		return nil, nil, exp.InternalError("HttpClient.DeleteJson failed, err: %v", err)
	}

	req, err := http.NewRequest(http.MethodDelete, c.getActualDomain(ctx)+path, bytes.NewReader(body))
	if err != nil {
		return nil, nil, exp.InternalError("HttpClient.DeleteJson failed, err: %v", err)
	}

	if headers == nil {
		headers = map[string][]string{}
	}
	headers[constants.HttpHeaderKeyContentType] = []string{constants.HttpHeaderValueJson}
	return c.doRequest(ctx, req, headers, body, midList)
}

func (c *HttpClient) appendContextAndHeaders(ctx context.Context, req *http.Request) context.Context {
	req.Header = utils.SetUserAndAuthTypeToHeaders(ctx, req.Header)

	// 添加 aPaaS 的 LaneID
	req.Header.Add(constants.HTTPHeaderKeyFaaSLaneID, utils.GetFaaSLaneIDFromCtx(ctx))

	// 添加 aPaaS 的环境标识
	req.Header.Add(constants.HTTPHeaderKeyFaaSEnvID, utils.GetFaaSEnvIDFromCtx(ctx))
	req.Header.Add(constants.HTTPHeaderKeyFaaSEnvType, fmt.Sprintf("%d", utils.GetFaaSEnvTypeFromCtx(ctx)))

	// 透传 BOE&PPE 环境标识
	env := utils.GetTTEnvFromCtx(ctx)
	if strings.HasPrefix(env, "ppe_") {
		req.Header.Add(constants.HttpHeaderKeyUsePPE, "1")
		req.Header.Add(constants.HttpHeaderKeyEnv, env)
	} else if strings.HasPrefix(env, "boe_") {
		req.Header.Add(constants.HttpHeaderKeyUseBOE, "1")
		req.Header.Add(constants.HttpHeaderKeyEnv, env)
	}
	ctx = utils.SetKEnvToCtxForRPC(ctx)

	// 透传 lane 隔离标识
	lane := utils.GetAPaaSLaneFromCtx(ctx)
	if lane != "" {
		req.Header.Add(constants.HTTPHeaderKeyAPaaSLane, lane)
	}

	req.Header.Add(constants.HttpHeaderKeyLogID, utils.GetLogIDFromCtx(ctx))

	// trace
	for k, v := range utils.GetTraceHeader(ctx) {
		req.Header.Set(k, v)
	}

	switch c.Type {
	case OpenAPIClient:
		req.Header.Add(constants.HttpHeaderKeySDKFuncMsg, getSDKFuncMsgValue(ctx))
	case FaaSInfraClient:
		req.Header.Add(constants.HttpHeaderKeyOrgID, utils.GetEnvOrgID())
	}

	// append aPaaS persist key
	utils.SetAPaaSPersistHeader(ctx, req.Header)
	ctx = utils.WithAPaaSPersistFaaSValue(ctx, constants.PersistFaaSKeyFaaSType, utils.GetFaaSType(ctx))
	if c.FromSDK != nil {
		ctx = utils.WithAPaaSPersistFaaSValue(ctx, constants.PersistFaaSKeyFromSDKName, c.FromSDK.GetSDKName())
		ctx = utils.WithAPaaSPersistFaaSValue(ctx, constants.PersistFaaSKeyFromSDKVersion, c.FromSDK.GetVersion())
	}
	req.Header.Add(constants.PersistFaaSKeySummarized, utils.GetAPaaSPersistFaaSMapStr(ctx))

	return ctx
}

// extractResponseInfo
func (c *HttpClient) extractResponseInfo(ctx context.Context, resp *http.Response) (map[string]interface{}, context.Context) {
	extra := make(map[string]interface{})

	if resp != nil && resp.Header != nil {
		logID := resp.Header.Get(constants.HttpHeaderKeyLogID)
		extra[constants.HttpHeaderKeyLogID] = logID
		ctx = utils.SetLogIDToCtx(ctx, logID)
	}

	return extra, ctx
}

func (c *HttpClient) transferToMeshReq(ctx context.Context, req *http.Request, psm, cluster string) (*http.Request, error) {
	meshReq, err := http.NewRequest(req.Method, "http://127.0.0.1"+req.URL.Path, req.Body)
	if err != nil {
		return nil, exp.InternalError("new meshReq failed, err: %v, logid: %v", err, utils.GetLogIDFromCtx(ctx))
	}

	if meshReq.Header == nil {
		meshReq.Header = map[string][]string{}
	}

	for key, values := range req.Header {
		for _, value := range values {
			meshReq.Header.Add(key, value)
		}
	}

	meshReq.Header.Set("destination-service", psm)
	meshReq.Header.Set("destination-cluster", cluster)
	meshReq.Header.Set("destination-request-timeout", strconv.FormatInt(utils.GetMeshDestReqTimeout(ctx), 10))

	return meshReq, nil
}

func (c *HttpClient) checkPodRateLimit(ctx context.Context) error {
	// 反压信号检测请求，不进行限流
	if checkPressureSdkReqTag(ctx) {
		return nil
	}

	// 重置限流配额
	quota := utils.GetPodRateLimitQuotaFromCtx(ctx)
	oldQuota := limiter.maxRequest
	if reset := limiter.ResetRateLimiter(quota); reset && utils.GetDebugTypeFromCtx(ctx) == 0 { // debug 态不输出此日志
		fmt.Println(fmt.Sprintf("%s rate limit reset from %d to %d, apiID: %s, tenantID: %d, namespace: %s",
			utils.GetFormatDate(), oldQuota, quota, utils.GetFuncAPINameFromCtx(ctx), utils.GetTenantIDFromCtx(ctx), utils.GetNamespaceFromCtx(ctx)))
	}

	// 未触发限流，请求放行
	if limiter.AllowRequest() {
		return nil
	}

	// 触发限流，记录日志
	rateLimitMsg := fmt.Sprintf("SDK request exceeded %d QPS per-instance rate limit, please reduce call frequency.", quota)
	rateLimitLog := utils.NewFormatLog(ctx, utils.LogLevelWarn, constants.RateLimitLogType, rateLimitMsg)
	if c.rateLimitLogCount < utils.LogCountLimit {
		c.rateLimitLogCount++
		fmt.Println(rateLimitLog.String())
	}

	// 触发限流，禁止访问
	if downgrade := utils.GetPodRateLimitDowngradeFromCtx(ctx); !downgrade {
		return fmt.Errorf(rateLimitMsg)
	}

	// 触发限流，降级通过
	return nil
}

func (c *HttpClient) logRequest(ctx context.Context, req *http.Request, resp *http.Response, reqErr error, reqBody, respBody []byte, startTime time.Time) {
	// debug 模式跳过日志打印
	if utils.IsDebug(ctx) {
		return
	}

	defer utils.PanicGuard(ctx)

	statusCode := -1
	if resp != nil {
		statusCode = resp.StatusCode
	}

	// 简要日志（通过开关控制）
	if utils.GetSDKCallLogSwitchFromCtx(ctx) {
		bizStatusCode := gjson.GetBytes(respBody, "code").String()
		sdkCallLogMsg := utils.SDKCallLogMessage{
			FaaSPlatform:  utils.GetFaaSPlatform(),
			APIName:       utils.GetFunctionNameFromCtx(ctx),
			Language:      utils.GetAPaaSPersistFaaSValueFromCtx(ctx, constants.PersistFaaSKeyFuncLanguage),
			SDKName:       utils.GetAPaaSPersistFaaSValueFromCtx(ctx, constants.PersistFaaSKeyFromSDKName),
			SDKVersion:    utils.GetAPaaSPersistFaaSValueFromCtx(ctx, constants.PersistFaaSKeyFromSDKVersion),
			SDKAPI:        utils.GetApiTimeoutMethodFromCtx(ctx),
			Host:          req.Host,
			HTTPCode:      strconv.Itoa(statusCode),
			BizStatusCode: bizStatusCode,
			Cost:          time.Since(startTime).Milliseconds(),
		}
		logMsgBytes, _ := json.Marshal(sdkCallLogMsg)
		sdkCallLog := utils.NewFormatLog(ctx, utils.LogLevelInfo, constants.SDKCallLogType, string(logMsgBytes))
		fmt.Println(sdkCallLog.String())
	}

	// 详细日志（通过开关控制）
	if utils.GetSDKCallLogDetailSwitchFromCtx(ctx) {
		isFileTransfer := isFileTransferRequest(req)
		isToken := isTokenRequest(req)

		var sb strings.Builder
		sb.WriteString(utils.GetFormatDate())
		sb.WriteString("\n🍓")
		sb.WriteString(req.Method)
		sb.WriteString(" ")
		sb.WriteString(req.URL.String())
		sb.WriteString(fmt.Sprintf("\n🍋%d %+v %s", statusCode, time.Since(startTime), utils.GetLogIDFromCtx(ctx)))
		if reqErr != nil {
			sb.WriteString(fmt.Sprintf("\n❌error: %v", reqErr))
		}
		sb.WriteString("\n🍏request header:")
		sb.WriteString(formatHeaderSafe(req.Header))
		sb.WriteString("\n🍏request body:")
		sb.WriteString(formatBodySafe(reqBody, isFileTransfer, isToken))
		sb.WriteString("\n🍎response header:")
		respHeader := ""
		if resp != nil {
			respHeader = formatHeaderSafe(resp.Header)
		}
		sb.WriteString(respHeader)
		sb.WriteString("\n🍎response body:")
		sb.WriteString(formatBodySafe(respBody, isFileTransfer, isToken))
		fmt.Println(sb.String())
	}
}

func checkPressureAndDecelerate(ctx context.Context) {
	// 反压信号检测请求，不进行降速
	if checkPressureSdkReqTag(ctx) {
		return
	}

	// pressureDecelerator 需要在 webframe 请求进入前调用 InitPressureDecelerator 方法进行初始化，否则无法降速
	if pressureDecelerator == nil {
		return
	}

	// 反压降速开关未开启，不进行降速
	if !utils.GetPressureNeedDecelerateFromCtx(ctx) {
		return
	}

	// 反压降速检测
	key := utils.GetAPaaSPersistFaaSPressureSignalId(ctx)
	sleepTime := pressureDecelerator.GetSleeptime(key)

	// 未触发降速
	if sleepTime <= 0 {
		return
	}

	// 记录降速日志
	msg := utils.SpeedDownMessage{
		Key:       key,
		SleepTime: sleepTime,
	}
	msgBytes, _ := json.Marshal(msg)
	speedDownLog := utils.NewFormatLog(ctx, utils.LogLevelWarn, constants.SpeedDownLogType, string(msgBytes))
	fmt.Println(speedDownLog.String())

	// 执行降速
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
}

func GetTimeoutCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	timeoutMap, ok1 := ctx.Value(constants.CtxKeyAPITimeoutMap).(map[string]int64)
	method, ok2 := ctx.Value(constants.CtxKeyAPITimeoutMethod).(string)
	if ok1 && ok2 {
		timeout, ok := timeoutMap[method]
		if ok && timeout > 0 {
			return context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
		}
	}
	timeout, ok := constants.APITimeoutMapDefault[method]
	if ok && timeout > 0 {
		return context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	}
	return context.WithTimeout(ctx, constants.APITimeoutDefault)
}

func getSDKFuncMsgValue(ctx context.Context) string {
	funcMsgMap := map[string]interface{}{}
	funcMsgMap["funcApiName"] = utils.GetFunctionNameFromCtx(ctx)
	marshal, err := utils.JsonMarshalBytes(funcMsgMap)
	if err != nil {
		return ""
	}
	return string(marshal)
}

// TimeoutDialer 设置连接&读写超时，非 0 才设置
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(ctx context.Context, net, addr string) (c net.Conn, err error) {
	return func(ctx context.Context, netw, addr string) (net.Conn, error) {
		var conn net.Conn
		var err error
		if cTimeout != 0 {
			conn, err = net.DialTimeout(netw, addr, cTimeout)
			if err != nil {
				return nil, err
			}
			if rwTimeout != 0 {
				_ = conn.SetDeadline(time.Now().Add(rwTimeout))
			}
		}
		return conn, nil
	}
}

const (
	KB                 = 1024            // KB
	MaxSize            = 10 * KB         // 10KB
	LargeBodyThreshold = 2 * 1024 * 1024 // 2MB
)

// isTokenRequest 判断是否为获取 token 的请求
func isTokenRequest(req *http.Request) bool {
	if req == nil || req.URL == nil {
		return false
	}

	return strings.Contains(req.URL.Path, OpenapiPathGetToken)
}

// isFileTransferRequest 判断是否为文件上传/下载请求
func isFileTransferRequest(req *http.Request) bool {
	if req == nil || req.URL == nil {
		return false
	}

	// 根据 Content-Type 判断
	contentType := req.Header.Get(constants.HttpHeaderKeyContentType)
	if strings.HasPrefix(contentType, "multipart/form-data") {
		return true
	}

	// 根据 URL 路径判断文件上传/下载接口
	path := req.URL.Path
	fileTransferPaths := []string{
		"/attachment/",
		"/api/attachment/",
	}
	for _, p := range fileTransferPaths {
		if strings.Contains(path, p) {
			return true
		}
	}

	return false
}

// formatHeaderSafe 安全地格式化 header，对敏感信息进行脱敏处理
func formatHeaderSafe(header http.Header) string {
	if header == nil {
		return ""
	}

	// 复制 header 避免修改原始数据
	safeHeader := make(http.Header)
	for k, v := range header {
		if k == constants.HttpHeaderKeyAuthorization {
			safeHeader[k] = []string{"***"}
		} else {
			safeHeader[k] = v
		}
	}

	return format.Any(safeHeader)
}

// formatBodySafe 安全地格式化 body，对于大 body 直接返回摘要信息避免 OOM，对于 token 请求跳过 body 避免泄露
func formatBodySafe(v interface{}, isFileTransfer, isToken bool) string {
	if v == nil {
		return ""
	}

	// 对于 token 请求，跳过 body 内容避免 token 泄露
	if isToken {
		return "[token request, body skipped]"
	}

	// 对于文件传输请求，直接跳过 body 内容
	if isFileTransfer {
		if b, ok := v.([]byte); ok {
			return fmt.Sprintf("[file content, size=%d bytes, skipped]", len(b))
		}
		return "[file content, skipped]"
	}

	// 对于 []byte 类型，先检查长度再格式化
	if b, ok := v.([]byte); ok {
		if len(b) > LargeBodyThreshold {
			return fmt.Sprintf("[large body, size=%d bytes, skipped]", len(b))
		}
	}

	str := format.Any(v)
	if len(str) > MaxSize {
		return str[:MaxSize] + fmt.Sprintf(">>> 🥥 truncated(%d) > 10KB", len(str))
	}

	return str
}
