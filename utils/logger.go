package utils

import (
	"time"
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

func GetFormatDate() string {
	return time.Now().Format("2006-01-02")
}

func GetFormatLogWithMessage(formatLog FormatLog, streamLogCount int64) string {
	if len(formatLog.Message) > LogLengthLimit {
		formatLog.Message = formatLog.Message[:LogLengthLimit] + LogLengthLimitTip
	}
	if streamLogCount == LogCountLimit {
		formatLog.Message = formatLog.Message + LogCountLimitTip
	}

	jsonContent, err := JsonMarshalBytes(formatLog)
	if err != nil {
		GetConsoleLogger(formatLog.LogID).Errorf("[Logger] getFormatLog failed, err: %v", err)
	}

	return string(jsonContent)
}
