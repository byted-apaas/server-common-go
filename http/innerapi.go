// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package http

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/byted-apaas/server-common-go/constants"
	cExceptions "github.com/byted-apaas/server-common-go/exceptions"
	"github.com/byted-apaas/server-common-go/structs"
	"github.com/byted-apaas/server-common-go/utils"
	"github.com/tidwall/gjson"
)

const (
	SCFileDownload = ""
	SCSuccess      = "0"

	ECInternalError  = "k_ec_000001"
	ECNoTenantID     = "k_ec_000002"
	ECNoUserID       = "k_ec_000003"
	ECUnknownError   = "k_ec_000004"
	ECOpUnknownError = "k_op_ec_00001"
	ECSystemBusy     = "k_op_ec_20001"
	ECSystemError    = "k_op_ec_20002"
	ECRateLimitError = "k_op_ec_20003"
	ECTokenExpire    = "k_ident_013000"
	ECIllegalToken   = "k_ident_013001"
	ECMissingToken   = "k_op_ec_10205"
)

func errorWrapper(body []byte, extra map[string]interface{}, err error) ([]byte, error) {
	if err != nil {
		return nil, cExceptions.ErrWrap(err)
	}

	code := gjson.GetBytes(body, "code").String()
	msg := gjson.GetBytes(body, "msg").String()
	switch code {
	case SCFileDownload:
		return body, nil
	case SCSuccess:
		data := gjson.GetBytes(body, "data")
		if data.Type == gjson.String {
			return []byte(data.Str), nil
		}
		return []byte(data.Raw), nil
	case ECInternalError, ECNoTenantID, ECNoUserID, ECUnknownError,
		ECOpUnknownError, ECSystemBusy, ECSystemError, ECRateLimitError,
		ECTokenExpire, ECIllegalToken, ECMissingToken:
		return nil, cExceptions.InternalError("%v ([%v] %v)", msg, code, utils.GetLogIDFromExtra(extra))
	default:
		return nil, cExceptions.InvalidParamError("%v ([%v] %v)", msg, code, utils.GetLogIDFromExtra(extra))
	}
}

func GetAppTokenHttp(ctx context.Context, clientID, clientSecret string) (*structs.AppTokenResp, error) {
	data := map[string]interface{}{
		"clientId":       clientID,
		"clientSecret":   clientSecret,
		"withTenantInfo": true,
	}

	body, err := errorWrapper(GetOpenapiClient().PostJson(ctx, OpenapiPathGetToken, nil, data))
	if err != nil {
		return nil, err
	}

	tokenResult := structs.AppTokenResp{}
	if err = json.Unmarshal(body, &tokenResult); err != nil {
		return nil, cExceptions.InternalError("[AppCredential] fetchToken Unmarshal TokenResult failed, err: %v", err)
	}

	return &tokenResult, nil
}

func GetFunctionMetaHttp(ctx context.Context, apiName string) (funcMeta *structs.FunctionMeta, err error) {
	headers := map[string][]string{
		constants.HttpHeaderKeyUser: {strconv.FormatInt(utils.GetUserIDFromCtx(ctx), 10)},
	}

	body, extra, err := GetOpenapiClient().PostJson(ctx, GetInnerAPIPathGetFunction(), headers, map[string]interface{}{"apiName": apiName}, AppTokenMiddleware, TenantAndUserMiddleware, ServiceIDMiddleware)
	if err != nil {
		return nil, cExceptions.ErrWrap(err)
	}

	data, err := errorWrapper(body, extra, err)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Detail *struct {
			APIID   string                 `json:"apiID"`
			APIName string                 `json:"apiName"`
			Input   []*structs.IOParamItem `json:"input"`
			Output  []*structs.IOParamItem `json:"output"`
		} `json:"detail"`
	}

	logid := utils.GetLogIDFromExtra(extra)
	if err := utils.JsonUnmarshalBytes(data, &resp); err != nil {
		return nil, cExceptions.InternalError("InvokeFunctionWithAuth failed, err: %v, logid: %v", err, logid)
	}

	if resp.Detail == nil {
		return nil, nil
	}

	return &structs.FunctionMeta{
		ApiName: resp.Detail.APIName,
		IOParam: structs.FunctionIOParam{
			Input:  resp.Detail.Input,
			Output: resp.Detail.Output,
		},
	}, nil
}
