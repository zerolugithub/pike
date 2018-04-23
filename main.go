package main

import (
	"time"

	"github.com/vicanso/pike/vars"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/vicanso/pike/cache"
	"github.com/vicanso/pike/middleware"
	"github.com/vicanso/pike/proxy"
)

func main() {
	// Echo instance
	e := echo.New()
	client := &cache.Client{
		Path: "/tmp/test.cache",
	}

	err := client.Init()
	if err != nil {
		panic(err)
	}
	defer client.Close()
	directors := make(proxy.Directors, 0)
	d := &proxy.Director{
		Name: "aslant",
		Ping: "/ping",
		Backends: []string{
			"http://127.0.0.1:5018",
		},
	}
	go d.StartHealthCheck(5 * time.Second)
	directors = append(directors, d)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if c.IsWebSocket() {
				return vars.ErrNotSupportWebSocket
			}
			return next(c)
		}
	})

	e.Use(customMiddleware.Identifier(client))

	e.Use(customMiddleware.DirectorPicker(directors))

	// e.Use(middleware.Gzip())

	e.Use(customMiddleware.ProxyWithConfig(customMiddleware.ProxyConfig{}))

	e.Use(customMiddleware.Dispatcher(client))

	// e.Use(func(c echo.Context) (err error) {
	// 	body := c.Get(vars.Body)
	// 	if body == nil {
	// 		// TODO ERROR
	// 		return errors.New("Get the response fail")
	// 	}

	// })

	// Routes

	// Start server
	e.Logger.Fatal(e.Start(":3015"))

}
