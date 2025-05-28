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
	EnvKBizIDC          = "KBizIDC"
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
	ReplaceObjectAPINameV3 = ":object_api_name"
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

	ExecuteID                   = "x-serverless-execute-id"
	FunctionAPIID               = "x-serverless-function-api-id"
	HTTPHeaderEnvoyRespFlag     = "x-envoy-response-flags"
	PodRateLimitQuotaHeader     = "x-serverless-sdk-pod-rate-limit-quota"
	PodRateLimitDowngradeHeader = "x-serverless-sdk-pod-rate-limit-downgrade"

	PressureNeedDecelerateHeader = "x-serverless-sdk-pressure-need-decelerate" // 反压中心是否需要降速，由CloudFunction下发该开关
	PressureConfigHeader         = "x-serverless-sdk-pressure-config"          // 反压中心相关配置，由CloudFunction下发该配置
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

	// PersistFaaSKeyRequestSource 溯源流量key
	PersistFaaSKeyRequestSource = APAAS_PERSIST_FAAS_PREFIX + "request-source"

	// RequestSourcePressureSignalId  RequestSource 中的 PressureSignalId，用于反压中心识别异步链路
	RequestSourcePressureSignalId = "pressureSignalId"

	// RequestSourceIsAsync  RequestSource 中的 IsAsync，用于反压中心标识是否异步链路
	RequestSourceIsAsync = "isAsync"
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

const (
	APaaSLogPrefix   = "apaas-log-prefix"
	APaaSLogSuffix   = "apaas-log-suffix"
	UserLogType      = "user"
	RateLimitLogType = "rate_limit" // SDK 限流
	SpeedDownLogType = "speed_down" // SDK 降速
)
