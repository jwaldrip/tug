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

If any command starts with `docker/` the rest will be interpreted as a docker image tag.

## Start the app

    $ tug start
    postgres | fixing permissions on existing onesory /var/lib/postgresql/data ... ok
    postgres | creating subdirectories ... ok
    postgres | selecting default max_connections ... 100
    postgres | selecting default shared_buffers ... 128MB
    web      | listening on 0.0.0.0:5000

## Container linking

Tug will set environment variables in the Docker [container linking format](https://docs.docker.com/userguide/dockerlinks/#environment-variables), like this:

    POSTGRES_PORT_5432_TCP=tcp://127.0.0.1:5000
    POSTGRES_PORT_5432_TCP_PROTO=tcp
    POSTGRES_PORT_5432_TCP_ADDR=127.0.0.1
    POSTGRES_PORT_5432_TCP_PORT=5000

##### Aliasing ENV vars

If your application expects env vars to be named differently, alias them in your `Tugfile`:

    web: env DATABASE_HOST=$POSTGRES_PORT_5432_ADDR bin/web

## Dockerfile

If your app has a `Dockerfile`, tug will use it to build and run your app in [Docker][docker] while setting up appropriate port forwarding and file synchronization.

For Tug to work most effectively your `Dockerfile` should include the following:

* The listening web port should be specified with an `EXPOSE` statement
* The app's code should be included from the local directory using an `ADD` statement
* The app's startup command should be defined using `CMD`
* The app's command in `Tugfile` should be empty.

##### Example Dockerfile

<pre>
FROM ruby:2.1.2

ENV PORT 3000
<b>EXPOSE 3000</b>

WORKDIR /app
<b>ADD . /app</b>

<b>CMD ["bundle", "exec", "unicorn", "-p", "$PORT"]</b>
</pre>

[docker]: https://www.docker.com/whatisdocker/
[golang]: http://golang.org/

##### Updated Tugfile

<pre>
web:      # from dockerfile
postgres: docker/postgres:9.3.5
redis:    docker/redis:2.8.9
</pre>
