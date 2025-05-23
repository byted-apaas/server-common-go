// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package constants

import "time"

const (
	HttpClientDialTimeoutDefault = 2 * time.Second
	HttpClientTLSTimeoutDefault  = 1 * time.Second

	RpcClientConnectTimeoutDefault = 3 * time.Second
	RpcClientRWTimeoutDefault      = 20 * time.Minute
	APITimeoutDefault              = 12 * time.Second

	AppTokenRefreshRemainTime = 20 * 60 * 1000 // 20m
)

const (
	GetAppToken         = "openapi_getAppToken"
	GetAllGlobalConfigs = "openapi_getAllGlobalConfigs"

	CreateRecord                 = "openapi_createRecord"
	CreateRecordV2               = "openapi_createRecordV2"
	CreateRecordV3               = "openapi_createRecordV3"
	BatchCreateRecord            = "openapi_batchCreateRecord"
	BatchCreateRecordV2          = "openapi_batchCreateRecordV2"
	BatchCreateRecordV3          = "openapi_batchCreateRecordV3"
	BatchCreateRecordAsync       = "openapi_batchCreateRecordAsync"
	UpdateRecord                 = "openapi_updateRecord"
	UpdateRecordV2               = "openapi_updateRecordV2"
	UpdateRecordV3               = "openapi_updateRecordV3"
	BatchUpdateRecord            = "openapi_batchUpdateRecord"
	BatchUpdateRecordV2          = "openapi_batchUpdateRecordV2"
	BatchUpdateRecordV3          = "openapi_batchUpdateRecordV3"
	BatchUpdateRecordAsync       = "openapi_batchUpdateRecordAsync"
	DeleteRecord                 = "openapi_deleteRecord"
	DeleteRecordV2               = "openapi_deleteRecordV2"
	DeleteRecordV3               = "openapi_deleteRecordV3"
	BatchDeleteRecord            = "openapi_batchDeleteRecord"
	BatchDeleteRecordV2          = "openapi_batchDeleteRecordV2"
	BatchDeleteRecordV3          = "openapi_batchDeleteRecordV3"
	BatchDeleteRecordAsync       = "openapi_batchDeleteRecordAsync"
	GetRecords                   = "openapi_getRecords"
	GetRecordsV2                 = "openapi_getRecordsV2"
	GetRecordV3                  = "openapi_getRecordV3"
	GetRecordsV3                 = "openapi_GetRecordsV3"
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

	CreateMessage                          = "openapi_createMessage"
	UpdateMessage                          = "openapi_updateMessage"
	MGetUserSettings                       = "openapi_mGetUserSettings"
	WorkflowUpdateVariables                = "openapi_workflowUpdateVariables"
	TerminateWorkflowInstance              = "openapi_terminateWorkflowInstance"
	GetFields                              = "openapi_getFieldsV5"
	GetField                               = "openapi_getFieldV5"
	GetUIMetaRecordDetail                  = "openapi_getUIMetaRecordDetail"
	GetUIMetaList                          = "openapi_getUIMetaList"
	GetMobileList                          = "openapi_getMobileList"
	UpdateMobileList                       = "openapi_updateMobileList"
	GetIdByApiName                         = "openapi_getIdByApiName"
	GetBatchAttachmentToken                = "openapi_getBatchAttachmentToken"
	GetDefaultIntegrationAppAccessToken    = "openapi_getDefaultIntegrationAppAccessToken"
	GetIntegrationAppAccessToken           = "openapi_getIntegrationAppAccessToken"
	GetDefaultIntegrationTenantAccessToken = "openapi_getDefaultIntegrationTenantAccessToken"
	GetIntegrationTenantAccessToken        = "openapi_getIntegrationTenantAccessToken"
	GetApprovalInstanceList                = "openapi_getApprovalInstanceList"
	GetApprovalInstance                    = "openapi_getApprovalInstance"

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
	FaaSInfraPSM    string // Deprecated
	InnerAPIPSM     string // Deprecated
	BOE             string // Deprecated
}

const (
	EnvTypeDev     string = "developmentboe" // Deprecated
	EnvTypeStaging string = "stagingboe"
	EnvTypeLr      string = "staging"
	EnvTypeGray    string = "gray"
	EnvTypeOnline  string = "online"

	EnvTypeStagingI18n string = "stagingboei18n"
	EnvTypeSG          string = "onlinesg"
	EnvTypeMY          string = "onlinemy"
)

var (
	// EnvConfMap 配置全是外网域名，实际消费不到，优先从 env 中消费
	EnvConfMap = map[string]PlatformConf{
		EnvTypeStaging: {
			OpenAPIDomain:   "",
			InnerAPIDomain:  "",
			FaaSInfraDomain: "",
			InnerAPIPSM:     "",
			FaaSInfraPSM:    "",
			BOE:             "boe",
		},
		EnvTypeLr: {
			OpenAPIDomain:   "https://oapi-kunlun-staging.bytedance.com",
			InnerAPIDomain:  "https://apaas-innerapi-lr.feishu-pre.cn",
			FaaSInfraDomain: "https://apaas-faasinfra-staging.bytedance.com",
			InnerAPIPSM:     "",
			FaaSInfraPSM:    "",
		},
		EnvTypeGray: {
			OpenAPIDomain:   "https://oapi-kunlun-gray.kundou.cn",
			InnerAPIDomain:  "https://apaas-innerapi.feishu-pre.cn",
			FaaSInfraDomain: "https://apaas-faasinfra-gray.kundou.cn",
			InnerAPIPSM:     "",
			FaaSInfraPSM:    "",
		},
		EnvTypeOnline: {
			OpenAPIDomain:   "https://oapi-kunlun.kundou.cn",
			InnerAPIDomain:  "https://apaas-innerapi.feishu.cn",
			FaaSInfraDomain: "https://apaas-faasinfra.kundou.cn",
			InnerAPIPSM:     "",
			FaaSInfraPSM:    "",
		},
		EnvTypeStagingI18n: {
			OpenAPIDomain:   "",
			InnerAPIDomain:  "",
			FaaSInfraDomain: "",
		},
		EnvTypeMY: {
			OpenAPIDomain:   "https://oapi-kunlun-my.byteintl.net",
			InnerAPIDomain:  "https://apaas-innerapi-my.byteintl.net",
			FaaSInfraDomain: "https://apaas-faasinfra-my.byteintl.net",
		},
		EnvTypeSG: {
			OpenAPIDomain:   "https://oapi-kunlun-my.byteintl.net",
			InnerAPIDomain:  "https://apaas-innerapi-my.byteintl.net",
			FaaSInfraDomain: "https://apaas-faasinfra-my.byteintl.net",
		},
	}
)

const (
	DefaultMeshDestReqTimeout = 60000
)
