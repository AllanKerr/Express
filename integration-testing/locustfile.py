from locust import HttpLocust, TaskSet

def profile(l):
    l.client.get("/test/", verify=False)

class UserBehavior(TaskSet):
    tasks = {profile: 1}

class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    min_wait = 5000
    max_wait = 9000
