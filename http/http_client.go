// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/byted-apaas/server-common-go/constants"
	exp "github.com/byted-apaas/server-common-go/exceptions"
	"github.com/byted-apaas/server-common-go/utils"
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

func (c *HttpClient) doRequest(ctx context.Context, req *http.Request, headers map[string][]string, midList []ReqMiddleWare) ([]byte, map[string]interface{}, error) {
	quota := utils.GetPodRateLimitQuotaFromCtx(ctx)
	oldQuota := limiter.maxRequest
	if reset := limiter.ResetRateLimiter(quota); reset && utils.GetDebugTypeFromCtx(ctx) == 0 { // debug 态不输出此日志
		fmt.Println(fmt.Sprintf("%s rate limit reset from %d to %d, apiID: %s, tenantID: %d, namespace: %s", utils.GetFormatDate(), oldQuota, quota, utils.GetFuncAPINameFromCtx(ctx), utils.GetTenantIDFromCtx(ctx), utils.GetNamespaceFromCtx(ctx)))
	}

	if !limiter.AllowRequest() {
		format := "request limit exceeded quota: %d qps"
		formatLog := utils.FormatLog{
			Level:         utils.LogLevelWarn,
			EventID:       utils.GetExecuteIDFromCtx(ctx),
			FunctionAPIID: utils.GetFunctionAPIIDFromCtx(ctx),
			LogID:         utils.GetLogIDFromCtx(ctx),
			Timestamp:     time.Now().UnixNano() / 1e3, // 使用微秒
			Message:       fmt.Sprintf(format, quota),
			TenantID:      utils.GetTenantIDFromCtx(ctx),
			TenantType:    utils.GetTenantTypeFromCtx(ctx),
			Namespace:     utils.GetNamespaceFromCtx(ctx),
			LogType:       constants.RateLimitLogType,
		}

		if c.rateLimitLogCount < utils.LogCountLimit {
			c.rateLimitLogCount++
			fmtMessage := utils.GetFormatLogWithMessage(formatLog, c.rateLimitLogCount)
			content := fmt.Sprintf("%s %s %s %s", utils.GetFormatDate(), constants.APaaSLogPrefix, fmtMessage, constants.APaaSLogSuffix)
			fmt.Println(content)
		}
		// 触发限流，禁止访问
		if downgrade := utils.GetPodRateLimitDowngradeFromCtx(ctx); !downgrade {
			return nil, nil, fmt.Errorf("request limit exceeded quota: %d qps", quota)
		}
		// 触发限流，降级通过
	}

	extra := map[string]interface{}{}

	if ctx == nil {
		ctx = context.Background()
	}

	for _, mid := range midList {
		err := mid(ctx, req)
		if err != nil {
			return nil, nil, err
		}
	}

	headers = utils.SetUserAndAuthTypeToHeaders(ctx, headers)

	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 添加Apaas的LaneID
	req.Header.Add(constants.HTTPHeaderKeyFaaSLaneID, utils.GetFaaSLaneIDFromCtx(ctx))

	// 添加ApaaS的环境标识
	req.Header.Add(constants.HTTPHeaderKeyFaaSEnvID, utils.GetFaaSEnvIDFromCtx(ctx))
	req.Header.Add(constants.HTTPHeaderKeyFaaSEnvType, fmt.Sprintf("%d", utils.GetFaaSEnvTypeFromCtx(ctx)))

	ctx = c.requestCommonInfo(ctx, req)

	// Timeout
	var cancel context.CancelFunc
	ctx, cancel = GetTimeoutCtx(ctx)
	defer cancel()

	var resp *http.Response
	var err error

	// OpenAPIClient
	psm, cluster := utils.GetOpenAPIPSMAndCluster(ctx)
	switch c.Type {
	case FaaSInfraClient:
		psm, cluster = utils.GetFaaSInfraPSMFromEnv()
	}

	// 连接层超时
	isUseMesh := utils.OpenMesh(ctx) && psm != "" && cluster != "" && c.MeshClient != nil
	_ = utils.InvokeFuncWithRetry(2, 5*time.Millisecond, func() error {
		if isUseMesh {
			var newReq *http.Request
			newReq, err = http.NewRequest(req.Method, "http://127.0.0.1"+req.URL.Path, req.Body)
			if err != nil {
				return err
			}

			if newReq.Header == nil {
				newReq.Header = map[string][]string{}
			}

			for key, values := range req.Header {
				for _, value := range values {
					newReq.Header.Add(key, value)
				}
			}

			// 走 mesh
			newReq.Header.Set("destination-service", psm)
			newReq.Header.Set("destination-cluster", cluster)
			newReq.Header.Set("destination-request-timeout", strconv.FormatInt(utils.GetMeshDestReqTimeout(ctx), 10))
			resp, err = c.MeshClient.Do(newReq.WithContext(ctx))
		} else {
			// 走 dns
			resp, err = c.Do(req.WithContext(ctx))
		}
		var opErr *net.OpError
		if errors.As(err, &opErr) && opErr.Op == "dial" && opErr.Timeout() {
			return err
		}
		return nil
	})

	var logid string
	if resp != nil && resp.Header != nil {
		logid = resp.Header.Get(constants.HttpHeaderKeyLogID)
		extra[constants.HttpHeaderKeyLogID] = logid
		if isUseMesh {
			meshFlag := resp.Header.Get(constants.HTTPHeaderEnvoyRespFlag)
			extra[constants.HTTPHeaderEnvoyRespFlag] = meshFlag
			urlPath := ""
			if req != nil && req.URL != nil {
				urlPath = req.URL.Path
			}
			fmt.Printf("doRequest mesh resp: logID: %v, mesh flag: %v, url path: %v", logid, meshFlag, urlPath)
		}
	}

	if err != nil {
		return nil, extra, exp.InternalError("doRequest failed, err: %v, logid: %v", err, logid)
	}

	if resp == nil {
		return nil, extra, exp.InternalError("doRequest failed: resp is nil, logid: %v", logid)
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	// Http resp body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, extra, exp.InternalError("doRequest readBody failed, err: %v, logid: %v", err, logid)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, extra, exp.InternalError("doRequest failed: statusCode is %d, data: %s, logid: %v", resp.StatusCode, string(data), logid)
	}

	return data, extra, nil
}

func (c *HttpClient) Get(ctx context.Context, path string, headers map[string][]string, midList ...ReqMiddleWare) ([]byte, map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, c.getActualDomain(ctx)+path, nil)
	if err != nil {
		return nil, nil, exp.InternalError("HttpClient.Get failed, err: %v", err)
	}

	return c.doRequest(ctx, req, headers, midList)
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
	return c.doRequest(ctx, req, headers, midList)
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
	return c.doRequest(ctx, req, headers, midList)
}

func (c *HttpClient) PostFormData(ctx context.Context, path string, headers map[string][]string, body *bytes.Buffer, midList ...ReqMiddleWare) ([]byte, map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodPost, c.getActualDomain(ctx)+path, body)
	if err != nil {
		return nil, nil, exp.InternalError("HttpClient.PostFormData failed, err: %v", err)
	}
	return c.doRequest(ctx, req, headers, midList)
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
	return c.doRequest(ctx, req, headers, midList)
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
	return c.doRequest(ctx, req, headers, midList)
}

func (c *HttpClient) requestCommonInfo(ctx context.Context, req *http.Request) context.Context {
	// common: ppe or boe
	env := utils.GetTTEnvFromCtx(ctx)
	if strings.HasPrefix(env, "ppe_") {
		req.Header.Add(constants.HttpHeaderKeyUsePPE, "1")
		req.Header.Add(constants.HttpHeaderKeyEnv, env)
	} else if strings.HasPrefix(env, "boe_") {
		req.Header.Add(constants.HttpHeaderKeyUseBOE, "1")
		req.Header.Add(constants.HttpHeaderKeyEnv, env)
	}

	// lane
	lane := utils.GetAPaaSLaneFromCtx(ctx)
	if lane != "" {
		req.Header.Add(constants.HTTPHeaderKeyAPaaSLane, lane)
	}

	// add aPaaS persist key
	utils.SetAPaaSPersistHeader(ctx, req.Header)

	req.Header.Add(constants.HttpHeaderKeyLogID, utils.GetLogIDFromCtx(ctx))

	// trace
	for k, v := range utils.GetTraceHeader(ctx) {
		req.Header.Set(k, v)
	}

	// divide open-api & faaS—infra
	switch c.Type {
	case OpenAPIClient:
		req.Header.Add(constants.HttpHeaderKeySDKFuncMsg, getSDKFuncMsgValue(ctx))
		// 运行时 faas 信息透传
		ctx = utils.WithAPaaSPersistFaaSValue(ctx, constants.PersistFaaSKeyFaaSType, utils.GetFaaSType(ctx))
		if c.FromSDK != nil {
			ctx = utils.WithAPaaSPersistFaaSValue(ctx, constants.PersistFaaSKeyFromSDKName, c.FromSDK.GetSDKName())
			ctx = utils.WithAPaaSPersistFaaSValue(ctx, constants.PersistFaaSKeyFromSDKVersion, c.FromSDK.GetVersion())
		}
		req.Header.Add(constants.PersistFaaSKeySummarized, utils.GetAPaaSPersistFaaSMapStr(ctx))
	case FaaSInfraClient:
		req.Header.Add(constants.HttpHeaderKeyOrgID, utils.GetEnvOrgID())
		req.Header.Add(constants.PersistFaaSKeySummarized, utils.GetAPaaSPersistFaaSMapStr(ctx))
	}

	return utils.SetKEnvToCtxForRPC(ctx)
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
