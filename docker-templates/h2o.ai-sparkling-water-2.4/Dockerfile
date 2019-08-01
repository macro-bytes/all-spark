FROM macrobytes/spark-2.4-standalone:latest

WORKDIR /spark

ENV SPARK_HOME=/spark/spark-2.4.3-bin-hadoop2.7
ENV MASTER="spark://localhost:7077"

EXPOSE 54321
EXPOSE 54322

ADD run /spark/run
ADD init /spark/init

RUN wget -q https://s3.amazonaws.com/h2o-release/sparkling-water/rel-2.4/13/sparkling-water-2.4.13.zip && \
    unzip sparkling-water-2.4.13.zip && \
    rm sparkling-water-2.4.13.zip && \
    chmod +x /spark/run && \
    chmod +x /spark/init
