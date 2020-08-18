FROM ubuntu:bionic

EXPOSE 7077
EXPOSE 7078
EXPOSE 8080
EXPOSE 8081

WORKDIR /allspark

ADD init init

RUN apt-get update && \
    apt-get install -y wget && \
    wget https://archive.apache.org/dist/spark/spark-2.4.4/spark-2.4.4-bin-hadoop2.7.tgz && \
    tar -xf spark-2.4.4-bin-hadoop2.7.tgz && \
    rm spark-2.4.4-bin-hadoop2.7.tgz && \
    wget https://cdn.azul.com/zulu/bin/zulu8.40.0.25-ca-jdk8.0.222-linux_x64.tar.gz && \
    tar -xvf zulu8.40.0.25-ca-jdk8.0.222-linux_x64.tar.gz && \
    rm zulu8.40.0.25-ca-jdk8.0.222-linux_x64.tar.gz && \
    chmod +x /allspark/init

ENV SPARK_WORKER_PORT=7078
ENV JAVA_HOME=/allspark/zulu8.40.0.25-ca-jdk8.0.222-linux_x64
ENV SPARK_HOME=/allspark/spark-2.4.4-bin-hadoop2.7

ENTRYPOINT /allspark/init
