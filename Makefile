all: setup_docker_network allspark_cli allspark_daemon

allspark_daemon:
	go build -o allspark_daemon ./allspark_orchestrator

allspark_cli:
	go build -o allspark_cli --tags cli ./allspark_orchestrator

setup_docker_network:
	docker network create -d bridge allspark_bridged_newtork || true

run_tests: clean pull_required_images setup_docker_network allspark_cli allspark_daemon
	export APP_EXIT_STATUS_PATH="/tmp/allspark_exit_status" && \
	docker run --name dev-spark-cluster -d -p 8080:8080 macrobytes/dev-spark-cluster:latest && \
	go test -timeout 2400s -count=1 -v ./monitor && \
	go test -timeout 2400s -count=1 -v ./cloud && \
	go test -timeout 2400s -count=1 -v ./api && \
	go test -timeout 2400s -count=1 -v ./util/netutil && \
	go test -timeout 2400s -count=1 -v ./util/serializer && \
	go test -timeout 2400s -count=1 -v ./datastore && \
    python3 dist/docker-setup/allspark-compute/test_run_monitor.py && \
	docker rm -f /dev-spark-cluster || true

pull_required_images:
	docker pull macrobytes/dev-spark-cluster:latest && \
	docker pull macrobytes/allspark-compute:latest

clean:
	rm /tmp/allspark_exit_status; \
	redis-cli flushall; \
	docker network rm allspark_bridged_newtork 2>/dev/null; \
	rm -f allspark_cli allspark_daemon /allspark/exit_status; \
	docker rm -f /dev-spark-cluster 2>/dev/null || true

build_orchestrator_image:
	docker build -t macrobytes/allspark-orchestration-service -f dist/docker-setup/allspark-orchestration-service/Dockerfile .