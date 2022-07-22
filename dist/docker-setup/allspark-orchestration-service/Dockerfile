FROM clearlinux/golang:1.18 AS builder

WORKDIR /allspark-orchestrator
ADD go.mod .
ADD go.sum .
COPY allspark_orchestrator allspark_orchestrator
COPY api api
COPY cloud cloud
COPY daemon daemon
COPY datastore datastore
COPY logger logger
COPY monitor monitor
COPY util util

RUN go build -o allspark_daemon ./allspark_orchestrator

FROM clearlinux:latest

EXPOSE 32418

WORKDIR /allspark

COPY --from=builder /allspark-orchestrator/allspark_daemon allspark_daemon
ADD dist/docker-setup/allspark-orchestration-service/config.local.json config.json

ENTRYPOINT ["/allspark/allspark_daemon"] 
CMD ["/allspark/config.json"]
