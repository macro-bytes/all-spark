FROM clearlinux:latest

EXPOSE 32418

ENV GOPATH=/allspark/allspark-orchestrator

WORKDIR /allspark

ADD allspark_orchestrator /allspark/allspark_orchestrator
ADD config.json /allspark/config.json

ENTRYPOINT /allspark/allspark_orchestrator /allspark/config.json
