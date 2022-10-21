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

type UserContext struct {
	Flow Flow `json:"flow"`
}

type Flow struct {
	Execution FlowExecution `json:"execution"`
}

type FlowExecution struct {
	ID int64 `json:"id"`
}
