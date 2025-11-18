package utils

import (
	"context"
	"strings"
	"time"

	"github.com/byted-apaas/server-common-go/constants"
)

const (
	LogLevelError = 4
	LogLevelWarn  = 5
	LogLevelInfo  = 6

	LogCountLimit     = 10000
	LogLengthLimit    = 10000
	LogLengthLimitTip = `\n... The log has been truncated because it exceeds the length limit.`
	LogCountLimitTip  = `The log has been discarded because it exceeded the limit of  10000`
)

type FormatLog struct {
	Level         int    `json:"level"`           // 日志级别, 4-error,5-warn,6-info
	EventID       string `json:"event_id"`        // 事件 ID，可观测需要
	FunctionAPIID string `json:"function_api_id"` // 函数 API ID
	LogID         string `json:"log_id"`          // 日志 ID，事件编号与日志编号有一一对应关系
	Timestamp     int64  `json:"timestamp"`       // 时间
	Message       string `json:"message"`         // 用户的日志内容，SDK 会对超长日志截断
	TenantID      int64  `json:"tenant_id"`       // 租户 ID
	TenantType    int64  `json:"tenant_type"`     // 租户 ID
	Namespace     string `json:"namespace"`       // 命名空间
	LogType       string `json:"log_type"`        // 日志类型
}

func NewFormatLog(ctx context.Context, level int, logType, message string) *FormatLog {
	return &FormatLog{
		Level:         level,
		EventID:       GetExecuteIDFromCtx(ctx),
		FunctionAPIID: GetFunctionAPIIDFromCtx(ctx),
		LogID:         GetLogIDFromCtx(ctx),
		Timestamp:     time.Now().UnixNano() / 1e3, // 使用微秒
		TenantID:      GetTenantIDFromCtx(ctx),
		TenantType:    GetTenantTypeFromCtx(ctx),
		Namespace:     GetNamespaceFromCtx(ctx),
		LogType:       logType,
		Message:       message,
	}
}

func (l *FormatLog) String() string {
	if len(l.Message) > LogLengthLimit {
		l.Message = l.Message[:LogLengthLimit] + LogLengthLimitTip
	}

	jsonContent, err := JsonMarshalBytes(l)
	if err != nil {
		GetConsoleLogger(l.LogID).Errorf("[Logger] FormatLog String failed, err: %v", err)
	}

	var sb strings.Builder
	sb.WriteString(GetFormatDate()) // 防止日志粘连
	sb.WriteString(constants.APaaSLogPrefix)
	sb.WriteString(string(jsonContent))
	sb.WriteString(constants.APaaSLogSuffix)

	return sb.String()
}

func GetFormatDate() string {
	return time.Now().Format("2006-01-02")
}

type SpeedDownMessage struct {
	Key       string `json:"key"`
	SleepTime int32  `json:"sleep_time"` // 降速时间，单位：毫秒
}
