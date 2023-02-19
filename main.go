package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	confPath := flag.String("conf", "/config/config.json", "config path")
	flag.Parse()

	file, err := os.ReadFile(*confPath)
	if err != nil {
		panic(err)
	}

	var jwts map[string]struct{}
	if err := json.Unmarshal(file, &jwts); err != nil {
		panic(err)
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		res := json.NewEncoder(os.Stdout)
		_ = res.Encode(req.Header)

		authHeader := req.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer") {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		tokens := strings.SplitN(authHeader, " ", 2)
		if len(tokens) != 2 {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if _, ok := jwts[tokens[1]]; !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	handler.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	const Port = ":8080"
	server := http.Server{Addr: Port, Handler: handler}
	go func() {
		stopSig := make(chan os.Signal, 1)
		signal.Notify(stopSig, syscall.SIGTERM, syscall.SIGINT)
		<-stopSig
		_ = server.Shutdown(context.TODO())
	}()

	fmt.Println("serving at", Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("server start failed", err)
	}
}
