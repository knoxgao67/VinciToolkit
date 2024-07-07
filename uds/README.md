### UDS
1. 用go实现 unix domain socket 通信
2. 并给出strace抓包建议

#### USAGE

start server & client
```shell
# server
go run server/main.go -p /tmp/1.sock
# client
go run client/main.go -p /tmp/1.sock
```

start http server & use curl to access
```shell
# server
go run server/main.go -p /tmp/1.sock -http

# curl
curl -v --unix-domain-socket /tmp/1.sock http://localhost/hello -d 'hello from client'

```


use strace to capture the system call
```shell
strace -p ${PID} -t -yy -f -e trace=read,write,close
```