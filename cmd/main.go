package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	target, err := url.Parse("http://localhost:3000")
	if err != nil {
		log.Fatal("Error parsing target URL: ", err)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			fmt.Println(r)
			r.SetURL(target)
			r.Out.Host = r.In.Host
		},
	}

	fmt.Println("Listening on port :8080")
	_ = http.ListenAndServe(":8080", proxy)
}
