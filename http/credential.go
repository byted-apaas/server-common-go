// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package http

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/byted-apaas/server-common-go/constants"
	exp "github.com/byted-apaas/server-common-go/exceptions"
	"github.com/byted-apaas/server-common-go/structs"
	"github.com/byted-apaas/server-common-go/utils"
)

type ICredential interface {
	getToken(ctx context.Context) (string, error)
	setSystemFlag(ctx context.Context, isSystem bool)
}

type AppCredential struct {
	id, secret string

	tenantInfo atomic.Value // TenantInfo
	token      atomic.Value // string
	expireTime atomic.Value // int64

	lock     sync.Mutex
	isSystem bool
}

func NewAppCredential(id, secret string) *AppCredential {
	return &AppCredential{
		id:     id,
		secret: secret,
		lock:   sync.Mutex{},
	}
}

func (c *AppCredential) GetID() string {
	return c.id
}

func (c *AppCredential) getToken(ctx context.Context) (string, error) {
	expireTime, ok := c.expireTime.Load().(int64)
	if ok && expireTime-utils.NowMils() > constants.AppTokenRefreshRemainTime {
		token, ok := c.token.Load().(string)
		if ok {
			return token, nil
		}
	}

	token, _, err := c.refresh(ctx)
	return token, err
}

func (c *AppCredential) setSystemFlag(ctx context.Context, isSystem bool) {
	c.isSystem = isSystem
}

func (c *AppCredential) GetTenantInfo(ctx context.Context) (*structs.Tenant, error) {
	if c == nil {
		credential, err := getFaaSCredential()
		if err != nil {
			return nil, exp.InternalError("get system credential failed: " + err.Error())
		}
		c = credential
	}
	tenant, ok := c.tenantInfo.Load().(*structs.Tenant)
	if ok {
		return tenant, nil
	}
	ctx = withPressureSdkReqTag(ctx)
	_, tenantInfo, err := c.refresh(ctx)
	return tenantInfo, err
}

func (c *AppCredential) refresh(ctx context.Context) (token string, tenantInfo *structs.Tenant, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	resp, err := c.fetchToken(ctx)
	if err != nil {
		return "", nil, err
	}
	c.token.Store(resp.AccessToken)
	c.expireTime.Store(resp.ExpireTime)
	tenant, ok := c.tenantInfo.Load().(*structs.Tenant)
	if !ok {
		tenantInfo = &structs.Tenant{
			ID:        resp.TenantInfo.ID,
			Name:      resp.TenantInfo.DomainName,
			Type:      resp.TenantInfo.TenantType,
			Namespace: resp.Namespace,
			Domain:    resp.TenantInfo.OutsideTenantInfo.OutsideDomainName,
		}
		c.tenantInfo.Store(tenantInfo)
	} else {
		tenantInfo = tenant
	}
	return resp.AccessToken, tenantInfo, nil
}

func (c *AppCredential) fetchToken(ctx context.Context) (result *structs.AppTokenResp, err error) {

	ctx = utils.SetApiTimeoutMethodToCtx(ctx, constants.GetAppToken)
	result, err = GetAppTokenHttp(ctx, c.id, c.secret)

	if err != nil {
		return nil, err
	}

	return result, nil
}

var (
	faaSTokenInstance *AppCredential
)

func getFaaSCredential() (*AppCredential, error) {
	if faaSTokenInstance != nil {
		return faaSTokenInstance, nil
	}

	appID, appSecret, err := utils.GetAppIDAndSecret()
	if err != nil {
		return nil, err
	}

	faaSTokenInstance = NewAppCredential(appID, appSecret)
	faaSTokenInstance.setSystemFlag(context.Background(), true)

	return faaSTokenInstance, nil
}

func getCredentialFromCtx(ctx context.Context) *AppCredential {
	if ctx == nil {
		return nil
	}
	tokenInstance, _ := ctx.Value(constants.CtxKeyCredential).(*AppCredential)
	return tokenInstance
}

func SetCredentialToCtx(ctx context.Context, credential *AppCredential) context.Context {
	if credential == nil {
		return ctx
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, constants.CtxKeyCredential, credential)
}
