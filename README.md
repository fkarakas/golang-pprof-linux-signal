# Goal
Activate pprof endpoint with the linux signal USR1 (shutdown with USR2)

# Why ?
Because you need to investigate an issue on a running GO program on production.

Of course, you can always activate pprof endpoint but this approch has 2 issues : 
- security: it is not recommended to expose an endpoint if you don't really need, most of security breach are comming from web server running in background and exposing internal information about a program.
- performance: even if the pprof sampling will not make a big impact on the running program, it has an impact on the performance.

# Linux signal

To interact with the go process we can use linux signal USR1 and USR2 which can be used as custom signal

So to activate the pprof:
```
$ kill -USR1 $PID
$ curl http://localhost:6060/debug/pprof/goroutine\?debug\=2
goroutine 14 [running]:
runtime/pprof.writeGoroutineStacks(0x1496a00, 0xc0001921c0, 0xc000078840, 0x0)

```

And to deactivate the pprof:
```
$ kill -USR2 $PID
$ curl http://localhost:6060/debug/pprof/goroutine\?debug\=2
curl: (7) Failed to connect to localhost port 6060: Connection refused
```
