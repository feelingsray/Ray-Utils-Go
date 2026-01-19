package manage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	mcp "trpc.group/trpc-go/trpc-mcp-go"
)

func TestRawMCPServer(t *testing.T) {
	s := mcp.NewServer("test", "1.0.0")
	handler := s.HTTPHandler()
	ts := httptest.NewServer(handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/sse")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestRegisterMCPServerRoutes(t *testing.T) {
	ctx := context.Background()

	mApi := func(engRouter *gin.RouterGroup) {}
	pApi := func(engRouter *gin.RouterGroup) {}
	feApi := func(engRouter *gin.RouterGroup) {}
	initCB := func(am *AppManage) error { return nil }
	doCB := func(code string, am *AppManage) {}
	destroyCB := func(am *AppManage) error { return nil }

	am, err := NewAppManage(ctx, "test_app_mcp", 8081, mApi, pApi, feApi, initCB, doCB, destroyCB, "/tmp/test_sys_dir_mcp", true, nil)
	assert.NoError(t, err)

	_, err = am.RegisterMCPServer("s1", "1.0.0")
	assert.NoError(t, err)
	_, err = am.RegisterMCPServer("s2", "1.0.0")
	assert.NoError(t, err)

	ts := httptest.NewServer(am.engRouter)
	defer ts.Close()

	for _, name := range []string{"s1", "s2"} {
		resp, err := http.Get(ts.URL + "/mcp/" + name + "/sse")
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusNotFound, resp.StatusCode)
		resp.Body.Close()
	}

	srv, ok := am.GetMCPServer("s1")
	assert.True(t, ok)
	assert.NotNil(t, srv)
}
