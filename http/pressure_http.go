package http

import "context"

// IPressureHttpClient 反压中心 http client
type IPressureHttpClient interface {

	// GetSleeptime 获取sleep时长，单位:ms，为0则表示不需要降速
	GetSleeptime(ctx context.Context, key string) (int64, error)

	// BatchGetSleeptime 批量获取sleep时长
	BatchGetSleeptime(ctx context.Context, keys []string) (map[string]int64, error)
}

// 反压中心 压力计算 http接口

type PressureHttpClient struct{}

func (c *PressureHttpClient) GetSleeptime(ctx context.Context, key string) (int64, error) {
	// todo unimplemented
	return 0, nil
}

func (c *PressureHttpClient) BatchGetSleeptime(ctx context.Context, keys []string) (map[string]int64, error) {
	// todo unimplemented
	return map[string]int64{}, nil
}

type MockPressureHttpClient struct{}

func (c *MockPressureHttpClient) GetSleeptime(ctx context.Context, key string) (int64, error) {
	return 1000, nil
}

func (c *MockPressureHttpClient) BatchGetSleeptime(ctx context.Context, keys []string) (map[string]int64, error) {
	data := make(map[string]int64, len(keys))
	for _, key := range keys {
		data[key] = 1000
	}
	return data, nil
}
