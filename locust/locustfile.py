from locust import HttpUser, task, between
from random import randint, random, choice
from dataclasses import dataclass
import json
import sqlite3


@dataclass
class Account:
    def __init__(self, id) -> None:
        self.id = id


db = sqlite3.connect("./data.db")


class BankUser(HttpUser):
    wait_time = between(10, 15)
    password = ""
    username = ""
    full_name = ""
    access_token = ""
    refresh_token = ""
    user_id = ""
    services = []

    @task
    def run_transaction(self):
        source = choice(self.services).id

        destination = db.execute(
            "select id from services order by RANDOM() limit 1"
        ).fetchone()[0]

        if source == destination:
            return

        account_response = self.client.get(
            f"/services/{source}",
            headers={"Authorization": f"Bearer {self.access_token}"},
        )
        if account_response.status_code != 200:
            return

        quantity = (random() ** 2) * (
            float(account_response.json()["balance"])
            + float(account_response.json()["init_balance"])
        )
        transaction_response = self.client.post(
            "/transactions",
            json={
                "currency": "USD",
                "amount": f"{quantity:.2f}",
                "source": source,
                "destination": destination,
            },
            headers={"Authorization": f"Bearer {self.access_token}"},
        )
        if transaction_response.status_code != 200:
            return

    def register(self):
        registration_response = self.client.post(
            "/users",
            json={
                "username": self.username,
                "password": self.password,
                "fullname": self.full_name,
            },
        )
        if registration_response.status_code != 200:
            return False

        return True

    def authenticate(self):
        auth_response = self.client.post(
            "/auth",
            json={
                "username": self.username,
                "password": self.password,
            },
        )

        if auth_response.status_code != 200:
            return False

        self.access_token = auth_response.json()["access_token"]
        self.refresh_token = auth_response.json()["refresh_token"]
        self.user_id = auth_response.json()["id"]

        return True

    def create_service(self, funds):
        service_response = self.client.post(
            f"/users/{self.user_id}/services",
            json={
                "type": "CHQ",
                "currency": "USD",
                "init_balance": f"{funds:.2f}",
            },
            headers={"Authorization": f"Bearer {self.access_token}"},
        )

        if service_response.status_code != 200:
            self.stop()
            raise Exception("could not create service")

        id = service_response.json()["id"]

        db.execute("insert into services values (?)", (id,))
        db.commit()

        self.services.append(Account(id))

    def get_services(self):
        services_response = self.client.get(
            f"/users/{self.user_id}/services",
            headers={"Authorization": f"Bearer {self.access_token}"},
        )

        data = [
            Account(json.loads(line)["id"])
            for line in services_response.iter_lines(decode_unicode=True)
            if not line
        ]

        for elem in data:
            db.execute("insert into services values (?)", (elem,))

        return data

    def on_start(self) -> None:
        result = db.execute(
            "select fullname, username, password from users order by RANDOM() limit 1"
        )
        self.full_name, self.username, self.password = result.fetchone()

        if not self.authenticate() and not self.register():
            self.stop()
            raise Exception("could not register")

        if not self.authenticate():
            self.stop()
            raise Exception("could not authenticate")

        self.services = self.get_services()
        if len(self.services) == 0:
            self.create_service(randint(100, 1_000_000))
