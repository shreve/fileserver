File Server
===========

[docker image](https://cloud.docker.com/repository/docker/shreve/fileserver)

I created this as a lightweight tool to access and manage files on a remote server. It's intended to be easy to use for both the administrator and end-users as a drop-in replacement for Nginx's autoindex.

All you need to do is expose a port and mount a directory onto /files.

```
docker run -p 80:80 --mount 'type=bind,source=/path/to/data,target=/files' shreve/fileserver
```

You can also run with docker-compose with the default config in this repo.

```
docker-compose up
```

This is written with vue.js and go.

## Features

* Looks nice
* Contains entire path from the root
* Compress and download entire directories from the browser
* More coming later I guess

## Development Setup

Change into `client/` then run `npm run serve -- --port=8081`

Change into `server/` then run `env FILES_ROOT=/path go run *.go`

Visit `http://localhost:8081`
