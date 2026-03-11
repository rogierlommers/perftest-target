import os
import random
import time
from locust import HttpUser, between, events, task

STRESS_START_AFTER_SECONDS = int(os.getenv("STRESS_START_AFTER_SECONDS", "120"))
TEST_START_MONO = None


@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    global TEST_START_MONO
    TEST_START_MONO = time.monotonic()


@events.test_stop.add_listener
def on_test_stop(environment, **kwargs):
    global TEST_START_MONO
    TEST_START_MONO = None


class PerfTestUser(HttpUser):
    # Allows overriding via env var, while still supporting --host in locust CLI.
    host = os.getenv("LOCUST_HOST", "http://localhost:3000")
    wait_time = between(0.1, 1.0)

    def on_start(self):
        # Fail fast if service is unavailable when a user starts.
        with self.client.get("/health", name="GET /health", catch_response=True) as response:
            if response.status_code != 200:
                response.failure(f"Expected 200, got {response.status_code}")

    @task(25)
    def get_users(self):
        with self.client.get("/users", name="GET /users", catch_response=True) as response:
            if response.status_code != 200:
                response.failure(f"Expected 200, got {response.status_code}")

    @task(10)
    def post_users(self):
        payload = {
            "username": f"user-{random.randint(1, 1_000_000)}"
        }
        with self.client.post("/users", json=payload, name="POST /users", catch_response=True) as response:
            if response.status_code != 201:
                response.failure(f"Expected 201, got {response.status_code}")

    @task(20)
    def get_tasks(self):
        with self.client.get("/tasks", name="GET /tasks", catch_response=True) as response:
            if response.status_code != 200:
                response.failure(f"Expected 200, got {response.status_code}")

    @task(25)
    def get_documents(self):
        with self.client.get("/documents", name="GET /documents", catch_response=True) as response:
            if response.status_code != 200:
                response.failure(f"Expected 200, got {response.status_code}")

    @task(10)
    def post_documents(self):
        payload = {
            "title": f"Document {random.randint(1, 1_000_000)}",
            "content": "Synthetic load-test content",
        }
        with self.client.post("/documents", json=payload, name="POST /documents", catch_response=True) as response:
            if response.status_code != 201:
                response.failure(f"Expected 201, got {response.status_code}")

    @task(12)
    def get_stress_cpu(self):
        if TEST_START_MONO is None:
            return

        if time.monotonic() - TEST_START_MONO < STRESS_START_AFTER_SECONDS:
            return

        # Aggressive stress-demo mode: hit CPU burn endpoint frequently and harder.
        path = "/stress/cpu?duration_ms=2000&workers=8"
        with self.client.get(path, name="GET /stress/cpu", catch_response=True) as response:
            if response.status_code != 200:
                response.failure(f"Expected 200, got {response.status_code}")
