export GOPATH=`pwd`
go test -v util/netutil
go test -v cloud
go test -v api
go test -v spark_monitor
