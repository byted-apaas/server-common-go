// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package utils

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/byted-apaas/server-common-go/structs"
)

const (
	red    = 31
	green  = 32
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct {
}

func (m *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	level, color := LevelInfo(entry.Level)
	msg := fmt.Sprintf("%s %s %s", level, entry.Time.Format("15:04:05.999"), entry.Message)
	if entry.Level == logrus.DebugLevel {
		msg = fmt.Sprintf("%s %s %s", entry.Time.Format("15:04:05.999"), level, entry.Message)
	}
	b.WriteString(fmt.Sprintf("\x1b[%dm%s\x1b[0m\n", color, msg))
	return b.Bytes(), nil
}

const (
	// Log Type
	Console = "Console"
	System  = "System"
	Result  = "Result"

	WebIDESource = "apaas-ide"
	ISWebIDE     = "IS_WEB_IDE"
)

type WebIDELogFormatter struct {
}

func (m *WebIDELogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 构建日志结构体
	log := structs.WebIDELog{
		Source:  WebIDESource,
		Time:    time.Now(),
		Type:    Console,
		Level:   WebIDELevelInfo(entry.Level),
		Message: entry.Message,
	}

	if entry.Level == logrus.DebugLevel {
		log.Type = Result
	}

	msg, err := JsonMarshalBytes(log)
	if err != nil {
		return nil, err
	}

	b.WriteString(fmt.Sprintf("%s\n", msg))
	return b.Bytes(), nil
}

func WebIDELevelInfo(level logrus.Level) string {
	switch level {
	case logrus.InfoLevel:
		return "info"
	case logrus.ErrorLevel:
		return "error"
	case logrus.WarnLevel:
		return "warn"
	default:
		return "info"
	}
}

func LevelInfo(level logrus.Level) (string, int) {
	switch level {
	case logrus.InfoLevel:
		return "[INFO] ", blue
	case logrus.WarnLevel:
		return "[WARN] ", yellow
	case logrus.ErrorLevel:
		return "[ERROR]", red
	case logrus.DebugLevel:
		return "Result ", green
	}

	return "[NULL] ", gray
}

type ConsoleLogger struct {
	logID string
}

func GetConsoleLogger(logIDs ...string) *ConsoleLogger {
	if os.Getenv(ISWebIDE) == "true" {
		logrus.SetFormatter(&WebIDELogFormatter{})
	} else {
		logrus.SetFormatter(&LogFormatter{})
	}
	logrus.SetLevel(logrus.TraceLevel)

	l := &ConsoleLogger{}
	if len(logIDs) > 0 && logIDs[0] != "" {
		l.logID = logIDs[0]
	}
	return l
}

func (c *ConsoleLogger) Infof(format string, args ...interface{}) {
	if c.logID != "" {
		format = fmt.Sprintf("%s %s", c.logID, format)
	}
	logrus.Infof(format, args...)
}

func (c *ConsoleLogger) Warnf(format string, args ...interface{}) {
	if c.logID != "" {
		format = fmt.Sprintf("%s %s", c.logID, format)
	}
	logrus.Warnf(format, args...)
}

func (c *ConsoleLogger) Errorf(format string, args ...interface{}) {
	if c.logID != "" {
		format = fmt.Sprintf("%s %s", c.logID, format)
	}
	logrus.Errorf(format, args...)
}

func (c *ConsoleLogger) Result(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}
