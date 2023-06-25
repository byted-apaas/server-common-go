// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package structs

import (
	"time"
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
	Flow Flow `json:"flow"`
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
