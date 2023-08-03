package main

import (
	"github.com/Chista-Framework/Chista/middlewares"
	"github.com/Chista-Framework/Chista/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	//config.Connect() -> Connect to db, implement when it'll be necessary
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(middlewares.LoggingMiddleware())
	router.Use(middlewares.CORSMiddleware())

	// Modules: ioc, phishing, leak, threat_monitor, blacklist, source, c2
	// Set the routes for each module
	routes.IocRoute(router)
	routes.PhishingRoute(router)
	routes.LeakRoute(router)
	routes.ThreatMonitorRoute(router)
	routes.BlacklistRoute(router)
	routes.SourceRoute(router)
	routes.C2Route(router)

	// Serve the API
	router.Run("0.0.0.0:7777")
}
