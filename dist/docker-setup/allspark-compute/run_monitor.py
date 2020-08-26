#!/usr/bin/python3

import os
import requests
import time
import json
from typing import Dict, Any

APP_EXIT_STATUS_PATH = "/allspark/exit_status"

def get_app_exit_status() -> str:
    """
    Returns the exit status of the Spark run script (either "ERROR" or "") 
    """
    if os.path.isfile(APP_EXIT_STATUS_PATH):
        with open(APP_EXIT_STATUS_PATH, "r") as fh:
            status = fh.read()
    else:
        status = ""
    return status

def get_cluster_status() -> Dict[str, Any]:
    """
    Returns the Spark cluster status
    :return: Dict[str, Any]
    """
    spark_status_url = "http://localhost:8080/json/"
    r = requests.get(url=spark_status_url)
    return r.json()

def run_monitor(cluster_id: str, callback_url: str):
    while True:
        try:
            data = {
                "ClusterID": cluster_id,
                "Status": get_cluster_status(),
                "AppExitStatus": get_app_exit_status(),
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
