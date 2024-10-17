package http

import (
	"context"
	"strconv"

	"github.com/byted-apaas/server-common-go/constants"
	cExceptions "github.com/byted-apaas/server-common-go/exceptions"
	"github.com/byted-apaas/server-common-go/structs"
	"github.com/byted-apaas/server-common-go/utils"
)

func GetAppTokenHttp(ctx context.Context, clientID, clientSecret string) (*structs.AppTokenResp, error) {
	data := map[string]interface{}{
		"clientId":       clientID,
		"clientSecret":   clientSecret,
		"withTenantInfo": true,
	}

	body, err := utils.ErrorWrapper(GetOpenapiClient().PostJson(ctx, OpenapiPathGetToken, nil, data))
	if err != nil {
		return nil, err
	}

	tokenResult := structs.AppTokenResp{}
	if err = utils.JsonUnmarshalBytes(body, &tokenResult); err != nil {
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

	data, err := utils.ErrorWrapper(body, extra, err)
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
