import random
import threading

import time
from locust import HttpLocust, TaskSet, task
from requests import get, post
from hosts import host_resolver

import sys

from lwwes.lwwes import Lwwes

passed_host = host_resolver.getHost(sys.argv)
post("{}/status/reset".format(passed_host))

def generate_request():
    operation = "ADD" if random.randint(0, 1) == 1 else "REMOVE"
    data = {"Value": {"Value": random.randint(0, 100)}, "Operation": operation}
    return data

class CrdtTasks(TaskSet):
    requests_count = 0
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
    time_to_wait = 2
    print("Waiting for {} seconds for results".format(time_to_wait))
    time.sleep(time_to_wait)
    with CrdtTasks.lock:
        requests_count = CrdtTasks.requests_count

    print(requests_count, WebsiteUser.host)
    current_value = get("{}/status/readable".format(WebsiteUser.host))
    # current_value = int(.content) // TODO check state
    # if current_value != requests_count:
    #     raise Exception("Inconsistent system! Expected {} got {}".format(requests_count, current_value))
    # print(current_value)
    print("Finish")


from locust.events import quitting

quitting_event = quitting
quitting_event += compare_results
