#!/usr/bin/python

import os
import requests
import time
import json

def run_monitor(cluster_id, callback_url):
    spark_status_url = "http://localhost:8080/json/"
    while True:
        try:
            r = requests.get(url=spark_status_url)
            cluster_status = r.json()
            
            data = {
                "ClusterID": cluster_id,
                "Status": cluster_status
            }

            requests.post(url=callback_url,
                        data=json.dumps(data))
        except:
            ...
        time.sleep(10)

if __name__ == "__main__":
    try:
        run_monitor(os.environ["CLUSTER_ID"],
                    os.environ["ALLSPARK_CALLBACK"])
    except:
        ...
