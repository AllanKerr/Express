from locust import HttpLocust, TaskSet
import time
import urllib3

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
        print(l.authorization)

def list_searches(l):
    headers = {'Authorization': l.authorization}
    path = get_path("/searches/v1/list")
    l.client.get(path, headers=headers, verify=False)

def list_campgrounds(l):
    path = get_path("/campgrounds/v1/list")
    l.client.get(path, verify=False)

class UserBehavior(TaskSet):
    tasks = {list_searches: 1, list_campgrounds: 2}

    def on_start(self):
        register(self)

class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    min_wait = 5000
    max_wait = 9000

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
