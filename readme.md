Comentarismo IMG Resizer Project - Golang <3

## About

This project is used for few tasks:

* Resize images and cache it with Redis with a TTL
* Create memes and cache it with Redis with a TTL
* Create gifs and cache it with Redis with a TTL


## Installing

* Requirements: Golang, Godep, Redis

The project should have all dependencies it needs inside the vendor folder. If you add new dependency, make sure to run godep save

## Running
* `make start`

## Stop
* `make stop`

## Logs
* `make log`


## Run with defaults for g7-box
```
$ GOPATH=/opt/gocode/ REDIS_HOST=g7-box REDIS_PORT=6379 godep go run main.go
```


## References:
```
* https://github.com/llgcode/draw2d
* https://github.com/fogleman/gg
```

ENV ARGS:
```
REDIS_HOST
REDIS_PASSWORD
REDIS_PORT
```

## Manual Test PROD
```
curl 'http://img.comentarismo.com/r/img/' --data 'url=https://i.ytimg.com/vi/YdOQGkQ1KFs/hqdefault.jpg&width=388&height=395&quality=50'
```

## Manual Test local
```
url 'http://localhost:3666/r/img/' --data 'url=https://i.ytimg.com/vi/YdOQGkQ1KFs/hqdefault.jpg&width=380&height=395&quality=50'
```
