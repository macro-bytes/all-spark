FROM alpine:latest

ADD run /spark/run

RUN  apk update && \
    apk add wget && \
    apk add bash && \
    apk add procps && \
    apk add coreutils && \
    apk add openjdk8 && \
    apk add python && \
    apk add python3 && \
    cd /spark && \
    wget https://www-us.apache.org/dist/spark/spark-2.4.3/spark-2.4.3-bin-hadoop2.7.tgz && \
    tar -xf spark-2.4.3-bin-hadoop2.7.tgz && \
    rm spark-2.4.3-bin-hadoop2.7.tgz && \
    chmod +x /spark/run

ENTRYPOINT /spark/run
