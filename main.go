package main

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	defer println("exiting main")

	exitchan := make(chan os.Signal, 1)
	signal.Notify(exitchan, syscall.SIGINT, syscall.SIGTERM)

	println("running ...")

	close := initPprof()

	<-exitchan

	close()
}

func initPprof() func() {
	startchan := make(chan os.Signal, 1)
	signal.Notify(startchan, syscall.SIGUSR1)

	stopchan := make(chan os.Signal, 1)
	signal.Notify(stopchan, syscall.SIGUSR2)

	exitchan := make(chan bool)

	var server *http.Server

	shutdown := func() {
		if server != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			println("shutting down pprof server")
			server.Shutdown(ctx)
		}
	}

	close := func() {
		shutdown()
		exitchan <- true
	}

	router := http.NewServeMux()
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	go func() {
		defer println("pprof launcher exited")
		for {
			select {
			case <-exitchan:
				return
			case <-stopchan:
				shutdown()
			case <-startchan:
				server = &http.Server{
					Addr:    "localhost:6060",
					Handler: router,
				}
				go func() {
					println("starting http pprof server")
					log.Println(server.ListenAndServe())
					println("http server pprof shutted down")
					server = nil
				}()
			}
		}
	}()

	return close
}
