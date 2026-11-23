package main

import (
	"github.com/zaheershaikh936/CommitLaunchBackend/cmd/internal/config"
	"github.com/zaheershaikh936/CommitLaunchBackend/cmd/internal/routes"
)

func main() {
	route := config.NewRouter()
	routes.ApplicationRoutes(route)
	config.ServerConfig(route)
}
