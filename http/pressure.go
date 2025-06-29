package http

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/byted-apaas/server-common-go/utils"
)

// PressureConfig 反压中心配置，由CloudFunction下发
type PressureConfig struct {
	MaxSleeptime   int64 `yaml:"MaxSleeptime" json:"MaxSleeptime"`     // 最大sleeptime，单位：ms，-1表示不限制，0转换默认值
	UpdateInterval int64 `yaml:"UpdateInterval" json:"UpdateInterval"` // 定时器更新周期，单位：ms，需要 > 0
	MaxKeyCapacity int   `yaml:"MaxKeyCapacity" json:"MaxKeyCapacity"` // 最大key容量，需要匹配反压中心http接口批量最大容量，超过会对key进行淘汰，需要 > 0
	EvictThreshold int64 `yaml:"EvictThreshold" json:"EvictThreshold"` // 淘汰阈值，当某个key未请求超过该阈值，则触发淘汰，单位：ms，需要 > 0
}

const (
	DefaultPressureMaxSleeptime   int64 = 1000          // 1s
	DefaultPressureUpdateInterval int64 = 5000          // 5s
	DefaultPressureMaxKeyCapacity int   = 1000          // 默认1000个
	DefaultPressureEvictThreshold int64 = 1 * 60 * 1000 // 1min
)

var (
	defaultPressureConfig = PressureConfig{
		MaxSleeptime:   DefaultPressureMaxSleeptime,
		UpdateInterval: DefaultPressureUpdateInterval,
		MaxKeyCapacity: DefaultPressureMaxKeyCapacity,
		EvictThreshold: DefaultPressureEvictThreshold,
	}
)

// NewPressureConfigByJsonStr new PressureConfig, replace illegal field values with default values
func NewPressureConfigByJsonStr(data string) *PressureConfig {
	if data == "" {
		return &defaultPressureConfig
	}

	var conf PressureConfig
	if err := json.Unmarshal([]byte(data), &conf); err != nil {
		fmt.Println("NewPressureConfigByJsonStr json unmarshal error : ", err.Error())
		conf = defaultPressureConfig
	}

	if conf.MaxSleeptime == 0 { // max_sleeptime < 0 表示不限制，max_sleeptime > 0，最大为 max_sleeptime
		conf.MaxSleeptime = DefaultPressureMaxSleeptime
	}
	if conf.UpdateInterval <= 0 {
		conf.UpdateInterval = DefaultPressureUpdateInterval
	}
	if conf.MaxKeyCapacity <= 0 {
		conf.MaxKeyCapacity = DefaultPressureMaxKeyCapacity
	}
	if conf.EvictThreshold <= 0 {
		conf.EvictThreshold = DefaultPressureEvictThreshold
	}
	return &conf
}

var (
	pressureDecelerator     *PressureDecelerator
	pressureDeceleratorOnce sync.Once
)

// InitPressureDecelerator 初始化 pressureDecelerator
func InitPressureDecelerator(ctx context.Context) {
	config := NewPressureConfigByJsonStr(utils.GetPressureConfigFromCtx(ctx))
	pressureDeceleratorOnce.Do(func() {
		//client := &MockPressureHttpClient{}
		client := &PressureHttpClient{}
		pressureDecelerator = NewPressureDecelerator(ctx, config, client)
		go func() {
			pressureDecelerator.RunUpdateTask() // 启动刷新任务
		}()
	})
	UpdatePressureConfig(config)
	UpdatePressureContext(ctx)
}

func UpdatePressureConfig(config *PressureConfig) {
	pressureDecelerator.setConfig(config)
}

func UpdatePressureContext(ctx context.Context) {
	pressureDecelerator.setContext(ctx)
}
