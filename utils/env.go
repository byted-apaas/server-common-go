package utils

import (
	"github.com/byted-apaas/server-common-go/constants"
	"os"
)

// GetGlobalAuthType 全局鉴权类型配置，优先级低于接口级
func GetGlobalAuthType() *string {
	authType := os.Getenv(constants.GlobalAuthTypeKey)
	if authType == "" {
		return nil
	}
	return &authType
}
