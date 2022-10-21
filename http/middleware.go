// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package http

import (
	"context"
	"net/http"

	"github.com/byted-apaas/server-common-go/constants"
	"github.com/byted-apaas/server-common-go/utils"
)

type ReqMiddleWare func(ctx context.Context, req *http.Request) error

func AppTokenMiddleware(ctx context.Context, req *http.Request) (err error) {
	if req == nil || req.Header == nil {
		return nil
	}

	credential := getCredentialFromCtx(ctx)
	if credential == nil {
		credential, err = getFaaSCredential()
		if err != nil {
			return err
		}
	}

	token, err := credential.getToken(ctx)
	if err != nil {
		return err
	}

	req.Header.Add(constants.HttpHeaderKeyAuthorization, token)
	return nil
}

func TenantAndUserMiddleware(ctx context.Context, req *http.Request) error {
	if req == nil || req.Header == nil {
		return nil
	}
	req.Header.Add(constants.HttpHeaderKeyTenant, utils.GetTenantName())
	req.Header.Add(constants.HttpHeaderKeyUser, "-1")
	return nil
}

func ServiceIDMiddleware(ctx context.Context, req *http.Request) error {
	if req == nil || req.Header == nil {
		return nil
	}
	req.Header.Add(constants.HttpHeaderKeyServiceID, utils.GetServiceID())
	return nil
}
