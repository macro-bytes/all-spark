export GOPATH=`pwd`
go test -count=1 -v util/netutil
go test -count=1 -v cloud
go test -count=1 -v api
go test -count=1 -v spark_monitor
go test -count=1 -v util/serializer
go test -count=1 -v datastore
