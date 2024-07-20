// Package version defines version of server-common-go.
package version

import (
	"sync"
)

// Version is server-common-go version.
const Version = "v0.0.30-beta.0" // v0.0.28 版本号已被使用，本次版本为 v0.0.29

const SDKName = "byted-apaas/server-common-go"

type ISDKInfo interface {
	GetVersion() string
	GetSDKName() string
}

type CommonSDKInfo struct{}

func (c *CommonSDKInfo) GetVersion() string {
	return Version
}

func (c *CommonSDKInfo) GetSDKName() string {
	return SDKName
}

var (
	commonSDKInfoOnce sync.Once
	commonSDKInfo     ISDKInfo
)

func GetCommonSDKInfo() ISDKInfo {
	if commonSDKInfo == nil {
		commonSDKInfoOnce.Do(func() {
			commonSDKInfo = &CommonSDKInfo{}
		})
	}
	return commonSDKInfo
}
