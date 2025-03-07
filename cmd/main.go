package main

import (
	"fmt"
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func reverseProxyHandler(ctx *fasthttp.RequestCtx) {
	// Forward request to backend service (hardcoded for now)
	// backendURL := "http://localhost:8081"
	req := &ctx.Request
	resp := &ctx.Response

	// Proxy request
	err := fasthttp.Do(req, resp)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
		ctx.SetBody([]byte("Backend service unavailable"))
	}
}

func main() {
	r := router.New()
	r.GET("/{proxyPath:*}", reverseProxyHandler)

	server := &fasthttp.Server{
		Handler: r.Handler,
		Name:    "GoAPI-Gateway",
	}

	fmt.Println("API Gateway running on :8080")
	log.Fatal(server.ListenAndServe(":8080"))
}
