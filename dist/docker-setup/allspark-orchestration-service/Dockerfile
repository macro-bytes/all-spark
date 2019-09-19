FROM clearlinux:latest

EXPOSE 32418

WORKDIR /allspark

ADD allspark_daemon allspark_daemon
ADD config.json config.json

ENTRYPOINT /allspark/allspark_daemon /allspark/config.json
