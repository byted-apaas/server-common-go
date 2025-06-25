package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/byted-apaas/server-common-go/constants"
	"github.com/byted-apaas/server-common-go/utils"
)

const (
	BatchQueryPressureSignalPath = "/v1/namespaces/:namespace/arch_service/pressure/batch_query"
)

// IPressureHttpClient 反压中心 http client
type IPressureHttpClient interface {

	// GetSleeptime 获取sleep时长，单位:ms，为0则表示不需要降速
	GetSleeptime(ctx context.Context, key string) (int32, error)

	// BatchGetSleeptime 批量获取sleep时长
	BatchGetSleeptime(ctx context.Context, keys []string) (map[string]int32, error)
}

// 反压中心 压力计算 http接口

type PressureHttpClient struct{}

func (c *PressureHttpClient) GetSleeptime(ctx context.Context, key string) (int32, error) {
	resp, err := c.BatchGetSleeptime(ctx, []string{key})
	if err != nil {
		return 0, err
	}
	if len(resp) == 0 {
		return 0, nil
	}
	return resp[key], nil
}

type BatchQueryPressureSignalReq struct {
	SignalList []string `json:"signal_list"`
}

type BatchQueryPressureSignalResp struct {
	Data struct {
		PressureSignalMap map[string]int32 `json:"pressure_signal_map"`
	} `json:"data"`
}

func (c *PressureHttpClient) BatchGetSleeptime(ctx context.Context, keys []string) (map[string]int32, error) {
	fmt.Printf("[pressure decelerator] keys: %+v\n", keys)
	req := BatchQueryPressureSignalReq{
		SignalList: keys,
	}
	path := strings.ReplaceAll(BatchQueryPressureSignalPath, constants.ReplaceNamespace, utils.GetNamespaceFromCtx(ctx))
	ctx = withPressureSdkReqTag(ctx)
	body, _, err := GetOpenapiClient().PostJson(ctx, path, nil, &req, AppTokenMiddleware, TenantAndUserMiddleware, ServiceIDMiddleware)
	if err != nil {
		fmt.Printf("BatchQueryPressureSignal PostJson error : %+v", err)
		return nil, err
	}
	var resp BatchQueryPressureSignalResp
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	fmt.Printf("BatchQueryPressureSignal resp.PressureSignalMap : %+v\n", resp.Data.PressureSignalMap)
	return resp.Data.PressureSignalMap, nil
}

// 带上反压中心请求tag，防止死循环
func withPressureSdkReqTag(ctx context.Context) context.Context {
	return context.WithValue(ctx, constants.CtxKeyPressureReqTag, true)
}

func checkPressureSdkReqTag(ctx context.Context) bool {
	val := ctx.Value(constants.CtxKeyPressureReqTag)
	if val == nil {
		return false
	}
	return true
}

type MockPressureHttpClient struct{}

func (c *MockPressureHttpClient) GetSleeptime(ctx context.Context, key string) (int32, error) {
	return 1000, nil
}

func (c *MockPressureHttpClient) BatchGetSleeptime(ctx context.Context, keys []string) (map[string]int32, error) {
	data := make(map[string]int32, len(keys))
	for _, key := range keys {
		data[key] = 1000
	}
	return data, nil
}
