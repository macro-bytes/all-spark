all: setup_docker_network allspark_cli allspark_daemon

allspark_daemon:
	go build -o allspark_daemon ./allspark_orchestrator

allspark_cli:
	go build -o allspark_cli --tags cli ./allspark_orchestrator

setup_docker_network:
	docker network create -d bridge allspark_bridged_newtork || true

run_tests: clean setup_docker_network allspark_cli allspark_daemon
	go test -count=1 -v ./monitor && \
	go test -count=1 -v ./cloud && \
	go test -count=1 -v ./api && \
	go test -count=1 -v ./util/netutil && \
	go test -count=1 -v ./util/serializer && \
	go test -count=1 -v ./datastore

clean:
	docker network rm allspark_bridged_newtork || true && \
	rm -f allspark_cli allspark_daemon
