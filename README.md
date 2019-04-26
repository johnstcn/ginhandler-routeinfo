# ginhandler-routeinfo

[![GoDoc](https://godoc.org/github.com/johnstcn/ginhandler-routeinfo?status.svg)](https://godoc.org/github.com/johnstcn/ginhandler-routeinfo)

A gin middleware to set the route path in the request context.

This middleware keeps an in-memory map of the routes handled by the Gin engine, and attempts to reverse-lookup the routing path based on the handler name.

For example, take the following code using a fictituous library `infoCat` to log request information:

```go
func handleHello(g *gin.Context) {
	name := g.Param("name")
	g.Status(http.StatusOK)
	g.Writer.WriteString("Hello " + name)
}

func logRequest() gin.HandlerFunc {
	return func(g *gin.Context) {
		metrics := map[string]string{
			"method": g.Request.Method,
			"endpoint": g.Request.URL.String(),
		}
		infoCat.Log(metrics)
	}
}

func main() {
	g := gin.New()
	g.Use(logRequest())
	g.Handle(http.MethodGet, "/hello/:name", handleHello)

	log.Fatal(http.ListenAndServe(":8000", g))
}
```

Running this and hitting `http://localhost:8000/hello/jimmy` will cause `infoCat` to log a potentially high-cardinality metric `endpoint:/hello/$USER`.

This middleware allows you to instead do the following:

```go
func logRequest() gin.HandlerFunc {
	return func(g *gin.Context) {
		endpoint := g.GetString(ginrouteinfo.RoutePathKey)
		metrics := map[string]string{
			"method": g.Request.Method,
			"endpoint": endpoint,
		}
		infoCat.Log(metrics)
	}
}

func main() {
	g := gin.New()
	g.Use(ginrouteinfo.WithRoutePath(g.Routes))
	g.Handle(http.MethodGet, "/hello/:name", handleHello)

	log.Fatal(http.ListenAndServe(":8000", g))
}
```

Now, `infoCat` will instead log `endpoint` as the route path `/hello/:name`.
