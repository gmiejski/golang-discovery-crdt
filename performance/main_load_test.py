import threading

import time
from locust import HttpLocust, TaskSet, task
from requests import get


class WebsiteTasks(TaskSet):
    variable = 0
    lock = threading.Lock()

    def on_start(self):
        self.client.post("/status/reset")
        time.sleep(1)

    @task
    def index(self):
        self.client.get("/status/update")
        with WebsiteTasks.lock: WebsiteTasks.variable += 1


class WebsiteUser(HttpLocust):
    task_set = WebsiteTasks
    min_wait = 1000
    max_wait = 1000


def on_my_event():
    print("Quitting 1")
    print(WebsiteTasks.variable, WebsiteUser.host)
    value = int(get("{}/status".format(WebsiteUser.host)).content)
    print(value)
    print("Quitting 2")


from locust.events import quitting

my_event = quitting
my_event += on_my_event
