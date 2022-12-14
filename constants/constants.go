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
