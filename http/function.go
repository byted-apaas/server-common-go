package http

import (
	"context"
	"fmt"

	"github.com/byted-apaas/server-common-go/constants"
	"github.com/byted-apaas/server-common-go/structs"
	"github.com/byted-apaas/server-common-go/utils"
)

//var (
//	GetFunctionMetaConfFetcher anycache.Fetcher
//)
//
//func init() {
//	GetFunctionMetaConfFetcher = anycache.New(cache.MustNewLocalBytesCacheLRU(5), codec.NewJson(codec.JsonImplIteratorDefault)).
//		WithTTL(time.Minute, time.Minute).
//		WithCacheNil(false).
//		BuildFetcherByLoader(
//			func(ctx context.Context, item interface{}, extraParam interface{}) string {
//				funcAPIName, _ := item.(string)
//				return funcAPIName
//			},
//			func(ctx context.Context, missedItem interface{}, extraParam interface{}) (interface{}, error) {
//				funcAPIName, _ := missedItem.(string)
//				metaConf, err := GetFunctionMetaFromRemote(ctx, funcAPIName)
//				if err != nil {
//					return nil, err
//				}
//				return metaConf, nil
//			},
//		)
//}

func GetFunctionMetaConfWithCache(ctx context.Context, apiName string) *structs.FunctionMeta {
	//funcMetaConf := structs.FunctionMeta{}
	//_, err := GetFunctionMetaConfFetcher.Get(ctx, apiName, &funcMetaConf)
	funcMetaConf, err := GetFunctionMetaFromRemote(ctx, apiName)
	if err != nil {
		fmt.Printf("GetFunctionMetaConfWithCache failed, apiName: %s, err: %+v", apiName, err)
		return nil
	}

	return funcMetaConf
}

func GetFunctionMetaFromRemote(ctx context.Context, apiName string) (functionMeta *structs.FunctionMeta, err error) {
	ctx = utils.SetApiTimeoutMethodToCtx(ctx, constants.GetFunction)
	functionMeta, err = GetFunctionMetaHttp(ctx, apiName)
	if err != nil {
		return nil, err
	}
	return functionMeta, nil
}
