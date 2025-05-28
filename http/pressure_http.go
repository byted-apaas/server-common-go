package http

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/byted-apaas/server-common-go/utils"
)

const (
	BatchQueryPressureSignalPath = "/batch_query"
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
	SignalList []string `json:"signalList"`
	TenantId   int64    `json:"TenantId"`
	Namespace  string   `json:"Namespace"`
}

type BatchQueryPressureSignalResp struct {
	PressureSignalMap map[string]int32 `json:"PressureSignalMap"`
}

func (c *PressureHttpClient) BatchGetSleeptime(ctx context.Context, keys []string) (map[string]int32, error) {
	tenant, err := utils.GetTenantFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if tenant == nil { // 兜底
		return nil, errors.New("get tenant info from ctx is nil")
	}
	req := &BatchQueryPressureSignalReq{
		SignalList: keys,
		TenantId:   tenant.ID,
		Namespace:  tenant.Namespace,
	}
	respByte, _, err := GetPressureSdkClient().PostJson(ctx, BatchQueryPressureSignalPath, nil, &req)
	if err != nil {
		return nil, err
	}
	var resp BatchQueryPressureSignalResp
	if err = json.Unmarshal(respByte, &resp); err != nil {
		return nil, err
	}
	return resp.PressureSignalMap, nil
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
