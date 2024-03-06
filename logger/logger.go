// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package logger

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
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

	LogLevelError = 4
	LogLevelWarn  = 5
	LogLevelInfo  = 6

	LogCountLimit     = 10000
	LogLengthLimit    = 10000
	LogLengthLimitTip = `\n... The log has been truncated because it exceeds the length limit.`
	LogCountLimitTip  = `The log has been discarded because it exceeded the limit of  10000`
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
	Domain     string         `json:"domain"`
	Type       int            `json:"type"`
	Level      int            `json:"level"`
	CreateTime int64          `json:"createTime"`
	RequestID  string         `json:"RequestID"`
	Sequence   int64          `json:"sequence"`
	Content    string         `json:"content"`
	ServiceID  string         `json:"serviceID"`
	EventID    string         `json:"eventID"`
	CtxInfo    logContextInfo `json:"ctxInfo"`

	Tags      []Tag     `json:"tags"`
	TagsI18n  []I18nTag `json:"tagsI18n"`
	ExtraInfo ExtraInfo `json:"extraInfo"`
}

type Logger struct {
	RequestID              string
	ctxInfo                logContextInfo
	lock                   sync.Mutex
	logs                   []string
	tags                   []Tag
	tagsI18n               []I18nTag
	extraInfo              ExtraInfo
	startTriggerTime       int64
	startRuntime           int64
	errorNum               int64
	infoNum                int64
	warnNum                int64
	sequence               int64
	isDebug                bool
	isLegacyLoggerDisabled bool
}

type logContextInfo struct {
	EventID      string
	ServiceID    string
	TenantID     string
	Namespace    string
	SourceID     string
	TriggerType  string
	FunctionName string
	Source       string
	InstanceID   string
}

func NewLogger(ctx context.Context) *Logger {
	l := &Logger{
		lock:             sync.Mutex{},
		logs:             make([]string, 0),
		RequestID:        utils.GetLogIDFromCtx(ctx),
		startTriggerTime: getFunctionLoggerExtraToCtx(ctx).StartTriggerTime,

		startRuntime: time.Now().UnixNano() / int64(time.Millisecond),
		errorNum:     0,
		infoNum:      0,
		warnNum:      0,
		isDebug:      utils.GetDebugTypeFromCtx(ctx) != 0,
		sequence:     1,
	}
	l.ctxInfo = l.getContextInfo(ctx)
	l.isLegacyLoggerDisabled = utils.GetLegacyLoggerDisabledFromCtx(ctx)
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
		l.addLog(fmt.Sprintf(format, args...), LogLevelInfo, NormalLog)
		fmt.Printf(format, args...)
	} else {
		utils.GetConsoleLogger().Infof(format, args...)
	}
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if !l.isDebug {
		atomic.AddInt64(&l.warnNum, 1)
		l.addLog(fmt.Sprintf(format, args...), LogLevelWarn, NormalLog)
		fmt.Printf(format, args...)
	} else {
		utils.GetConsoleLogger().Warnf(format, args...)
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if !l.isDebug {
		atomic.AddInt64(&l.errorNum, 1)
		l.addLog(fmt.Sprintf(format, args...), LogLevelError, NormalLog)
		fmt.Printf(format, args...)
	} else {
		utils.GetConsoleLogger().Errorf(format, args...)
	}
}

func Send(ctx context.Context, l *Logger) {
	if l.isLegacyLoggerDisabled {
		return
	}
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
	l.addLog("", LogLevelInfo, AggregationLog)

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

func (l *Logger) fmtStreamNormalLog(content string, level int) Log {
	log := Log{
		Domain:     LogDomain,
		RequestID:  l.RequestID,
		Level:      level,
		CreateTime: TimeNowMils(),
		Sequence:   l.getSequence(),
		Content:    content,
		Tags:       make([]Tag, 0),
		TagsI18n:   make([]I18nTag, 0),
		ExtraInfo:  ExtraInfo{},
		ServiceID:  l.ctxInfo.ServiceID,
		EventID:    l.ctxInfo.EventID,
		CtxInfo:    l.ctxInfo,
	}
	return log
}

func (l *Logger) sendStreamLog(content string, level int) {
	log := l.fmtStreamNormalLog(content, level)
	logContent, err := json.Marshal(log)
	if err != nil {
		fmt.Println("invalid content")
	}
	fmt.Println(logContent)
}

func (l *Logger) addLog(content string, level int, logType int) {
	l.sendStreamLog(content, level)
	if l.isLegacyLoggerDisabled {
		return
	}
	if logType == NormalLog && len(l.logs) >= LogCountLimit {
		return
	}

	if len(content) > LogLengthLimit {
		content = content[:LogLengthLimit] + LogLengthLimitTip
	}

	log := Log{
		Domain:     LogDomain,
		RequestID:  l.RequestID,
		Type:       logType,
		Level:      level,
		CreateTime: TimeNowMils(),
		Sequence:   l.getSequence(),
		Content:    content,
		Tags:       make([]Tag, 0),
		TagsI18n:   make([]I18nTag, 0),
		ExtraInfo:  ExtraInfo{},
	}

	// 聚合日志
	if logType == AggregationLog {
		curTime := TimeNowMils()
		log.Tags = l.tags
		log.TagsI18n = l.tagsI18n
		log.ExtraInfo = l.extraInfo
		log = l.tagsAddNum(log)
		if l.errorNum > 0 {
			log.Level = LogLevelError
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

	if len(l.logs) < LogCountLimit {
		if len(l.logs) == LogCountLimit-1 {
			log.Content = LogCountLimitTip
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
		},
	}
}

func (l *Logger) getContextInfo(ctx context.Context) logContextInfo {
	info := logContextInfo{
		EventID:      utils.GetEventID(ctx),
		ServiceID:    utils.GetServiceIDFromCtx(ctx),
		TenantID:     strconv.FormatInt(utils.GetTenantIDFromCtx(ctx), 10),
		Namespace:    utils.GetNamespaceFromCtx(ctx),
		SourceID:     getFunctionLoggerExtraToCtx().SourceID,
		TriggerType:  utils.GetTriggerTypeFromCtx(ctx),
		FunctionName: utils.GetFunctionNameFromCtx(ctx),
		Source:       strconv.Itoa(utils.GetSourceTypeFromCtx(ctx)),
		InstanceID:   strconv.FormatInt(getFunctionLoggerExtraToCtx(ctx).InstanceID, 10),
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
