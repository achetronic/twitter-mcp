package main

import (
	"log"

	"twitter-mcp/internal/globals"
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

	// 2. Create a new MCP server
	mcpServer := server.NewMCPServer(
		appCtx.Config.Server.Name,
		appCtx.Config.Server.Version,
		server.WithToolCapabilities(true),
	)

	// 3. Add Twitter tools to your MCP server
	tm := tools.NewToolsManager(tools.ToolsManagerDependencies{
		AppCtx:        appCtx,
		McpServer:     mcpServer,
		TwitterClient: twitterClient,
	})
	tm.AddTools()

	// 4. Start stdio server
	appCtx.Logger.Info("starting twitter-mcp stdio server")
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatal(err)
	}
}
