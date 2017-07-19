import threading

import time
from locust import HttpLocust, TaskSet, task
from requests import get


class CrdtTasks(TaskSet):
    variable = 0
    lock = threading.Lock()

    def on_start(self):
        self.client.post("/status/reset")
        time.sleep(1)

    @task
    def index(self):
        self.client.post("/status/update")
        with CrdtTasks.lock: CrdtTasks.variable += 1


class WebsiteUser(HttpLocust):
    task_set = CrdtTasks
    min_wait = 1000
    max_wait = 1000


def compare_results():
    print("Performance test ended")
    with CrdtTasks.lock:
        variable = CrdtTasks.variable
    print(variable, WebsiteUser.host)
    value = int(get("{}/status".format(WebsiteUser.host)).content)
    print(value)
    print("Quitting 2")



from locust.events import quitting

quitting_event = quitting
quitting_event += compare_results
