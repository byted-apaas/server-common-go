package http

import (
	"context"
	"fmt"

	"github.com/byted-apaas/server-common-go/structs"
	"github.com/byted-apaas/server-common-go/utils"
)

// AppendParamsUnauthFields 处理参数中的权限
// 提取权限信息
func AppendParamsUnauthFields(ctx context.Context, funcAPIName string, inputOrOutput string, params interface{}, unauthFieldsMap map[string]interface{}) (context.Context, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	param := map[string]interface{}{}
	err := utils.Decode(params, &param)
	if err != nil {
		fmt.Printf("[AppendParamsUnauthFields] Failed, err: %+v\n", err)
		return ctx, nil
	}

	if len(param) == 0 {
		return ctx, nil
	}

	metaConf := utils.GetFunctionMetaConfFromCtx(ctx, funcAPIName)
	if metaConf == nil {
		metaConf = GetFunctionMetaConfWithCache(ctx, funcAPIName)
	}
	if metaConf == nil {
		return ctx, nil
	}

	var ioParams []*structs.IOParamItem
	if inputOrOutput == "input" {
		ioParams = metaConf.IOParam.Input
	} else if inputOrOutput == "output" {
		ioParams = metaConf.IOParam.Output
	}

	for _, p := range ioParams {
		if p.Type == "Record" {
			AppendUnauthFieldRecord(ctx, p.ObjectAPIName, param[p.Key], utils.ParseStrList(unauthFieldsMap[p.Key]))
		} else if p.Type == "RecordList" {
			AppendUnauthFieldRecordList(ctx, p.ObjectAPIName, param[p.Key], utils.ParseStrsList(unauthFieldsMap[p.Key]))
		}
	}
	return ctx, nil
}

// CalcParamsNeedPermission 计算参数需要的权限，返回权限信息
func CalcParamsNeedPermission(ctx context.Context, funcAPIName string, inputOrOutput string, params interface{}) *structs.Permission {
	if params == nil {
		return nil
	}

	data := map[string]interface{}{}
	err := utils.Decode(params, &data)
	if err != nil {
		fmt.Printf("CalcParamsNeedPermission decode params failed, funcAPIName: %s, err: %+v\n", funcAPIName, err)
		return nil
	}

	metaConf := utils.GetFunctionMetaConfFromCtx(ctx, funcAPIName)
	if metaConf == nil {
		metaConf = GetFunctionMetaConfWithCache(ctx, funcAPIName)
	}
	if metaConf == nil {
		return nil
	}

	var ioParams []*structs.IOParamItem
	if inputOrOutput == "input" {
		ioParams = metaConf.IOParam.Input
	} else if inputOrOutput == "output" {
		ioParams = metaConf.IOParam.Output
	}

	perm := structs.Permission{UnauthFields: map[string]interface{}{}}
	for _, p := range ioParams {
		if p.Type == "Record" {
			id := utils.GetRecordID(data[p.Key])
			perm.UnauthFields[p.Key] = utils.GetRecordUnauthFieldByObjectAndRecordID(ctx, p.ObjectAPIName, id)
		} else if p.Type == "RecordList" {
			var newRecords []*structs.RecordOnlyID
			err := utils.Decode(data[p.Key], &newRecords)
			if err != nil {
				fmt.Printf("CalcParamsNeedPermission decode RecordList failed, err: %+v\n", err)
				return nil
			}

			var unauthFieldsList [][]string
			hasUnauthFields := false
			for _, record := range newRecords {
				unauthFields := utils.GetRecordUnauthFieldByObjectAndRecordID(ctx, p.ObjectAPIName, record.GetID())
				if len(unauthFields) > 0 {
					hasUnauthFields = true
				}
				unauthFieldsList = append(unauthFieldsList, unauthFields)
			}

			if hasUnauthFields {
				perm.UnauthFields[p.Key] = unauthFieldsList
			}
		}
	}

	if len(perm.UnauthFields) > 0 {
		return &perm
	}

	return nil
}

func AppendUnauthFieldRecord(ctx context.Context, objectAPIName string, recordOrID interface{}, unauthFields []string) {
	if objectAPIName == "" || recordOrID == nil || len(unauthFields) == 0 {
		return
	}

	id, ok := recordOrID.(int64)
	if !ok {
		id = utils.GetRecordID(recordOrID)
	}
	if id <= 0 {
		return
	}

	unauthFieldMap := utils.GetRecordUnauthField(ctx)
	if unauthFieldMap == nil {
		return
	}

	if v, ok := unauthFieldMap[objectAPIName]; !ok || v == nil {
		unauthFieldMap[objectAPIName] = map[int64][]string{}
	}

	unauthFieldMap[objectAPIName][id] = unauthFields
}

func AppendUnauthFieldRecordList(ctx context.Context, objectAPIName string, records interface{}, unauthFieldsList [][]string) {
	if objectAPIName == "" || records == nil || len(unauthFieldsList) == 0 {
		return
	}

	var newRecords []*structs.RecordOnlyID
	err := utils.Decode(records, &newRecords)
	if err != nil {
		fmt.Printf("AppendUnauthFieldRecordList failed, err: %+v\n", err)
		return
	}

	if len(newRecords) != len(unauthFieldsList) {
		fmt.Printf("len(record)(%d) != len(unauthFieldsList)(%d)\n", len(newRecords), len(unauthFieldsList))
		return
	}

	for i := 0; i < len(newRecords); i++ {
		record := newRecords[i]
		unauthFields := unauthFieldsList[i]
		AppendUnauthFieldRecord(ctx, objectAPIName, record.GetID(), unauthFields)
	}
}
