package manage

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	mcp "trpc.group/trpc-go/trpc-mcp-go"
)

func TestRegisterMCPServerRoutes(t *testing.T) {
	ctx := context.Background()

	mApi := func(engRouter *gin.RouterGroup) {}
	pApi := func(engRouter *gin.RouterGroup) {}
	feApi := func(engRouter *gin.RouterGroup) {}
	initCB := func(am *AppManage) error { return nil }
	doCB := func(code string, am *AppManage) {}
	destroyCB := func(am *AppManage) error { return nil }

	// AppManage internal port 8081, MCP ports 13001, 13002
	am, err := NewAppManage(ctx, "test_app_mcp", 8081, mApi, pApi, feApi, initCB, doCB, destroyCB, "/tmp/test_sys_dir_mcp", true, nil)
	assert.NoError(t, err)

	_, err = am.RegisterMCPServer("s1", "1.0.0", 13001)
	assert.NoError(t, err)
	_, err = am.RegisterMCPServer("s2", "1.0.0", 13002)
	assert.NoError(t, err)

	time.Sleep(500 * time.Millisecond) // Wait for servers to start

	// Check s1
	resp, err := http.Get("http://localhost:13001/s1/sse")
	if err != nil {
		t.Logf("Failed to connect to s1: %v", err)
	} else {
		assert.NotEqual(t, http.StatusNotFound, resp.StatusCode)
		resp.Body.Close()
	}

	// Check s2
	resp2, err := http.Get("http://localhost:13002/s2/sse")
	if err != nil {
		t.Logf("Failed to connect to s2: %v", err)
	} else {
		assert.NotEqual(t, http.StatusNotFound, resp2.StatusCode)
		resp2.Body.Close()
	}

	srv, ok := am.GetMCPServer("s1")
	assert.True(t, ok)
	assert.NotNil(t, srv)
}

func TestMCPServerToolCall(t *testing.T) {
	ctx := context.Background()

	// 1. Setup AppManage
	mApi := func(engRouter *gin.RouterGroup) {}
	pApi := func(engRouter *gin.RouterGroup) {}
	feApi := func(engRouter *gin.RouterGroup) {}
	initCB := func(am *AppManage) error { return nil }
	doCB := func(code string, am *AppManage) {}
	destroyCB := func(am *AppManage) error { return nil }

	// AppManage port 8082, MCP port 13003
	am, err := NewAppManage(ctx, "test_tool_app", 8082, mApi, pApi, feApi, initCB, doCB, destroyCB, "/tmp/test_tool_sys", true, nil)
	assert.NoError(t, err)

	// 2. Register MCP Server
	port := 13003
	s, err := am.RegisterMCPServer("tool-server", "1.0.0", port)
	assert.NoError(t, err)

	// 3. Register Tool
	tool := mcp.NewTool("hello",
		mcp.WithDescription("Say hello"),
	)
	s.RegisterTool(tool, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewTextResult("Hello, world!"), nil
	})

	time.Sleep(500 * time.Millisecond) // Wait for server to start

	// 4. Create Client
	// Use the SSE endpoint: http://localhost:13003/tool-server/sse
	endpoint := fmt.Sprintf("http://localhost:%d/tool-server/sse", port)
	t.Logf("Connecting to SSE endpoint: %s", endpoint)

	client, err := mcp.NewSSEClient(endpoint, mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	})
	assert.NoError(t, err)
	defer client.Close()

	// 5. Initialize Client
	initReq := &mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: "2024-11-05",
			ClientInfo: mcp.Implementation{
				Name:    "test-client",
				Version: "1.0.0",
			},
			Capabilities: mcp.ClientCapabilities{},
		},
	}
	_, err = client.Initialize(ctx, initReq)
	assert.NoError(t, err)

	err = client.SendInitialized(ctx)
	assert.NoError(t, err)

	// 6. Call Tool
	result, err := client.CallTool(ctx, &mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "hello",
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Hello, world!", result.Content[0].(mcp.TextContent).Text)
}
