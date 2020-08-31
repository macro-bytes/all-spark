all: setup_docker_network allspark_cli allspark_daemon

allspark_daemon:
	go build -o allspark_daemon ./allspark_orchestrator

allspark_cli:
	go build -o allspark_cli --tags cli ./allspark_orchestrator

setup_docker_network:
	docker network create -d bridge allspark_bridged_newtork || true

run_tests: clean setup_docker_network allspark_cli allspark_daemon
	docker run --name dev-spark-cluster -d -p 8080:8080 macrobytes/dev-spark-cluster:latest && \
	go test -timeout 1200s -count=1 -v ./monitor && \
	go test -timeout 1200s -count=1 -v ./cloud && \
	go test -timeout 1200s -count=1 -v ./api && \
	go test -timeout 1200s -count=1 -v ./util/netutil && \
	go test -timeout 1200s -count=1 -v ./util/serializer && \
	go test -timeout 1200s -count=1 -v ./datastore && \
    python3 dist/docker-setup/allspark-compute/test_run_monitor.py && \
    docker rm -f dev-spark-cluster

clean:
	docker network rm allspark_bridged_newtork || true && \
	rm -f allspark_cli allspark_daemon && /allspark/app_exit_status && \
	docker rm -f dev-spark-cluster || true
