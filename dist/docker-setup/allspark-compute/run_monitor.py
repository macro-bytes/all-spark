#!/usr/bin/python3

import os
import requests
import time
import json
from typing import Dict, Any

APP_EXIT_STATUS_PATH = os.environ.get("APP_EXIT_STATUS_PATH", "/allspark/exit_status")

def get_app_exit_status() -> str:
    """
    Returns the exit status of the Spark run script (either "ERROR" or "") 
    """
    if os.path.isfile(APP_EXIT_STATUS_PATH):
        with open(APP_EXIT_STATUS_PATH, "r") as fh:
            status = fh.read().strip()
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

def get_local_status() -> Dict[str, Any]:
    """
    Returns a running status if localhost:4040/api/v1/applications is reachable; othe
    :return: Dict[str, Any]
    """
    idle_cluster_state = {
        "url" : "spark://simulated-local-mode-cluster:7077",
        "workers" : [ ],
        "aliveworkers" : 0,
        "cores" : 0,
        "coresused" : 0,
        "memory" : 0,
        "memoryused" : 0,
        "activeapps" : [ ],
        "completedapps" : [ ],
        "activedrivers" : [ ],
        "completeddrivers" : [ ],
        "status" : "ALIVE"
    }

    running_cluster_state = {
        "url" : "spark://simulated-local-mode-cluster::7077",
        "workers" : [ ],
        "aliveworkers" : 0,
        "cores" : 0,
        "coresused" : 0,
        "memory" : 0,
        "memoryused" : 0,
        "activeapps" : [ {
            "id" : "app-0-0000",
            "starttime" : 0,
            "name" : "pyspark-shell",
            "cores" : 0,
            "user" : "spark-user",
            "memoryperslave" : 0,
            "submitdate" : "Thu Jan 13 17:47:31 GMT 2022",
            "state" : "WAITING",
            "duration" : 0
        } ],
        "completedapps" : [ ],
        "activedrivers" : [ ],
        "completeddrivers" : [ ],
        "status" : "ALIVE"
    }

    spark_status_url = "http://localhost:4040/api/v1/applications"
    try:
        r = requests.get(url=spark_status_url)
        if r.status_code != 200:
            return idle_cluster_state
        return running_cluster_state
    except:
        return idle_cluster_state


def run_monitor(cluster_id: str, callback_url: str, cluster_mode):
    while True:
        status = get_cluster_status() if cluster_mode > 0 else get_local_status()
        try:
            data = {
                "ClusterID": cluster_id,
                "Status": status,
                "AppExitStatus": get_app_exit_status(),
            }

            requests.post(url=callback_url,
                          data=json.dumps(data))
        except:
            ...
        time.sleep(10)

if __name__ == "__main__":
    try:
        cluster_mode = True if int(os.environ["EXPECTED_WORKERS"]) > 0 else False
        run_monitor(os.environ["CLUSTER_ID"],
                    os.environ["ALLSPARK_CALLBACK"],
                    cluster_mode)
    except:
        ...
