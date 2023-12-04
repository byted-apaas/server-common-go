// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package constants

const (
	EnvKEnvironment     = "ENV"
	EnvKOrgID           = "KOrgID"
	EnvKTenantName      = "KTenantName"
	EnvKNamespace       = "KNamespace"
	EnvKSvcID           = "KSvcID"
	EnvKClientID        = "KClientID"
	EnvKClientSecret    = "KClientSecret"
	EnvKOpenApiDomain   = "KOpenApiDomain"
	EnvKFaaSInfraDomain = "KFaaSInfraDomain"
	EnvKInnerAPIDomain  = "KInnerAPIDomain"
	EnvKFaaSInfraPSM    = "KFaaSInfraPSM"
	EnvKInnerAPIPSM     = "KInnerAPIPSM"
	EnvKLGWPSM          = "KLGWPSM"
	EnvKLGWCluster      = "KLGWCluster"
	EnvKFaaSScene       = "KFaaSScene"
	EnvKFaaSType        = "KFaaSType"
	EnvKMeshHttp        = "TCE_ENABLE_HTTP_SIDECAR_EGRESS"
	EnvKMeshUDS         = "WITH_HTTP_MESH_EGRESS_UDS"
	EnvKSocketAddr      = "SERVICE_MESH_HTTP_EGRESS_ADDR"
)

const (
	ReplaceNamespace       = ":namespace"
	ReplaceObjectAPIName   = ":objectAPIName"
	ReplaceFieldAPIName    = ":fieldAPIName"
	ReplaceFunctionAPIName = ":functionAPIName"
	ReplaceRecordID        = ":recordID"
	ReplaceFileID          = ":fileID"
	ReplaceExecutionID     = ":executionId"
	ReplaceAPIName         = ":apiName"
)

const (
	HttpHeaderKeyTenant        = "Tenant"
	HttpHeaderKeyUser          = "User"
	HttpHeaderKeyServiceID     = "X-Kunlun-Service-Id"
	HttpHeaderKeyAuthorization = "Authorization"
	HttpHeaderKeyContentType   = "Content-Type"
	HttpHeaderKeyLogID         = "X-Tt-Logid"

	HttpHeaderKeyOrgID       = "X-Kunlun-Org-Id"
	HttpHeaderKeySDKFuncMsg  = "Rpc-Persist-Kunlun-Faassdk"
	HttpHeaderKeyEnv         = "x-tt-env"
	HttpHeaderKeyUsePPE      = "x-use-ppe"
	HttpHeaderKeyUseBOE      = "x-use-boe"
	HTTPHeaderKeyAPaaSLane   = "Rpc-Persist-Lane-C-Apaas-Lane"
	HTTPHeaderKeyFaaSLaneID  = "X-Ae-Lane"
	HTTPHeaderKeyFaaSEnvID   = "X-Ae-EnvID"
	HTTPHeaderKeyFaaSEnvType = "X-Ae-EnvType"

	HTTPHeaderKeyAuthType = "Rpc-Persist-AUTH-TYPE"
	AuthTypeKey           = "AUTH_TYPE"
	GlobalAuthTypeKey     = "GlobalAuthType"
	AuthTypeSystem        = "system"
	AuthTypeUser          = "user"
	AuthTypeMixUserSystem = "mix_user_system" // OQL 场景需要，传该值时只有 select 过权限，where、orderBy 和 groupBy 都不影响

	HTTPHeaderKeyTLBEnv = "X-Tlb-Env"     // TLB 分流标签
	TLBEnvOAPILGWGray   = "oapi_lgw_gray" // 请求至 openapi 域名后，TLB 分流至 LGW
)

const (
	HttpHeaderValueJson = "application/json"
	HttpHeaderValueBson = "application/bson"
)

const (
	MetaInfoFaaSSdkFuncMsgKey = "KUNLUN_FAASSDK"
	MetaInfoAPaaSLaneKey      = "LANE_C_APAAS_LANE"
)

const (
	DebugTypeInvoke = 0
	DebugTypeOnline = 1
	DebugTypeLocal  = 2
)

const (
	APAAS_PERSIST_PREFIX      = "rpc-persist-"
	PersistAPaaSKeySummarized = "rpc-persist-apaas-summarized"
	APAAS_PERSIST_FAAS_PREFIX = "x-apaas-persist-faas-"
	PersistFaaSKeySummarized  = "x-apaas-persist-faas-summarized"

	PersistFaaSKeyFaaSType = "x-apaas-persist-faas-type"

	PersistFaaSKeyOpenAPIRoutingType = "x-apaas-persist-faas-openapi-routing-type"

	PersistFaaSKeyFromSDKName    = "x-apaas-persist-faas-from-sdk-name"
	PersistFaaSKeyFromSDKVersion = "x-apaas-persist-faas-from-sdk-version"
)

const (
	OpenAPIRoutingTypeToLGW      = "toLGW"
	OpenAPIRoutingTypeToInnerAPI = "toInnerAPI"
)

const (
	FaaSTypeFunction     = "function"
	FaaSTypeMicroService = "microService"
	FaaSTypeOpenSDK      = "openSDK"
)

const (
	FunctionMetaConfCacheTableKey = "function-meta-conf-cache-table"
)
