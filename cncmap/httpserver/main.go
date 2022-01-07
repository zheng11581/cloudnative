package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	//flag.Set("v", "4")
	flag.Parse()
	glog.V(2).Info("Starting http server...")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthy", healthyHandler)
	http.HandleFunc("/header", headerHandler)
	http.HandleFunc("/logging", loggingHandler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func healthyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entering healthy handler...")
	io.WriteString(w, "ok\n")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entering root handler...")
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("Hello, [%s]\n", user))
	} else {
		io.WriteString(w, fmt.Sprintf("Hello, [stranger]\n"))
	}
	io.WriteString(w, "Reading request headers to response...\n")
	header := r.Header
	for k, v := range header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}

}

func headerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entering header handler...")
	io.WriteString(w, "Reading request headers to response...\n")
	header := r.Header
	resHeader := w.Header()
	for k, v := range header {
		resHeader[k] = v
	}

	version := os.Getenv("VERSION")
	resHeader["Version"] = []string{version}

	for k, v := range resHeader {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
}

func loggingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entering logging handler...")
	fmt.Println("未完成")
	io.WriteString(w, "ok")

}
