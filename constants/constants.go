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
	EnvKFaaSScene       = "KFaaSScene"
	EnvKFaaSType        = "KFaaSType"
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

	HttpHeaderKeyOrgID      = "X-Kunlun-Org-Id"
	HttpHeaderKeySDKFuncMsg = "Rpc-Persist-Kunlun-Faassdk"
	HttpHeaderKeyEnv        = "x-tt-env"
	HttpHeaderKeyUsePPE     = "x-use-ppe"
	HttpHeaderKeyUseBOE     = "x-use-boe"
	HttpHeaderKeyAPaaSLane  = "Rpc-Persist-Lane-C-Apaas-Lane"

	HTTPHeaderKeyAuthType = "Rpc-Persist-AUTH-TYPE"
	AuthTypeKey           = "AUTH_TYPE"
	AuthTypeSystem        = "system"
	AuthTypeUser          = "user"
	AuthTypeMixUserSystem = "mix_user_system" // OQL 场景需要，传该值时只有 select 过权限，where、orderBy 和 groupBy 都不影响
)

const (
	HttpHeaderValueJson = "application/json"
	HttpHeaderValueBson = "application/bson"
)

const (
	MetaInfoFaaSSdkFuncMsgKey = "KUNLUN_FAASSDK"
)

const (
	DebugTypeInvoke = 0
	DebugTypeOnline = 1
	DebugTypeLocal  = 2
)

const (
	APAAS_PERSIST_FAAS_PREFIX = "x-apaas-persist-faas-"
	PersistFaaSKeySummarized  = "x-apaas-persist-faas-summarized"

	PersistFaaSKeyFaaSType = "x-apaas-persist-faas-type"
)

const (
	FaaSTypeFunction     = "function"
	FaaSTypeMicroService = "microService"
	FaaSTypeOpenSDK      = "openSDK"
)
