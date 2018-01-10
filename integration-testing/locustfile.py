from locust import HttpLocust, TaskSet
import time
import urllib3
import campgrounds
from random import randrange
import json

def get_path(path):
    return "/testapi" + path

def username():
    if 'num' not in username.__dict__:
        username.num = int(time.time())
    username.num += 1
    return "express-" + str(username.num)

def register(l):
    r = l.client.post("/oauth2/register", {
        "username": username(),
        "password":"password",
        "confirm-password":"password"
    }, verify=False)
    if r.status_code == 200:
        l.authorization = "bearer " + r.json()["access_token"]

def list_searches(l):
    headers = {'Authorization': l.authorization}
    path = get_path("/searches/v1/list")
    r = l.client.get(path, headers=headers, verify=False)

def list_campgrounds(l):
    path = get_path("/campgrounds/v1/list")
    l.client.get(path, verify=False)

def add_search(l):
    headers = {
        'Authorization': l.authorization,
        'content-type': 'application/json'
    }
    items = []
    for i in range (0, randrange(1,5)):
        items.append(campgrounds.random_campground())

    path = get_path("/searches/v1/add")
    l.client.post(path,  data=json.dumps({
        "campgrounds": items,
        "rangeStart":"2018-01-01",
        "rangeEnd":"2018-12-01",
        "nights":randrange(1,5)
    }), headers=headers, verify=False)

class UserBehavior(TaskSet):
    tasks = {list_searches: 1, list_campgrounds: 1, add_search: 1}

    def on_start(self):
        register(self)

class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    min_wait = 5000
    max_wait = 9000

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
