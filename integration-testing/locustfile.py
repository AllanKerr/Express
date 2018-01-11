from locust import HttpLocust, TaskSet
import time
import urllib3
import campgrounds
from random import *
import json
import string
from locust.exception import StopLocust

def get_path(path):
    return "/testapi" + path

def random_username():
    if 'num' not in random_username.__dict__:
        random_username.num = int(time.time())
    random_username.num += 1
    return "express-" + str(random_username.num)


 #***************************************************************************************
 #    Python Snippets: How to Generate Random String
 #    Author: Jackson Cooper
 #    Date: Jan. 10, 2018
 #    Code version: 1.0
 #    Availability: https://www.pythoncentral.io/python-snippets-how-to-generate-random-string/
 #
 #***************************************************************************************/
def random_password():
    min_char = 8
    max_char = 12
    allchar = string.ascii_letters + string.punctuation + string.digits
    password = "".join(choice(allchar) for x in range(randint(min_char, max_char)))
    return password

def register(l):
    username = random_username()
    password = random_password()

    r = l.client.post("/oauth2/register", {
        "username": username,
        "password": password,
        "confirm-password": password
    }, verify=False)

    if r.status_code == 200:
        l.username = username
        l.password = password
        l.authorization = "bearer " + r.json()["access_token"]
        l.refresh = r.json()["refresh_token"]

def login(l):
    r = l.client.post("/oauth2/login", {
        "username": l.username,
        "password": l.password,
    }, verify=False)
    if r.status_code == 200:
        l.authorization = "bearer " + r.json()["access_token"]
        l.refresh = r.json()["refresh_token"]

def refresh(l):
    r = l.client.post("/oauth2/token", {
        "grant_type": "refresh_token",
        "refresh_token": l.refresh,
    }, auth=('admin', 'demo-password'), verify=False)
    if r.status_code == 200:
        l.authorization = "bearer " + r.json()["access_token"]
        l.refresh = r.json()["refresh_token"]

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
    tasks = {list_searches: 30, list_campgrounds: 30, add_search: 15, login: 2, refresh: 2, register: 1}

    def on_start(self):
        self.count = 0
        register(self)

class WebsiteUser(HttpLocust):
    weight = 30
    task_set = UserBehavior
    min_wait = 5000
    max_wait = 9000

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
