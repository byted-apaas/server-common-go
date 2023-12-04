module github.com/byted-apaas/server-common-go

go 1.16

require (
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/json-iterator/go v1.1.12
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/muesli/cache2go v0.0.0-20221011235721-518229cd8021
	github.com/sirupsen/logrus v1.9.0
	github.com/tidwall/gjson v1.9.3
	go.mongodb.org/mongo-driver v1.8.3
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0

// Deprecated
retract v0.0.23
// Deprecated
retract v0.0.24