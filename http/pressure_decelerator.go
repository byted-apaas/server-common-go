package http

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// 反压中心 - 降速

// PressureDecelerator 反压中心降速缓存结构体
type PressureDecelerator struct {
	config atomic.Value // *PressureConfig
	ctx    atomic.Value // context.Context
	client IPressureHttpClient

	cache sync.Map // map[string]*PressureDeceleratorItem
	size  int64    // current cache size

	ticker   *time.Ticker
	updating int32 // ticker update task mutex
}

type PressureDeceleratorItem struct {
	first       sync.Once // first time load
	key         string
	sleeptime   int32 // unit: ms
	lastReqTime int64 // last request time, unit: ms
}

func NewPressureDecelerator(ctx context.Context, config *PressureConfig, client IPressureHttpClient) *PressureDecelerator {
	if config == nil {
		config = &defaultPressureConfig
	}

	pd := &PressureDecelerator{
		client: client,
		size:   0,
	}
	pd.setConfig(config)
	pd.setContext(ctx)
	//if pd.config.UpdateInterval <= 0 {
	//	panic("pressure decelerator config update interval must be a positive number")
	//}
	pd.ticker = time.NewTicker(time.Duration(config.UpdateInterval) * time.Millisecond)

	return pd
}

func (pd *PressureDecelerator) getContext() context.Context {
	return pd.ctx.Load().(context.Context)
}

func (pd *PressureDecelerator) setContext(ctx context.Context) {
	pd.ctx.Store(ctx)
}

func (pd *PressureDecelerator) getConfig() *PressureConfig {
	return pd.config.Load().(*PressureConfig) // pd.config need to be not nil
}

func (pd *PressureDecelerator) setConfig(config *PressureConfig) {
	pd.config.Store(config)
}

func (pd *PressureDecelerator) GetSleeptime(key string) int32 {

	if key == "" { // if pd == nil || key == "" { return 0 }
		return 0
	}

	value, ok := pd.cache.Load(key)
	if !ok {
		var loaded bool
		if value, loaded = pd.cache.LoadOrStore(key, &PressureDeceleratorItem{
			key:         key,
			lastReqTime: getCurrentTimestampMs(),
		}); !loaded {
			atomic.AddInt64(&pd.size, 1)
		}
	}

	item := value.(*PressureDeceleratorItem)
	fmt.Println("GetSleeptime start key: ", key) // 2
	item.first.Do(func() {
		fmt.Println("GetSleeptime inner start key: ", key) // 1
		pd.updateOne(item)
		fmt.Println("GetSleeptime inner end key: ", key) // 0
	})
	fmt.Println("GetSleeptime end key: ", key) // 0
	atomic.StoreInt64(&item.lastReqTime, getCurrentTimestampMs())

	return atomic.LoadInt32(&item.sleeptime)
}

func (pd *PressureDecelerator) updateOne(item *PressureDeceleratorItem) {

	st, err := pd.client.GetSleeptime(pd.getContext(), item.key)
	if err != nil {
		fmt.Println("PressureDecelerator http client GetSleeptime error : ", err.Error())
	}
	if maxSleeptime := pd.getConfig().MaxSleeptime; maxSleeptime > 0 { // max_sleeptime < 0 表示不限制，max_sleeptime > 0，最大为 max_sleeptime
		st = minInt32(st, int32(maxSleeptime))
	}

	atomic.StoreInt32(&item.sleeptime, st)
}

func (pd *PressureDecelerator) RunUpdateTask() {
	// pd.config.UpdateInterval > 0 && pd.ticker != nil
	for range pd.ticker.C {
		go func() {
			if atomic.CompareAndSwapInt32(&pd.updating, 0, 1) { // 定时器更新任务间互斥
				defer atomic.StoreInt32(&pd.updating, 0) // 解锁
				now := getCurrentTimestampMs()
				updateKeys := make([]string, 0, atomic.LoadInt64(&pd.size)+20) // 多设置20个预留，可能中途有新增的
				sortKeys := make(map[string]int64, atomic.LoadInt64(&pd.size))
				evictKeys := make([]string, 0, atomic.LoadInt64(&pd.size)) // 还是记录淘汰key列表，使用updateKeys取反会把中途新增的新key也淘汰掉
				pd.cache.Range(func(key, value interface{}) bool {
					item := value.(*PressureDeceleratorItem)
					if lastReqTime := atomic.LoadInt64(&item.lastReqTime); now-lastReqTime <= pd.getConfig().EvictThreshold {
						updateKeys = append(updateKeys, item.key)
						sortKeys[item.key] = lastReqTime
					} else {
						evictKeys = append(evictKeys, item.key)
					}
					return true
				})
				pd.evictKeys(evictKeys)
				if len(updateKeys) == 0 { // 更新key列表为空，则直接返回
					return
				}
				sort.Slice(updateKeys, func(i, j int) bool { // 按last_req_time降序排序
					return sortKeys[updateKeys[i]] > sortKeys[updateKeys[j]]
				})
				if maxKeyCap := pd.getConfig().MaxKeyCapacity; len(updateKeys) > maxKeyCap { // 仍超过最大容量，淘汰最早的
					pd.evictKeys(updateKeys[maxKeyCap:])
					updateKeys = updateKeys[:maxKeyCap]
				}
				res, err := pd.client.BatchGetSleeptime(pd.getContext(), updateKeys)
				if err != nil {
					fmt.Printf("PressureDecelerator update ticker [%d] error : %+v\n", now, err)
					//res = make(map[string]int64) // nil map 不可写，但可读，因此不做处理
				}
				for _, key := range updateKeys {
					value, ok := pd.cache.Load(key)
					if !ok || value == nil {
						continue
					}
					item := value.(*PressureDeceleratorItem)
					if maxSleeptime := pd.getConfig().MaxSleeptime; maxSleeptime > 0 {
						atomic.StoreInt32(&item.sleeptime, minInt32(res[key], int32(maxSleeptime)))
					} else {
						atomic.StoreInt32(&item.sleeptime, res[key])
					}
				}
			}
		}()
	}
}

// evictKeys 淘汰keys
// 当前淘汰策略：1.last_req_time超出阈值的key；2.当前缓存超过最大容量，淘汰last_req_time最小的key
func (pd *PressureDecelerator) evictKeys(keys []string) {
	//atomic.AddInt64(&pd.size, -int64(len(keys)))
	for _, key := range keys {
		pd.cache.Delete(key)
		atomic.AddInt64(&pd.size, -1)
	}
}

func (pd *PressureDecelerator) StopUpdateTask() {
	pd.ticker.Stop()
}

func getCurrentTimestampMs() int64 {
	// time.Now().UnixNano() / int64(time.Millisecond)
	return time.Now().UnixNano() / 1e6
}

func minInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
