export GOPATH=$(PWD)

all: allspark_cli allspark_daemon

allspark_daemon:
	go build -o allspark_daemon allspark_orchestrator

allspark_cli:
	go build -o allspark_cli --tags cli allspark_orchestrator

install_prerequisites:
	go get github.com/docker/docker/client && \
	go get github.com/aws/aws-sdk-go/... && \
	go get -u github.com/Azure/azure-sdk-for-go/... && \
	go get -u github.com/go-redis/redis

run_tests: allspark_cli allspark_daemon
	go test -count=1 -v util/netutil && \
	go test -count=1 -v cloud && \
	go test -count=1 -v api && \
	go test -count=1 -v monitor && \
	go test -count=1 -v util/serializer && \
	go test -count=1 -v datastore

clean:
	rm -f allspark_cli allspark_daemon
