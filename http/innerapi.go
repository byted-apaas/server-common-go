// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package http

//func getClientOptions() []client.Option {
//	return []client.Option{
//		client.WithRPCTimeout(constants.RpcClientRWTimeoutDefault),
//		client.WithConnectTimeout(constants.RpcClientConnectTimeoutDefault),
//		client.WithCluster("default"),
//		client.WithLongConnection(connpool.IdleConfig{
//			MaxIdlePerAddress: 10,
//			MaxIdleGlobal:     1000,
//			MaxIdleTimeout:    60 * time.Second,
//		})}
//}
//
//var (
//	innerAPICliMap        sync.Map
//	innerAPICliMapForAuth sync.Map
//)

//func GetInnerAPICli(ctx context.Context, opts ...client.Option) (client innerapiservice.Client, err error) {
//	psm := utils.GetInnerAPIPSM(ctx)
//	boe := utils.GetBoe(ctx)
//	cli, ok := innerAPICliMap.Load(psm + boe)
//	if !ok {
//		client, err = innerapiservice.NewClient(psm, append(opts, getClientOptions()...)...)
//		if err != nil {
//			return nil, err
//		}
//		innerAPICliMap.Store(psm+boe, client)
//		return client, nil
//	}
//	client, ok = cli.(innerapiservice.Client)
//	if !ok {
//		return nil, cExceptions.InternalError("client assert failed, cli: %v", cli)
//	}
//	return client, nil
//}

//func GetInnerAPICliForAuth(ctx context.Context, opts ...client.Option) (client innerapiservice.Client, err error) {
//	psm := utils.GetInnerAPIPSM(ctx)
//	boe := utils.GetBoe(ctx)
//	cli, ok := innerAPICliMapForAuth.Load(psm + boe)
//	if !ok {
//		client, err = innerapiservice.NewClient(psm, append(opts, getClientOptions()...)...)
//		if err != nil {
//			return nil, err
//		}
//		innerAPICliMapForAuth.Store(psm+boe, client)
//		return client, nil
//	}
//	client, ok = cli.(innerapiservice.Client)
//	if !ok {
//		return nil, cExceptions.InternalError("client for auth assert failed, cli: %v", cli)
//	}
//	return client, nil
//}

//func GetAppTokenInnerAPI(ctx context.Context, clientID, clientSecret string) (result *identity.GetAppTokenResponse, err error) {
//	req := identity.NewGetAppTokenRequest()
//	req.ClientID = clientID
//	req.ClientSecret = clientSecret
//	req.WithTenantInfo = utils.BoolPtr(true)
//
//	cli, err := GetInnerAPICliForAuth(ctx)
//	if err != nil {
//		return nil, cExceptions.ErrWrap(err)
//	}
//
//	ctx = utils.SetKEnvToCtxForRPC(ctx)
//	ctx = utils.SetAPaaSLaneToCtxForRPC(ctx)
//
//	var cancel context.CancelFunc
//	ctx, cancel = GetTimeoutCtx(ctx)
//	defer cancel()
//	resp, err := cli.GetAppToken(ctx, req)
//	if err != nil {
//		return nil, cExceptions.InternalError("call innerAPI GetAppToken() failed: %+v", err)
//	}
//
//	if resp.BaseResp.KStatusCode != "" {
//		msg := resp.BaseResp.KStatusMessage
//		if resp.BaseResp.StatusMessage != "" {
//			msg = resp.BaseResp.StatusMessage
//		}
//		return nil, cExceptions.NewErrWithCodeV2(resp.BaseResp.KStatusCode, msg, utils.GetLogIDFromCtx(ctx))
//	}
//
//	return resp, nil
//}

//func GetAppTokenRpc(ctx context.Context, clientID, clientSecret string) (result *structs.AppTokenResp, err error) {
//	rpcRes, err := GetAppTokenInnerAPI(ctx, clientID, clientSecret)
//	if err != nil {
//		return nil, err
//	}
//
//	var tInfo = &common.TenantInfo{}
//	if rpcRes.TenantInfo != nil {
//		tInfo = rpcRes.TenantInfo
//	}
//	var outsideDomainName string
//	if tInfo.OutsideTenantInfo != nil {
//		outsideDomainName = tInfo.OutsideTenantInfo.OutsideDomainName
//	}
//	result = &structs.AppTokenResp{
//		AccessToken: rpcRes.AccessToken,
//		ExpireTime:  rpcRes.ExpireTime,
//		Namespace:   rpcRes.Namespace,
//		TenantInfo: structs.TenantInfo{
//			ID:         tInfo.TenantID,
//			DomainName: tInfo.DomainName,
//			TenantName: tInfo.TenantName,
//			TenantType: int64(tInfo.TenantType),
//			OutsideTenantInfo: struct {
//				OutsideDomainName string `json:"outsideDomainName"`
//			}{OutsideDomainName: outsideDomainName},
//		},
//	}
//
//	return result, nil
//}

//func GetFunctionMetaRpc(ctx context.Context, apiName string) (funcMeta *structs.FunctionMeta, err error) {
//	if ctx == nil {
//		ctx = context.Background()
//	}
//	ctx, err = RebuildRpcCtx(utils.SetApiTimeoutMethodToCtx(ctx, constants.GetFunction))
//	if err != nil {
//		return nil, err
//	}
//
//	var cancel context.CancelFunc
//	ctx, cancel = GetTimeoutCtx(ctx)
//	defer cancel()
//
//	req := cloudfunction.NewGetFunctionRequest()
//	ctx = metainfo.WithValue(ctx, constants.HttpHeaderKeyUser, strconv.FormatInt(utils.GetUserIDFromCtx(ctx), 10))
//	req.Namespace = utils.GetNamespace()
//	req.APIName = &apiName
//	cli, err := GetInnerAPICli(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	resp, err := cli.GetFunction(ctx, req)
//	if err != nil {
//		return nil, cExceptions.InternalError("call innerAPI GetAppToken() failed: %+v", err)
//	}
//	if resp.BaseResp.KStatusCode != "" {
//		msg := resp.BaseResp.KStatusMessage
//		if resp.BaseResp.StatusMessage != "" {
//			msg = resp.BaseResp.StatusMessage
//		}
//		return nil, cExceptions.NewErrWithCodeV2(resp.BaseResp.KStatusCode, msg, utils.GetLogIDFromCtx(ctx))
//	}
//
//	if resp.FunctionDetail == nil {
//		return nil, nil
//	}
//
//	funcMeta = &structs.FunctionMeta{ApiName: resp.FunctionDetail.GetAPIName()}
//	for _, input := range resp.FunctionDetail.Input {
//		v := structs.IOParamItem{
//			Key:  input.Key,
//			Type: input.Type,
//		}
//		if input.ObjectApiName != nil {
//			v.ObjectAPIName = *input.ObjectApiName
//		}
//		funcMeta.IOParam.Input = append(funcMeta.IOParam.Input, &v)
//	}
//
//	for _, output := range resp.FunctionDetail.Output {
//		v := structs.IOParamItem{
//			Key:  output.Key,
//			Type: output.Type,
//		}
//		if output.ObjectApiName != nil {
//			v.ObjectAPIName = *output.ObjectApiName
//		}
//		funcMeta.IOParam.Output = append(funcMeta.IOParam.Output, &v)
//	}
//	return funcMeta, nil
//}
