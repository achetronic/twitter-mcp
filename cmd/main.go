package main

import (
	"log"
	"net/http"
	"time"

	"twitter-mcp/internal/globals"
	"twitter-mcp/internal/handlers"
	"twitter-mcp/internal/middlewares"
	"twitter-mcp/internal/tools"
	"twitter-mcp/internal/twitter"

	"github.com/mark3labs/mcp-go/server"
)

func main() {

	// 0. Process the configuration
	appCtx, err := globals.NewApplicationContext()
	if err != nil {
		log.Fatalf("failed creating application context: %v", err.Error())
	}

	// 1. Initialize Twitter client
	twitterClient := twitter.NewClient(
		appCtx.Config.Twitter.APIKey,
		appCtx.Config.Twitter.APIKeySecret,
		appCtx.Config.Twitter.AccessToken,
		appCtx.Config.Twitter.AccessTokenSecret,
		appCtx.Config.Twitter.BearerToken,
	)

	// 2. Initialize middlewares that need it
	accessLogsMw := middlewares.NewAccessLogsMiddleware(middlewares.AccessLogsMiddlewareDependencies{
		AppCtx: appCtx,
	})

	jwtValidationMw, err := middlewares.NewJWTValidationMiddleware(middlewares.JWTValidationMiddlewareDependencies{
		AppCtx: appCtx,
	})
	if err != nil {
		appCtx.Logger.Info("failed starting JWT validation middleware", "error", err.Error())
	}

	// 3. Create a new MCP server
	mcpServer := server.NewMCPServer(
		appCtx.Config.Server.Name,
		appCtx.Config.Server.Version,
		server.WithToolCapabilities(true),
	)

	// 4. Initialize handlers for later usage
	hm := handlers.NewHandlersManager(handlers.HandlersManagerDependencies{
		AppCtx: appCtx,
	})

	// 5. Add Twitter tools to your MCP server
	tm := tools.NewToolsManager(tools.ToolsManagerDependencies{
		AppCtx:        appCtx,
		McpServer:     mcpServer,
		Middlewares:   []middlewares.ToolMiddleware{},
		TwitterClient: twitterClient,
	})
	tm.AddTools()

	// 6. Wrap MCP server in a transport (stdio, HTTP, SSE)
	switch appCtx.Config.Server.Transport.Type {
	case "http":
		httpServer := server.NewStreamableHTTPServer(mcpServer,
			server.WithHeartbeatInterval(30*time.Second),
			server.WithStateLess(false))

		// Register it under a path, then add custom endpoints.
		// Custom endpoints are needed as the library is not feature-complete according to MCP spec requirements
		// Ref: https://modelcontextprotocol.io/specification/2025-06-18/basic/authorization#overview
		mux := http.NewServeMux()
		mux.Handle("/mcp", accessLogsMw.Middleware(jwtValidationMw.Middleware(httpServer)))

		if appCtx.Config.OAuthAuthorizationServer.Enabled {
			mux.Handle("/.well-known/oauth-authorization-server"+appCtx.Config.OAuthAuthorizationServer.UrlSuffix,
				accessLogsMw.Middleware(http.HandlerFunc(hm.HandleOauthAuthorizationServer)))
		}

		if appCtx.Config.OAuthProtectedResource.Enabled {
			mux.Handle("/.well-known/oauth-protected-resource"+appCtx.Config.OAuthProtectedResource.UrlSuffix,
				accessLogsMw.Middleware(http.HandlerFunc(hm.HandleOauthProtectedResources)))
		}

		// Start StreamableHTTP server with proper timeouts for long-lived connections
		httpSrv := &http.Server{
			Addr:              appCtx.Config.Server.Transport.HTTP.Host,
			Handler:           mux,
			ReadHeaderTimeout: 10 * time.Second,
			IdleTimeout:       0, // Disable idle timeout for SSE/streaming connections
		}

		appCtx.Logger.Info("starting StreamableHTTP server", "host", appCtx.Config.Server.Transport.HTTP.Host)
		err := httpSrv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}

	default:
		// Start stdio server
		appCtx.Logger.Info("starting stdio server")
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatal(err)
		}
	}
}
