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
and contain it in a container.

There is also a `docker-compose.yml` file. This file will be responsible for creating our `main.go` container and MongoDB container. 
Once they are created, those two will be contained in a container cluster as if they are in the same LAN so they can communicate with each other.

To execute the `docker-compose.yml` file, run the following command

```bash 
docker compose run

docker compose stop # to pause

docker compose down # to remove

docker rmi dbmq_app # to remove the image for dbmq
```

If you are not sure what to delete or it says `image not found`, run `docker images -a` to list all the existing images. 

You should run `docker rmi dbmq_app` if you make some changes in `main.go` or in any of the files in `handlers`. If you don't remove, the new changes won't be included in the image. You can imagine a docker image is like an ISO to install a certain Linux distribution like Ubuntu.

Once those two containers are up and running, open MongoDB Compass. If you want to connect to your MongoDB database, you need to type this URL in the URL field and press connect.

```bash
mongodb://localhost:27011
```

Why is the port 27011? Because in the docker compose, the external port is set to 27011 and the internal port to 27017 in the `docker-compose.yml` file. Since our machine is out of the docker cluster, we can only connect to the database from 27011.

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

