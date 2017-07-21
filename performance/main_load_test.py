import threading

import time
from locust import HttpLocust, TaskSet, task
from requests import get, post
from hosts import host_resolver

import sys
passed_host = host_resolver.getHost(sys.argv)
post("{}/status/reset".format(passed_host))


class CrdtTasks(TaskSet):
    requests_count = 0
    lock = threading.Lock()

    def on_start(self):
        time.sleep(1)

    @task
    def index(self):
        with CrdtTasks.lock:
            CrdtTasks.requests_count += 1
        self.client.post("/status/update")


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
    current_value = int(get("{}/status".format(WebsiteUser.host)).content)
    if current_value != requests_count:
        raise Exception("Inconsistent system! Expected {} got {}".format(requests_count, current_value))
    print(current_value)
    print("Finish")


from locust.events import quitting

quitting_event = quitting
quitting_event += compare_results
