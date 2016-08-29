

run main g7-box
```
$ GOPATH=/opt/gocode/ REDIS_HOST=g7-box REDIS_PORT=6379 godep go run main.go 
```

ENV ARGS:
```
REDIS_HOST
REDIS_PASSWORD
REDIS_PORT
```

build main
```
$ GOPATH=/opt/gocode/ godep go build main.go 
```


run executable (OSX)
`./main`
 

```
$ nohup ./imgresizer & .
```

Find PID
```
netstat -tulpn | grep 3666

```

Kill with PID

```
kill -9 PID
```



```
flynn create imgresizer 
Created imgresizer
```

```
flynn log -f
```

Test PROD
```
curl 'http://img.comentarismo.com/r/img/' --data 'url=https://i.ytimg.com/vi/YdOQGkQ1KFs/hqdefault.jpg&width=388&height=395&quality=50'
```

Test local
```
url 'http://localhost:3666/r/img/' --data 'url=https://i.ytimg.com/vi/YdOQGkQ1KFs/hqdefault.jpg&width=380&height=395&quality=50'
```
