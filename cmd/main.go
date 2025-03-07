package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/OsGift/go-api-gateway/internal/middleware"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Route struct {
	Path    string `yaml:"path"`
	Backend string `yaml:"backend"`
}

type Config struct {
	Routes []Route `yaml:"routes"`
}

var routeMap = make(map[string]string)

func loadRoutes() {
	data, err := os.ReadFile("config/routes.yaml")
	if err != nil {
		log.Fatal("Failed to read config file:", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Failed to parse YAML:", err)
	}

	for _, route := range config.Routes {
		routeMap[route.Path] = route.Backend
	}

	fmt.Println("Routes loaded:", routeMap)
}

func reverseProxyHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())

	backend, exists := routeMap[path]
	if !exists {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBody([]byte("Route not found"))
		return
	}

	req := &ctx.Request
	resp := &ctx.Response

	req.SetRequestURI(backend + path) // Forward request

	err := fasthttp.Do(req, resp)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
		ctx.SetBody([]byte("Backend service unavailable"))
	}
}

func main() {
	loadRoutes() // Load routes from YAML

	r := router.New()
	r.GET("/{proxyPath:*}", middleware.RateLimitMiddleware(reverseProxyHandler))

	server := &fasthttp.Server{
		Handler: r.Handler,
		Name:    "GoAPI-Gateway",
	}

	fmt.Println("API Gateway running on :8080")
	log.Fatal(server.ListenAndServe(":8080"))
}
