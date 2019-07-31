FROM alpine:latest

ADD init /spark/init
ADD run /spark/run

ENV SPARK_WORKER_PORT=7078

EXPOSE 7077
EXPOSE 7078
EXPOSE 8080
EXPOSE 8081

RUN  apk update && \
    apk add wget && \
    apk add bash && \
    apk add procps && \
    apk add coreutils && \
    apk add openjdk8 && \
    apk add python && \
    apk add python3 && \
    apk add curl && \
    apk add html2text && \
    cd /spark && \
    wget https://www-us.apache.org/dist/spark/spark-2.4.3/spark-2.4.3-bin-hadoop2.7.tgz && \
    tar -xf spark-2.4.3-bin-hadoop2.7.tgz && \
    rm spark-2.4.3-bin-hadoop2.7.tgz && \
    chmod +x /spark/init && \
    chmod +x /spark/run

ENTRYPOINT /spark/init
