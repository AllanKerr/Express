from locust import HttpLocust, TaskSet
import time

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
        l.authorization = "Authorization: bearer " + r.json()["access_token"]
        print(l.authorization)

def test(l):
    print("empty")

class UserBehavior(TaskSet):
    tasks = {test: 1} #{profile: 1}

    def on_start(self):
        register(self)


class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    min_wait = 5000
    max_wait = 9000
