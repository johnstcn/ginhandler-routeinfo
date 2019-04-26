package ginrouteinfo

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const testHandlerPath = "/test/:param"

func handleGetTest(g *gin.Context) {
	routePath := g.GetString(RoutePathKey)
	if routePath == "" {
		g.Status(http.StatusNotFound)
		return
	}
	if routePath != testHandlerPath {
		g.Status(http.StatusInternalServerError)
		return
	}
	g.Status(http.StatusOK)
}

func TestWithRoutePath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	g := gin.New()
	g.Use(WithRoutePath(g.Routes))
	g.Handle(http.MethodGet, testHandlerPath, handleGetTest)

	req := httptest.NewRequest(http.MethodGet, "/test/something", nil)
	rr := httptest.NewRecorder()
	g.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Result().StatusCode)
}

func BenchmarkWithRoutePath(b *testing.B) {
	// hit with a random string each time
	randStr := func(len int) string {
		buf := make([]byte, len)
		_, _ = rand.Read(buf)
		str := base64.StdEncoding.EncodeToString(buf)
		return str[:len]
	}
	gin.SetMode(gin.TestMode)
	g := gin.New()
	g.Use(WithRoutePath(g.Routes))
	g.Handle(http.MethodGet, testHandlerPath, handleGetTest)

	rr := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test/" + randStr(12), nil)
		g.ServeHTTP(rr, req)
	}
}