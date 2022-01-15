package main

import (
	"fmt"
	"github.com/golang/glog"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	//flag.Set("v", "4")
	////flag.Parse()
	//glog.V(2).Info("Starting http server...")
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/header", header)
	mux.HandleFunc("/logging", logging)
	mux.HandleFunc("/", index)
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entering healthz handler...")
	//glog.V(2).Infof("The /healthz return: ")
	io.WriteString(w, "ok\n")
}

func index(w http.ResponseWriter, r *http.Request) {
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

func header(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entering header handler...")
	os.Setenv("VERSION", "1.0.0")
	version := os.Getenv("VERSION")
	w.Header().Set("VERSION", version)
	glog.V(2).Infof("VERSION is %s", version)
	for k, v := range r.Header {
		for _, vv := range v {
			w.Header().Set(k, vv)
		}
	}
}

func logging(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entering logging handler...")
	ip := ClientIP(r)
	io.WriteString(w, fmt.Sprintf("IP is: %s", ip))
	glog.V(2).Infof("IP is: %s", ip)

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
