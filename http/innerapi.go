// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package http

import (
	"context"
	"encoding/json"

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
