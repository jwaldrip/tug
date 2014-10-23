# tug

Use Docker for development

## Prerequisites

* [Docker][docker] >= 1.3
* [Golang][golang] >= 1.3
* Working `DOCKER_HOST` (non-Linux users see [boot2docker](http://boot2docker.io/))

## Installation

    $ go get github.com/nitrous-io/tug
    
## Create a Tugfile

    web:      bin/web -p $PORT
    postgres: docker/postgres:9.3.5
    redis:    docker/redis:2.8.9

## Start the app

    $ tug start
    postgres | fixing permissions on existing directory /var/lib/postgresql/data ... ok
    postgres | creating subdirectories ... ok
    postgres | selecting default max_connections ... 100
    postgres | selecting default shared_buffers ... 128MB
    web      | listening on 0.0.0.0:5000

## Container linking

Tug will set environment variables in the Docker [container linking](https://docs.docker.com/userguide/dockerlinks/#environment-variables) format, like this:

    POSTGRES_PORT_5432_TCP=tcp://127.0.0.1:5000
    POSTGRES_PORT_5432_TCP_PROTO=tcp
    POSTGRES_PORT_5432_TCP_ADDR=127.0.0.1
    POSTGRES_PORT_5432_TCP_PORT=5000

##### Aliasing ENV vars

If your application expects env vars to be named differently, alias them in your Tugfile command:

    web: env DATABASE_HOST=$POSTGRES_PORT_5432_ADDR bin/web

## Dockerfile

If your repo has a `Dockerfile`, tug will use it and run your app in [Docker][docker] setting up the proper port forwards and directory synchronization.

Your `Dockerfile` should do the following:

* Inject the code into the container with an `ADD` statement
* Expose the listening web port using an `EXPOSE` statement
* Start the app using a `CMD` statement

##### Example Dockerfile

    FROM ruby:2.1.2
    ENV PORT 3000
    EXPOSE 3000
    ADD . /app
    CMD bin/start-my-web-app -p 3000

[docker]: https://www.docker.com/whatisdocker/
[golang]: http://golang.org/
