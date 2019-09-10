#!/usr/bin/python

import os
import requests
import time

def run_monitor(cluster_id, callback_url):
    spark_status_url = "http://localhost:8080/json/"
    while True:
        r = requests.get(url=spark_status_url)
        cluster_status = r.json()
        
        data = {
            "ClusterID": cluster_id,
            "Status": cluster_status
        }

        requests.post(url=callback_url,
                      data=data)
        
        time.sleep(10)

if __name__ == "__main__":
    run_monitor(os.environ["CLUSTER_ID"],
                os.environ["CALLBACK_URL"])