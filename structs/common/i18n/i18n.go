package i18n

type I18ns = []*I18n

type I18n struct {
	LanguageCode int64  `thrift:"LanguageCode,1,required" frugal:"1,required,i64" json:"language_code" mapstructure:"language_code"`
	Text         string `thrift:"Text,2,required" frugal:"2,required,string" json:"text" mapstructure:"text"`
}

type I18nCnUs struct {
	ZhCn string `thrift:"ZhCn,1,required" frugal:"1,required,string" json:"zh_CN"`
	EnUs string `thrift:"EnUs,2,required" frugal:"2,required,string" json:"en_US"`
}
