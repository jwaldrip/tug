# tug

Use Docker for development

## Prerequisites

* [Docker][docker] >= 1.3
* [Golang][golang] >= 1.3
* Working Docker host (non-Linux users see [boot2docker](http://boot2docker.io/))

## Installation

    $ go get github.com/nitrous-io/tug

## Set up your application
    
### Create a Tugfile

    web:      bin/web -p $PORT
    postgres: docker/postgres:9.3.5
    redis:    docker/redis:2.8.9

> Any command that starts with `docker/` will be interpreted as a docker image tag.

### Create `bin/bootstrap`

If your app needs to do any setup before it starts, put it in a `bin/bootstrap` file:

    $ cat bin/bootstrap
    #!/bin/sh
    bundle exec rake db:migrate

> Make sure your `bin/bootstrap` has a +x bit set.

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

##### Example Dockerfile

<pre>
FROM ruby:2.1.2

ENV PORT 3000
<b>EXPOSE 3000</b>

WORKDIR /app
<b>ADD . /app</b>
</pre>

[docker]: https://www.docker.com/whatisdocker/
[golang]: http://golang.org/

## Contributors

Tug is sponsored by [Nitrous.IO](https://www.nitrous.io/) and built by [these contributors](https://github.com/nitrous-io/tug/graphs/contributors).
