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

func GetFaaSInfraPSMFromEnv() (psm string, cluster string) {
	return os.Getenv(constants.EnvKFaaSInfraPSM), "default"
}

func GetInnerAPIPSMFromEnv() string {
	return os.Getenv(constants.EnvKInnerAPIPSM)
}

func GetLGWPSMFromEnv() string {
	return os.Getenv(constants.EnvKLGWPSM)
}

func GetLGWClusterFromEnv() string {
	return os.Getenv(constants.EnvKLGWCluster)
}

func GetIfPrintRequestCurl() bool {
	return os.Getenv(constants.EnvPrintRequest) == "true"
}
