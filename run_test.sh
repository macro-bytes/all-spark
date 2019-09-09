export GOPATH=`pwd`

go build -tags cli allspark_orchestrator # cli
go build allspark_orchestrator           # daemon 

go test -count=1 -v util/netutil && \
go test -count=1 -v cloud && \
go test -count=1 -v api && \
go test -count=1 -v monitor && \
go test -count=1 -v util/serializer && \
go test -count=1 -v datastore
