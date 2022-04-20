FROM ubuntu:bionic

EXPOSE 7077
EXPOSE 7078
EXPOSE 8080
EXPOSE 8081

WORKDIR /allspark

ADD init init
ADD run run
ADD requirements.txt requirements.txt
ADD run_monitor.py run_monitor.py

RUN apt-get update && \
    apt-get install -y python3.7 \
        libpython3.7-dev \
        python3.7-dev \
        python3-pip \
        jq \
        vim \
        netcat \
        git \
        curl \
        wget \
        multitail \
        unzip && \
    ln -s /usr/bin/python3.7 /usr/bin/python && \
    update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.7 1 && \
    pip3 install --upgrade pip && \
    pip3 install -r requirements.txt && \
    rm requirements.txt && \
    wget https://archive.apache.org/dist/spark/spark-2.4.4/spark-2.4.4-bin-hadoop2.7.tgz && \
    tar -xf spark-2.4.4-bin-hadoop2.7.tgz && \
    rm spark-2.4.4-bin-hadoop2.7.tgz && \
    wget https://cdn.azul.com/zulu/bin/zulu8.40.0.25-ca-jdk8.0.222-linux_x64.tar.gz && \
    tar -xvf zulu8.40.0.25-ca-jdk8.0.222-linux_x64.tar.gz && \
    rm zulu8.40.0.25-ca-jdk8.0.222-linux_x64.tar.gz && \
    chmod +x /allspark/init && \
    chmod +x /allspark/run && \
    chmod +x /allspark/run_monitor.py

ENV SPARK_WORKER_PORT=7078
ENV JAVA_HOME=/allspark/zulu8.40.0.25-ca-jdk8.0.222-linux_x64
ENV SPARK_HOME=/allspark/spark-2.4.4-bin-hadoop2.7

ENTRYPOINT /allspark/init
