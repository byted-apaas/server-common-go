module github.com/byted-apaas/server-common-go

go 1.16

require (
	github.com/json-iterator/go v1.1.12
	github.com/muesli/cache2go v0.0.0-20221011235721-518229cd8021
	github.com/sirupsen/logrus v1.9.0
	github.com/tidwall/gjson v1.9.3
	go.mongodb.org/mongo-driver v1.8.3
)

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0
