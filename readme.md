
curl test
```
curl 'http://localhost:3666/img/' --data 'url=https%3A%2F%2Fimg.washingtonpost.com%2Frf%2Fimage_1484w%2F2010-2019%2FWashingtonPost%2F2015%2F10%2F19%2FNational-Politics%2FImages%2FCongress_Budget-06d04-3776.jpg&width=300&height=400&quality=10' --compressed
```

run main
```
$ GOPATH=/opt/gocode/ godep go run main.go 
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