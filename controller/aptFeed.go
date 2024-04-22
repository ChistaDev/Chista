package controller

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Chista-Framework/Chista/database"
	"github.com/Chista-Framework/Chista/helpers"
	"github.com/Chista-Framework/Chista/logger"
	"github.com/Chista-Framework/Chista/models"
	"github.com/gin-gonic/gin"
)

var (
	config = &models.APTFeedConfigDB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
)

func GetAptFeed(ctx *gin.Context) {
	err := helpers.InitiliazeWebSocketConnection()
	if err != nil {
		logger.Log.Errorf("Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true")
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true"})
		return
	}
	time.Sleep(1 * time.Second)
	defer helpers.CloseWSConnection()

	// Checking the verbosity condition.
	if ctx.Query("verbosity") != "" {
		var err error
		verbosity, err = strconv.Atoi(ctx.Query("verbosity"))
		if err != nil {
			logger.Log.Errorf("Cannot parse verbosity level %v", err)
			helpers.SendMessageWS("Source", "Cannot parse verbosity level, default using (No Verbosity).", "error")
		} else {
			logger.Log.Infof("Verbosity level %v", verbosity)
			helpers.SendMessageWS("Source", fmt.Sprintf("Verbosity level %v", verbosity), "info")
		}
	} else {
		verbosity = 0
		helpers.SendMessageWS("Source", "Using default verbosity option. (No Verbose)", "info")
	}
	helpers.VERBOSITY = verbosity

	// Connect to the database
	db, err := database.ConnectDB(config)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to connect to the database"})
		logger.Log.Errorf("Failed to connect to the database: %v", err)
		helpers.SendMessageWS("APT Feed", "Failed to connect to the database", "error")
	}

	fmt.Println(db.Name())

}
