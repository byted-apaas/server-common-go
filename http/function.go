package http

import (
	"context"
	"fmt"
	"time"

	"github.com/muesli/cache2go"

	"github.com/byted-apaas/server-common-go/constants"
	"github.com/byted-apaas/server-common-go/structs"
	"github.com/byted-apaas/server-common-go/utils"
)

func GetFunctionMetaConfWithCache(ctx context.Context, apiName string) *structs.FunctionMeta {
	// from local cache
	funcMetaConf := getFunctionMetaConfFromLocalCache(apiName)
	if funcMetaConf != nil {
		return funcMetaConf
	}

	// from remote
	funcMetaConf, err := GetFunctionMetaFromRemote(ctx, apiName)
	if err != nil || funcMetaConf == nil {
		fmt.Printf("GetFunctionMetaConfWithCache failed, apiName: %s, err: %+v", apiName, err)
		return nil
	}

	// save cache
	addFunctionMetaConfToLocalCache(apiName, funcMetaConf)

	return funcMetaConf
}

func getFunctionMetaConfFromLocalCache(funcAPIName string) *structs.FunctionMeta {
	cacheTable := cache2go.Cache(constants.FunctionMetaConfCacheTableKey)
	cacheItem, err := cacheTable.Value(funcAPIName)
	if err != nil {
		return nil
	}
	if value, ok := cacheItem.Data().(*structs.FunctionMeta); ok {
		return value
	}
	return nil
}

func addFunctionMetaConfToLocalCache(funcAPIName string, funcMetaConf *structs.FunctionMeta) {
	cacheTable := cache2go.Cache(constants.FunctionMetaConfCacheTableKey)

	cacheTable.Add(funcAPIName, time.Minute, funcMetaConf)
}

func GetFunctionMetaFromRemote(ctx context.Context, apiName string) (functionMeta *structs.FunctionMeta, err error) {
	ctx = utils.SetApiTimeoutMethodToCtx(ctx, constants.GetFunction)
	functionMeta, err = GetFunctionMetaHttp(ctx, apiName)
	if err != nil {
		return nil, err
	}
	return functionMeta, nil
}
