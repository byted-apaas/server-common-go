package structs

type FunctionMeta struct {
	ApiName string          `json:"apiName"`
	IOParam FunctionIOParam `json:"io_param"`
}

type FunctionIOParam struct {
	Input  []*IOParamItem `json:"input"`
	Output []*IOParamItem `json:"output"`
}

type IOParamItem struct {
	Key           string `json:"key" `
	Type          string `json:"type"`
	ObjectAPIName string `json:"objectApiName"`
}
