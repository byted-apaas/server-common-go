package reference

import "github.com/byted-apaas/server-common-go/structs/common/avatar"

type LookupWithAvatar struct {
	ID       int64          `thrift:"ID,1,required" frugal:"1,required,i64" json:"id" mapstructure:"id"`
	Name     *string        `thrift:"Name,2,optional" frugal:"2,optional,string" json:"name" mapstructure:"name"`
	Avatar   *avatar.Avatar `thrift:"Avatar,3,optional" frugal:"3,optional,avatar.Avatar" json:"avatar" mapstructure:"avatar"`
	TenantID *int64         `thrift:"TenantID,4,optional" frugal:"4,optional,i64" json:"tenant_id" mapstructure:"tenant_id"`
	Email    *string        `thrift:"Email,5,optional" frugal:"5,optional,string" json:"email" mapstructure:"email"`
}
