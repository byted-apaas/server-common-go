// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package constants

import "time"

const (
	HttpClientDialTimeoutDefault = 2 * time.Second
	HttpClientTLSTimeoutDefault  = 1 * time.Second

	RpcClientConnectTimeoutDefault = 1 * time.Second
	RpcClientRWTimeoutDefault      = 20 * time.Minute
	APITimeoutDefault              = 12 * time.Second

	AppTokenRefreshRemainTime = 20 * 60 * 1000 // 20m
)

const (
	GetAppToken         = "openapi_getAppToken"
	GetAllGlobalConfigs = "openapi_getAllGlobalConfigs"

	CreateRecord                 = "openapi_createRecord"
	CreateRecordV2               = "openapi_createRecordV2"
	BatchCreateRecord            = "openapi_batchCreateRecord"
	BatchCreateRecordV2          = "openapi_batchCreateRecordV2"
	BatchCreateRecordAsync       = "openapi_batchCreateRecordAsync"
	UpdateRecord                 = "openapi_updateRecord"
	UpdateRecordV2               = "openapi_updateRecordV2"
	BatchUpdateRecord            = "openapi_batchUpdateRecord"
	BatchUpdateRecordV2          = "openapi_batchUpdateRecordV2"
	BatchUpdateRecordAsync       = "openapi_batchUpdateRecordAsync"
	DeleteRecord                 = "openapi_deleteRecord"
	DeleteRecordV2               = "openapi_deleteRecordV2"
	BatchDeleteRecord            = "openapi_batchDeleteRecord"
	BatchDeleteRecordV2          = "openapi_batchDeleteRecordV2"
	BatchDeleteRecordAsync       = "openapi_batchDeleteRecordAsync"
	GetRecords                   = "openapi_getRecords"
	GetRecordsV2                 = "openapi_getRecordsV2"
	ModifyRecordsWithTransaction = "openapi_modifyRecordsWithTransaction"
	Oql                          = "openapi_oql"
	GetExecutionUserTaskInfo     = "openapi_getExecutionUserTaskInfo"
	GetExecutionInfo             = "openapi_getExecutionInfo"
	RevokeExecution              = "openapi_revokeExecution"
	ExecuteFlow                  = "openapi_executeFlow"

	UploadAttachment     = "openapi_uploadAttachment"
	UploadAttachmentV2   = "openapi_uploadAttachmentV2"
	DownloadAttachment   = "openapi_downloadAttachment"
	DownloadAttachmentV2 = "openapi_downloadAttachmentV2"
	DownloadAvatar       = "openapi_downloadAvatar"
	UploadAvatar         = "openapi_uploadAvatar"

	InvokeFuncWithAuth = "openapi_invokeFuncWithAuth"
	GetFunction        = "openapi_getFunction"

	CreateMessage             = "openapi_createMessage"
	UpdateMessage             = "openapi_updateMessage"
	MGetUserSettings          = "openapi_mGetUserSettings"
	WorkflowUpdateVariables   = "openapi_workflowUpdateVariables"
	TerminateWorkflowInstance = "openapi_terminateWorkflowInstance"
	GetFields                 = "openapi_getFieldsV5"
	GetField                  = "openapi_getFieldV5"
	GetUIMetaRecordDetail     = "openapi_getUIMetaRecordDetail"
	GetUIMetaList             = "openapi_getUIMetaList"
	GetMobileList             = "openapi_getMobileList"
	UpdateMobileList          = "openapi_updateMobileList"
	GetIdByApiName            = "openapi_getIdByApiName"
	GetBatchAttachmentToken   = "openapi_getBatchAttachmentToken"

	GetServiceToken         = "faasinfra_getServiceToken"
	SendLog                 = "faasinfra_sendLog"
	InvokeFuncSync          = "faasinfra_invokeFuncSync"
	CreateAsyncTask         = "faasinfra_createAsyncTask"
	CreateDistributedTask   = "faasinfra_CreateDistributedTask"
	CreateAsyncTaskV2       = "faasinfra_createAsyncTaskV2"
	CreateDistributedTaskV2 = "faasinfra_CreateDistributedTaskV2"
	InvokeMicroserviceSync  = "faasinfra_invokeMicroserviceSync"
	InvokeMicroserviceAsync = "faasinfra_invokeMicroserviceAsync"
	RequestMongodb          = "faasinfra_requestMongodb"
	RequestFile             = "faasinfra_requestFile"
	RequestRedis            = "faasinfra_requestRedis"
)

// APITimeoutMapDefault millSeconds
var APITimeoutMapDefault = map[string]int64{
	"openapi_uploadAttachment":     50 * 1000,
	"openapi_uploadAttachmentV2":   50 * 1000,
	"openapi_downloadAttachment":   50 * 1000,
	"openapi_downloadAttachmentV2": 50 * 1000,
	"openapi_downloadAvatar":       30 * 1000,
	"openapi_uploadAvatar":         30 * 1000,
	"openapi_executeFlow":          25 * 1000,

	"openapi_invokeFuncWithAuth": 16 * 60 * 1000,

	"faasinfra_invokeFuncSync":         16 * 60 * 1000,
	"faasinfra_invokeMicroserviceSync": 16 * 60 * 1000,

	"faasinfra_requestFile": 30 * 1000,
}

type PlatformConf struct {
	OpenAPIDomain   string
	InnerAPIDomain  string
	FaaSInfraDomain string
	InnerAPIPSM     string
	BOE             string
}

const (
	EnvTypeLr     string = "staging"
	EnvTypeGray   string = "gray"
	EnvTypeOnline string = "online"
)

var (
	EnvConfMap = map[string]PlatformConf{
		EnvTypeLr:     {"", "https://apaas-innerapi-lr.feishu-pre.cn", "https://apaas-faasinfra-staging.bytedance.com", "", ""},
		EnvTypeGray:   {"", "https://apaas-innerapi.feishu-pre.cn", "https://apaas-faasinfra-gray.kundou.cn", "", ""},
		EnvTypeOnline: {"", "https://apaas-innerapi.feishu.cn", "https://apaas-faasinfra.kundou.cn", "", ""},
	}
)
