# Golang + MongoDB + RabbitMQ

## Tools
1. Docker Desktop \
   [Download Docker Desktop here](https://docs.docker.com/desktop/mac/install/).
   If you are using M1 chip, click the one with `Mac with Apple chip`
2. Docker Compose \
   Open your Docker Desktop and click `Setting`. In the `Setting`, check `Use Docker Compose V2`.
3. MongoDB + MongoDB Compass


## Setup

You will notice there a file called `Dockerfile`. This file will be responsible for building our `main.go` file 
and contain in a container.

There is also a `docker-compose.yml` file. This file will be responsible for creating our `main.go` container and MongoDB container.

To execute the `docker-compose.yml` file, run the following command

```bash 
docker compose run
```

The command above will ramp up 2 containers. The first container will be the `main.go` container and the second container will be the MongoDB container.

Once those two containers are up and running, open MongoDB Compass. If you want to connect to your MongoDB database, you need to type this URL in the URL field and press connect.

```bash
mongodb://localhost:27011
```

Why is the port 27011? Because in the docker compose, the external port is set to 27011 and the internal port to 27017 in the `docker-compose.yml` file.

```yaml
...
  mongo:
    image: mongo:latest
    container_name: mongo
    volumes:
      - ./mongo-volume:/data/db
    ports:
      - 27011:27017 # here
...
```

In your MongoDB Compass, create a new database and name it the way you wish. In your new database, create a new collection and name it the way you wish. Populate the collection with some data.

To make sure that our API is working, open a terminal and run the following command

```bash
curl -X GET http://localhost:9090/GET/
```

If you want to create a new user, run the following command

```bash
curl -X POST http://localhost:9090/POST/ -d '{"name":"John", "age":30}'

curl -X POST http://localhost:9090/POST/ -d "@data.json"
```

Have fun :)

