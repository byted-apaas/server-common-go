package avatar

import "github.com/byted-apaas/server-common-go/structs/common/i18n"

type Avatar struct {
	Source  string            `thrift:"Source,1,required" frugal:"1,required,string" json:"source" mapstructure:"source"`
	Image   map[string]string `thrift:"Image,2,optional" frugal:"2,optional,map<string:string>" json:"image" mapstructure:"image"`
	Color   *string           `thrift:"Color,3,optional" frugal:"3,optional,string" json:"color" mapstructure:"color"`
	Content i18n.I18ns        `thrift:"Content,4,optional" frugal:"4,optional,list<i18n.I18n>" json:"content" mapstructure:"content"`
	ColorID *string           `thrift:"ColorID,5,optional" frugal:"5,optional,string" json:"color_id" mapstructure:"color_id"`
}
