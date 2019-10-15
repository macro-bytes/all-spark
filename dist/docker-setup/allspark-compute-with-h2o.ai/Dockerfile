FROM macrobytes/allspark-compute:latest

ADD init init

ADD h2odriver-sw3.26.6-2-extended.jar.* h2o/
RUN cat h2o/h2odriver-sw3.26.6-2-extended.jar.* >>h2odriver-sw3.26.6-2-extended.jar && \
    rm -Rf h2o/

