// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package logger

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/byted-apaas/server-common-go/constants"
	"github.com/byted-apaas/server-common-go/http"
	"github.com/byted-apaas/server-common-go/structs"
	"github.com/byted-apaas/server-common-go/utils"
)

const (
	LogDomain = "lowcode_func_log"

	NormalLog      = 1
	AggregationLog = 2
)

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type I18nTag struct {
	Key   string       `json:"key"`
	Value structs.I18n `json:"value"`
}

type ExtraInfo struct {
	FunctionVersionID int64        `json:"functionVersionID"`
	SourceLabel       structs.I18n `json:"sourceLabel"`
	ObjectLabel       structs.I18n `json:"objectLabel"`
	TriggerTimeCost   int64        `json:"triggerTimeCost"`
	RuntimeCost       int64        `json:"runtimeCost"`
}

type Log struct {
	Domain          string `json:"domain"`
	Type            int    `json:"type"`
	Level           int    `json:"level"`
	CreateTime      int64  `json:"createTime"` // 毫秒级时间戳
	RequestID       string `json:"RequestID"`
	Sequence        int64  `json:"sequence"`
	Content         string `json:"content"`
	CreateTimeMicro *int64 `json:"createTimeMicro"` // 微秒级时间戳

	Tags      []Tag     `json:"tags"`
	TagsI18n  []I18nTag `json:"tagsI18n"`
	ExtraInfo ExtraInfo `json:"extraInfo"`
}

type Logger struct {
	RequestID        string
	executeID        string
	functionAPIID    string
	tenantID         int64
	namespace        string
	tenantType       int64
	lock             sync.Mutex
	logs             []string
	tags             []Tag
	tagsI18n         []I18nTag
	extraInfo        ExtraInfo
	startTriggerTime int64
	startRuntime     int64
	errorNum         int64
	infoNum          int64
	warnNum          int64
	sequence         int64
	isDebug          bool
	streamLogCount   int64
}

func NewLogger(ctx context.Context) *Logger {
	l := &Logger{
		lock:             sync.Mutex{},
		logs:             make([]string, 0),
		RequestID:        utils.GetLogIDFromCtx(ctx),
		executeID:        utils.GetExecuteIDFromCtx(ctx),
		functionAPIID:    utils.GetFunctionAPIIDFromCtx(ctx),
		tenantID:         utils.GetTenantIDFromCtx(ctx),
		namespace:        utils.GetNamespaceFromCtx(ctx),
		tenantType:       utils.GetTenantTypeFromCtx(ctx),
		startTriggerTime: getFunctionLoggerExtraToCtx(ctx).StartTriggerTime,

		startRuntime:   time.Now().UnixNano() / int64(time.Millisecond),
		errorNum:       0,
		infoNum:        0,
		warnNum:        0,
		isDebug:        utils.GetDebugTypeFromCtx(ctx) != 0,
		sequence:       1,
		streamLogCount: 0,
	}

	if !l.isDebug {
		l.tags = l.getTags(ctx)
		l.tagsI18n = l.getTagsI18(ctx)
		l.extraInfo = l.getExtraInfo(ctx)
	}

	return l
}

func NewConsoleLogger(ctx context.Context) *Logger {
	return &Logger{
		RequestID: utils.GetLogIDFromCtx(ctx),
		isDebug:   true,
	}
}

func SetLogger(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, constants.CtxKeyLogger, l)
}

func GetLogger(ctx context.Context) *Logger {
	l, ok := ctx.Value(constants.CtxKeyLogger).(*Logger)
	if !ok || l == nil {
		return NewConsoleLogger(ctx)
		// utils.GetConsoleLogger().Errorf("GetLogger failed !")
		// panic("[Logger Usage Error] please make sure that your context parameter in GetLogger() method inherits from the functions Handler, rather than self-built context or an empty context.")
	}

	return l
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if !l.isDebug {
		atomic.AddInt64(&l.infoNum, 1)
		l.addLog(fmt.Sprintf(format, args...), utils.LogLevelInfo, NormalLog)
		if l.streamLogCount < utils.LogCountLimit {
			l.streamLogCount++
			content := fmt.Sprintf("%s %s %s %s", utils.GetFormatDate(), constants.APaaSLogPrefix, l.getFormatLog(utils.LogLevelInfo, format, args...), constants.APaaSLogSuffix)
			fmt.Println(content)
		}
	} else {
		utils.GetConsoleLogger().Infof(format, args...)
	}
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if !l.isDebug {
		atomic.AddInt64(&l.warnNum, 1)
		l.addLog(fmt.Sprintf(format, args...), utils.LogLevelWarn, NormalLog)
		if l.streamLogCount < utils.LogCountLimit {
			l.streamLogCount++
			content := fmt.Sprintf("%s %s %s %s", utils.GetFormatDate(), constants.APaaSLogPrefix, l.getFormatLog(utils.LogLevelWarn, format, args...), constants.APaaSLogSuffix)
			fmt.Println(content)
		}
	} else {
		utils.GetConsoleLogger().Warnf(format, args...)
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if !l.isDebug {
		atomic.AddInt64(&l.errorNum, 1)
		l.addLog(fmt.Sprintf(format, args...), utils.LogLevelError, NormalLog)
		if l.streamLogCount < utils.LogCountLimit {
			l.streamLogCount++
			content := fmt.Sprintf("%s %s %s %s", utils.GetFormatDate(), constants.APaaSLogPrefix, l.getFormatLog(utils.LogLevelError, format, args...), constants.APaaSLogSuffix)
			fmt.Println(content)
		}
	} else {
		utils.GetConsoleLogger().Errorf(format, args...)
	}
}

func (l *Logger) getFormatLog(level int, format string, args ...interface{}) string {
	formatLog := utils.FormatLog{
		Level:         level,
		EventID:       l.executeID,
		FunctionAPIID: l.functionAPIID,
		LogID:         l.RequestID,
		Timestamp:     time.Now().UnixNano() / 1e3, // 使用微秒
		Message:       fmt.Sprintf(format, args...),
		TenantID:      l.tenantID,
		TenantType:    l.tenantType,
		Namespace:     l.namespace,
		LogType:       constants.UserLogType,
	}
	return utils.GetFormatLogWithMessage(formatLog, l.streamLogCount)
}

func Send(ctx context.Context, l *Logger) {
	if l.isDebug {
		return
	}
	var (
		err         error
		data        []byte
		compressLog string
	)
	defer func() {
		if err != nil {
			utils.GetConsoleLogger(l.RequestID).Errorf("[Logger] Send failed, err: %v", err)
		}
	}()

	if len(l.logs) == 0 {
		return
	}
	l.addLog("", utils.LogLevelInfo, AggregationLog)

	data, err = utils.JsonMarshalBytes(l.logs)
	if err != nil {
		return
	}
	compressLog, err = CompressForDeflate(data)
	if err != nil {
		return
	}

	err = http.SendLog(ctx, map[string]string{"compressData": compressLog})
}

func (l *Logger) addLog(content string, level int, logType int) {
	if logType == NormalLog && len(l.logs) >= utils.LogCountLimit {
		return
	}

	if len(content) > utils.LogLengthLimit {
		content = content[:utils.LogLengthLimit] + utils.LogLengthLimitTip
	}

	log := Log{
		Domain:          LogDomain,
		RequestID:       l.RequestID,
		Type:            logType,
		Level:           level,
		CreateTime:      TimeNowMils(),
		CreateTimeMicro: TimeNowMicros(), // 用于旧日志转发到可观测
		Sequence:        l.getSequence(),
		Content:         content,
		Tags:            make([]Tag, 0),
		TagsI18n:        make([]I18nTag, 0),
		ExtraInfo:       ExtraInfo{},
	}

	// 聚合日志
	if logType == AggregationLog {
		curTime := TimeNowMils()
		log.Tags = l.tags
		log.TagsI18n = l.tagsI18n
		log.ExtraInfo = l.extraInfo
		log = l.tagsAddNum(log)
		if l.errorNum > 0 {
			log.Level = utils.LogLevelError
		}
		log.ExtraInfo.TriggerTimeCost = curTime - l.startTriggerTime
		log.ExtraInfo.RuntimeCost = curTime - l.startRuntime
		if log.ExtraInfo.TriggerTimeCost <= 0 {
			log.ExtraInfo.TriggerTimeCost = log.ExtraInfo.RuntimeCost
		}
		b, _ := utils.JsonMarshalBytes(log)
		l.lock.Lock()
		l.logs = append(l.logs, string(b))
		l.lock.Unlock()
		return
	}

	if len(l.logs) < utils.LogCountLimit {
		if len(l.logs) == utils.LogCountLimit-1 {
			log.Content = utils.LogCountLimitTip
		}
		b, _ := utils.JsonMarshalBytes(log)
		l.lock.Lock()
		l.logs = append(l.logs, string(b))
		l.lock.Unlock()
	}
}

func (l *Logger) getTags(ctx context.Context) []Tag {
	return []Tag{
		{
			Key:   "tenantID",
			Value: strconv.FormatInt(utils.GetTenantIDFromCtx(ctx), 10),
		}, {
			Key:   "namespace",
			Value: utils.GetNamespaceFromCtx(ctx),
		}, {
			Key:   "sourceID",
			Value: getFunctionLoggerExtraToCtx(ctx).SourceID,
		}, {
			Key:   "triggerType",
			Value: utils.GetTriggerTypeFromCtx(ctx),
		}, {
			Key:   "functionName",
			Value: utils.GetFunctionNameFromCtx(ctx),
		}, {
			Key:   "source",
			Value: strconv.Itoa(utils.GetSourceTypeFromCtx(ctx)),
		}, {
			Key:   "instanceID",
			Value: strconv.FormatInt(getFunctionLoggerExtraToCtx(ctx).InstanceID, 10),
		}, {
			Key:   "tenantType",
			Value: strconv.FormatInt(utils.GetTenantTypeFromCtx(ctx), 10),
		}, {
			Key:   "functionAPIID",
			Value: utils.GetFunctionAPIIDFromCtx(ctx),
		},
	}
}

func (l *Logger) getTagsI18(ctx context.Context) []I18nTag {
	return []I18nTag{
		{
			Key:   "functionLabel",
			Value: getFunctionLoggerExtraToCtx(ctx).FunctionLabel,
		},
	}
}

func (l *Logger) getExtraInfo(ctx context.Context) ExtraInfo {
	return ExtraInfo{
		FunctionVersionID: getFunctionLoggerExtraToCtx(ctx).FunctionVersionID,
		SourceLabel:       getFunctionLoggerExtraToCtx(ctx).SourceLabel,
		ObjectLabel:       getFunctionLoggerExtraToCtx(ctx).ObjectLabel,
	}
}

func (l *Logger) getSequence() int64 {
	return atomic.AddInt64(&l.sequence, 1)
}

func (l *Logger) tagsAddNum(log Log) Log {
	log.Tags = append(log.Tags, Tag{
		Key:   "infoNum",
		Value: strconv.FormatInt(l.infoNum, 10),
	})

	log.Tags = append(log.Tags, Tag{
		Key:   "warnNum",
		Value: strconv.FormatInt(l.warnNum, 10),
	})

	log.Tags = append(log.Tags, Tag{
		Key:   "errorNum",
		Value: strconv.FormatInt(l.errorNum, 10),
	})

	return log
}

func TimeNowMils() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func TimeNowMicros() *int64 {
	t := time.Now().UnixNano() / int64(time.Microsecond)
	return &t
}

func CompressForDeflate(b []byte) (string, error) {
	var buf bytes.Buffer

	w := zlib.NewWriter(&buf)
	if _, e := w.Write(b); e != nil {
		return "", e
	}
	w.Close()

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// FunctionLoggerExtra Only Logger 使用的参数
type FunctionLoggerExtra struct {
	FunctionLabel     structs.I18n `json:"functionLabel"`
	SourceID          string       `json:"sourceID"`
	SourceLabel       structs.I18n `json:"sourceLabel"`
	FunctionVersionID int64        `json:"functionVersionID"`
	ObjectLabel       structs.I18n `json:"objectLabel"`
	InstanceID        int64        `json:"instanceID"`
	StartTriggerTime  int64        `json:"startTriggerTime"`
}

func SetFunctionLoggerExtraToCtx(ctx context.Context, extra FunctionLoggerExtra) context.Context {
	return context.WithValue(ctx, constants.CtxKeyFLoggerExtra, extra)
}

func getFunctionLoggerExtraToCtx(ctx context.Context) FunctionLoggerExtra {
	extra := FunctionLoggerExtra{}
	v := ctx.Value(constants.CtxKeyFLoggerExtra)
	_ = utils.Decode(v, &extra)
	return extra
}
