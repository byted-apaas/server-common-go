package http

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPressureDecelerator(t *testing.T) {
	ctx := context.Background()
	conf := &PressureConfig{
		MaxSleeptime:   500,
		MaxKeyCapacity: 3,
		EvictThreshold: 2000,
		UpdateInterval: 1000,
	}
	cli := &MockPressureHttpClient{}
	pd := NewPressureDecelerator(ctx, conf, cli)
	go func() { pd.RunUpdateTask() }()
	defer pd.StopUpdateTask()
	assert.Equal(t, int32(500), pd.GetSleeptime("key1"))
	time.Sleep(200 * time.Millisecond)
	_, ok := pd.cache.Load("key1")
	assert.True(t, ok)
	time.Sleep(time.Duration(conf.EvictThreshold+conf.UpdateInterval) * time.Millisecond)
	_, ok = pd.cache.Load("key1")
	assert.False(t, ok)
	pd.GetSleeptime("key1")
	time.Sleep(10 * time.Millisecond)
	pd.GetSleeptime("key2")
	time.Sleep(10 * time.Millisecond)
	pd.GetSleeptime("key3")
	pd.GetSleeptime("key4")
	time.Sleep(time.Duration(conf.UpdateInterval) * time.Millisecond)
	_, ok = pd.cache.Load("key1")
	assert.False(t, ok)
	_, ok = pd.cache.Load("key2")
	assert.True(t, ok)
}

/**
goos: darwin
goarch: arm64
pkg: github.com/byted-apaas/server-common-go/http
BenchmarkPressureDecelerator
PressureDecelerator ticker update task started at 1748921442889
PressureDecelerator ticker update task started at 1748921443389
BenchmarkPressureDecelerator-10    	19787380	        58.91 ns/op
PASS
*/
func BenchmarkPressureDecelerator(b *testing.B) {
	ctx := context.Background()
	conf := &PressureConfig{
		MaxSleeptime:   500,
		MaxKeyCapacity: 5,
		EvictThreshold: 2000,
		UpdateInterval: 500,
	}
	cli := &MockPressureHttpClient{}
	pd := NewPressureDecelerator(ctx, conf, cli)
	go func() { pd.RunUpdateTask() }()
	defer pd.StopUpdateTask()
	keys := []string{"key1", "key2", "key3", "key4"}
	for i := 0; i < b.N; i++ {
		pd.GetSleeptime(keys[i%4])
	}
}

/**
goos: darwin
goarch: arm64
pkg: github.com/byted-apaas/server-common-go/http
BenchmarkPressureDeceleratorWithEvictKeys
PressureDecelerator ticker update task started at 1748921584988
PressureDecelerator ticker update task started at 1748921585488
BenchmarkPressureDeceleratorWithEvictKeys-10    	18219337	        62.66 ns/op
PASS
*/
func BenchmarkPressureDeceleratorWithEvictKeys(b *testing.B) {
	ctx := context.Background()
	conf := &PressureConfig{
		MaxSleeptime:   500,
		MaxKeyCapacity: 5,
		EvictThreshold: 2000,
		UpdateInterval: 500,
	}
	cli := &MockPressureHttpClient{}
	pd := NewPressureDecelerator(ctx, conf, cli)
	go func() { pd.RunUpdateTask() }()
	defer pd.StopUpdateTask()
	keys := []string{"key1", "key2", "key3", "key4", "key5", "key6", "key7", "key8"}
	for i := 0; i < b.N; i++ {
		pd.GetSleeptime(keys[i%8])
	}
}
