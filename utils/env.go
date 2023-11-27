package utils

import (
	"os"

	"github.com/byted-apaas/server-common-go/constants"
)

// GetGlobalAuthType 全局鉴权类型配置，优先级低于接口级
func GetGlobalAuthType() *string {
	authType := os.Getenv(constants.GlobalAuthTypeKey)
	if authType == "" {
		return nil
	}
	return &authType
}

func GetFaaSPlatform() string {
	return os.Getenv(constants.EnvKFaaSType)
}

func GetOpenAPIDomainName() string {
	return os.Getenv("KOpenApiDomain")
}

func GetFaaSInfraDomainName() string {
	return os.Getenv("KOpenApiDomain")
}
