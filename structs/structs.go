// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package structs

import (
	rawJson "encoding/json"
	"strconv"
	"time"

	"github.com/byted-apaas/server-common-go/structs/common/i18n"
	"github.com/byted-apaas/server-common-go/structs/common/reference"
)

type I18n []map[string]interface{}

type HttpConfig struct {
	Domain             string
	MaxIdleConn        int
	MaxIdleConnPerHost int
	IdleConnTimeout    time.Duration
}

type AppTokenResult struct {
	Code string       `json:"code"`
	Msg  string       `json:"msg"`
	Data AppTokenResp `json:"data"`
}

type TenantInfo struct {
	ID                int64  `json:"id"`
	DomainName        string `json:"domainName"`
	TenantName        string `json:"tenantName"`
	TenantType        int64  `json:"tenantType"`
	OutsideTenantInfo struct {
		OutsideDomainName string `json:"outsideDomainName"`
	} `json:"outsideTenantInfo"`
}

type AppInfo struct {
	Namespace   string                      `json:"namespace"`
	Label       i18n.I18nCnUs               `json:"label"`
	Description i18n.I18nCnUs               `json:"description"`
	CreatedAt   int64                       `json:"createdAt"`
	CreatedBy   *reference.LookupWithAvatar `json:"createdBy"`
}

type EventInfo struct {
	Type       string        `json:"type"`
	Name       i18n.I18nCnUs `json:"name"`
	ApiName    string        `json:"apiName"`
	InstanceId int64         `json:"instanceId"`
}

type AppTokenResp struct {
	AccessToken string     `json:"accessToken"`
	ExpireTime  int64      `json:"expireTime"`
	Namespace   string     `json:"namespace"`
	TenantInfo  TenantInfo `json:"tenantInfo"`
}

type RPCCliConf struct {
	Psm         string        `yaml:"Psm" json:"Psm"`
	DebugAddr   string        `yaml:"DebugAddr" json:"DebugAddr"`
	Cluster     string        `yaml:"Cluster" json:"Cluster"`
	IDC         string        `yaml:"IDC" json:"IDC"`
	Timeout     time.Duration `yaml:"Timeout" json:"Timeout"`
	ConnTimeout time.Duration `yaml:"ConnTimeout" json:"ConnTimeout"`
}

// UserContext 上下文参数
type UserContext struct {
	Flow       Flow `json:"flow"`
	Permission struct {
		UnauthFields map[string]interface{} `json:"_unauthFields"`
	} `json:"permission"`
}

type Flow struct {
	Execution FlowExecution `json:"execution"`
	APIName   string        `json:"apiName"`
	// Deprecated: 已废弃
	Variables map[string]CfVariable `json:"variables"`
}

// Deprecated: 已废弃 CfVariable 流程变量
type CfVariable struct {
	Value     interface{} `json:"value"`
	FieldType string      `json:"type"`
	VarType   string      `json:"varType"`
}

// FlowExecution 流程实例相关
type FlowExecution struct {
	// 流程实例 ID
	ID int64 `json:"id"`
}

type WebIDELog struct {
	Source  string    `json:"source"`
	Time    time.Time `json:"time"`
	Type    string    `json:"type"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
}
type Permission struct {
	UnauthFields map[string]interface{} `json:"_unauthFields"`
}

type RecordOnlyID struct {
	ID interface{} `json:"_id"`
}

func (r RecordOnlyID) GetID() (id int64) {
	switch r.ID.(type) {
	case int64:
		id, _ = r.ID.(int64)
	case string:
		idStr, _ := r.ID.(string)
		id, _ = strconv.ParseInt(idStr, 10, 64)
	case rawJson.Number:
		isNumber, _ := r.ID.(rawJson.Number)
		id, _ = isNumber.Int64()
	}
	return id
}

type ParamUnauthField struct {
	Type             string     `json:"type"`
	UnauthFields     []string   `json:"unauthFields"`
	UnauthFieldsList [][]string `json:"unauthFieldsList"`
}
