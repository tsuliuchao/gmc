package main

import (
	"net/http"
	"time"

	"github.com/snail007/gmc/util/timeutil"

	"github.com/snail007/gmc"
)

func main() {
	api := gmc.NewAPIServer(":8030")
	// add a middleware typed 1 to filter all request registered in router,
	// exclude 404 requests.
	api.AddMiddleware1(func(c gmc.C, s *gmc.APIServer) (isStop bool) {
		s.Logger().Printf("before request %s", c.Request.RequestURI)
		return false
	})
	// add a middleware typed 2 to logging every request registered in router,
	// exclude 404 requests.
	api.AddMiddleware2(func(c gmc.C, s *gmc.APIServer) (isStop bool) {
		s.Logger().Printf("after request %s %d %d %s", c.Request.Method, c.StatusCode(), c.WriteCount(), c.Request.RequestURI)
		return false
	})
	// sets a function to handle 404 requests.
	api.Handle404(func(c gmc.C) {
		c.Write("404")
	})
	// sets a function to handle panic error.
	api.Handle500(func(c gmc.C, err interface{}) {
		c.WriteHeader(http.StatusInternalServerError)
		c.Write("panic error : ", err)
	})
	// sets an api in url path: /hello
	// add more api , just call api.API() repeatedly
	api.API("/hello", func(c gmc.C) {
		api.Logger().Printf("request %s", c.Request.RequestURI)
		c.Write("hello world!")
	}).API("/hi", func(c gmc.C) {
		c.Write("hi!")
	})
	// trigger a panic error
	api.API("/error", func(c gmc.C) {
		a := 0
		a /= a
	})
	api.API("/version", func(c gmc.C) {
		c.Write(1.1)
	})
	// routing by group is supported
	group1 := api.Group("/v1")
	group1.API("/time", func(c gmc.C) {
		c.Write(timeutil.TimeToStr(time.Now()))
	})

	app := gmc.NewAPP().SetNoneMainConfigFile(true)

	err := app.ParseConfig()
	if err != nil {
		panic(err)
	}

	app.AddService(gmc.ServiceItem{
		Service: api,
		BeforeInit: func(cfg *gmc.Config) (err error) {
			api.PrintRouteTable(nil)
			return
		},
	})

	app.Logger().Panic(app.Run())
}