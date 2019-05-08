File Server
===========

[docker image](https://cloud.docker.com/repository/docker/shreve/fileserver)

I created this as a lightweight tool to access and manage files on a remote server. It's intended to be easy to use for both the administrator and end-users as a drop-in replacement for Nginx's autoindex.

All you need to do is expose a port and mount a directory onto /files.

```
docker run -p 80:80 --mount 'type=bind,source=/path/to/data,target=/files' shreve/fileserver
```

This is written with vue.js and go.
