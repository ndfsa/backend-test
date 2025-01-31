import sqlite3
from faker import Faker

fake = Faker()
db = sqlite3.connect("data.db")
db.execute("drop table if exists users")
db.execute("drop table if exists services")
db.execute("create table users(fullname varchar, username varchar, password varchar)")
db.execute("create table services(id varchar)")
db.commit()

data = []
for i in range(1, 1_000):
    data.append((fake.name(), fake.user_name(), fake.password()))

db.executemany("insert into users values(?, ?, ?)", data)
db.commit()
