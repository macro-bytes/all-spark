FROM clearlinux:latest

EXPOSE 7077
EXPOSE 7078
EXPOSE 8080
EXPOSE 8081

WORKDIR /allspark
ADD init init
ADD run run
ADD requirements.txt requirements.txt
ADD run_monitor.py run_monitor.py

RUN swupd bundle-add curl \
        git \
        jq \
        python3-basic \
        wget \
        sysadmin-basic \
        unzip && \
    pip install -r requirements.txt && \
    rm requirements.txt && \
    wget https://www.apache.org/dist/spark/spark-2.4.4/spark-2.4.4-bin-hadoop2.7.tgz && \
    tar -xf spark-2.4.4-bin-hadoop2.7.tgz && \
    rm spark-2.4.4-bin-hadoop2.7.tgz && \
    wget https://s3.amazonaws.com/h2o-release/sparkling-water/rel-2.4/13/sparkling-water-2.4.13.zip && \
    unzip sparkling-water-2.4.13.zip && \
    rm sparkling-water-2.4.13.zip && \
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
