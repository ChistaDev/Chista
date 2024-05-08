package main

import (
	"time"

	"fmt"

	"github.com/Chista-Framework/Chista/controller"
	"github.com/Chista-Framework/Chista/helpers"
	"github.com/Chista-Framework/Chista/middlewares"
	"github.com/Chista-Framework/Chista/models"
	"github.com/Chista-Framework/Chista/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load .env file
	err := godotenv.Load(".ENV")
	if err != nil {
		fmt.Printf("Error while openning ENV file. Error: %s", err)
		return
	}

	// Create a task model for perodic functions.
	tasks := map[string]models.PeriodicFunctions{
		"Ransomware Profiles Data Check": {Fn: controller.GetRansomProfileData, Interval: 24 * time.Hour},
		"Ransom Data Check":              {Fn: controller.GetRansomwatchData, Interval: 24 * time.Hour},
		"Apt Profiles Data Check":        {Fn: controller.ScheduleAptData, Interval: 6 * time.Hour},
		"Phishing Monitor Task":          {Fn: controller.PeriodicPhishingMonitorTask, Interval: 250 * time.Second},
	}

	// Create a separate channel and quit signal for each function.
	quit := make(chan struct{})

	// Running periodic functions using goroutines
	helpers.RunPeriodicly(tasks, quit)

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
	routes.ThreatProfileRoute(router)
	routes.BlacklistRoute(router)
	routes.SourceRoute(router)
	routes.C2Route(router)
	routes.ActivitiesRoute(router)

	// Serve the APIg
	router.Run("localhost:7777")

	// Send quit signals to stop each function from running
	close(quit)

	// Wait for the main program to finish to allow all periodic functions to complete.
	time.Sleep(2 * time.Second)
}
