package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"httpserver/metrics"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	flag.Set("v", "4")
	flag.Parse()
	defer glog.Flush()
	glog.V(2).Info("Starting http server...")
	metrics.Register()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//log.Println("Starting http server...")
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/header", header)
	mux.HandleFunc("/logging", logging)
	mux.HandleFunc("/hello", hello)
	mux.Handle("/metrics", promhttp.Handler())

	go func() {
		err := http.ListenAndServe(":80", mux)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for sig := range c {
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			glog.Info("Waiting 5 seconds for grace shutdown...")
			time.Sleep(time.Second * 5)
			ShutDown()
		}
	}
	glog.Info("Server Exited Properly")
}

func ShutDown() {
	glog.Info("Starting quit...")
	os.Exit(0)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	log.Print("Entering healthz handler...")
	io.WriteString(w, "ok\n")
}

func hello(w http.ResponseWriter, r *http.Request) {
	if r.URL.RequestURI() == "/favicon.ico" {
		return
	}
	glog.Info("Entering root handler...")
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	delay := randInt(2000, 5000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("Hello, [%s]\n", user))
	} else {
		io.WriteString(w, fmt.Sprintf("Hello, [stranger]\n"))
	}
	io.WriteString(w, "Reading request headers to response...\n")
	header := r.Header
	header.Set("Delay", fmt.Sprintf("%d ms", delay))
	for k, v := range header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}

}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func header(w http.ResponseWriter, r *http.Request) {
	glog.Info("Entering header handler...")
	version := os.Getenv("VERSION")
	w.Header().Set("VERSION", version)
	glog.Infof("VERSION is %s", version)
	for k, v := range r.Header {
		for _, vv := range v {
			w.Header().Set(k, vv)
		}
	}
}

func logging(w http.ResponseWriter, r *http.Request) {
	glog.Info("Entering logging handler...")
	ip := ClientIP(r)
	io.WriteString(w, fmt.Sprintf("IP is: %s", ip))
	glog.Infof("IP is: %s", ip)

}

func getCurrentIP(r *http.Request) string {
	// 这里也可以通过X-Forwarded-For请求头的第一个值作为用户的ip
	// 但是要注意的是这两个请求头代表的ip都有可能是伪造的
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		// 当请求头不存在即不存在代理时直接获取ip
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	return ip
}

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
