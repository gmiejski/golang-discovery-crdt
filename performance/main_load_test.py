import random
import threading

import time
from locust import HttpLocust, TaskSet, task
from requests import get, post

import crdt_app
from crdt_app.data_consistency_test import test_data_consistency
from hosts import locust_host_resolver

import sys

from lwwes.lwwes import Lwwes
from crdt_app.host_resolver import get_all_hosts

passed_host = locust_host_resolver.getLocustHost(sys.argv)
for host in get_all_hosts(passed_host):
    print("Resetting host {}".format(host))
    post("{}/status/reset".format(host))

def generate_request():
    operation = "ADD" if random.randint(0, 1) == 1 else "REMOVE"
    data = {"Value": {"Value": random.randint(0, 100)}, "Operation": operation}
    return data


class CrdtTasks(TaskSet):
    lock = threading.Lock()
    lwwes = Lwwes()

    def on_start(self):
        time.sleep(1)

    @task
    def index(self):
        generated_data = generate_request()
        with CrdtTasks.lock:
            CrdtTasks.lwwes.add(generated_data["Value"]["Value"], generated_data["Operation"])
        self.client.post("/status/update", json=generated_data)


class WebsiteUser(HttpLocust):
    task_set = CrdtTasks
    min_wait = 100
    max_wait = 1000
    host = passed_host


def compare_results():
    print("Performance test ended")
    time_to_wait = 3
    print("Waiting for {} seconds for results".format(time_to_wait))
    time.sleep(time_to_wait)
    test_data_consistency(WebsiteUser.host, sorted(CrdtTasks.lwwes.get()))
    print("Finish")


def get_server_value(host):
    current_value = get("{}/status/readable".format(host))
    json_rs = current_value.json()
    return sorted(list(map(lambda x: int(x), json_rs["Values"])))


from locust.events import quitting

quitting_event = quitting
quitting_event += compare_results
