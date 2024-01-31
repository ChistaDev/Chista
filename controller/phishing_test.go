package controller

import (
	"fmt"
	"testing"

	"github.com/Chista-Framework/Chista/helpers"
	"github.com/Chista-Framework/Chista/logger"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func initWsConnection() error {
	API_ONLY := "true"

	// Establish a WebSocket connection.
	//var err error
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:7778/ws", nil)
	if err != nil {
		logger.Log.Warnf("WebSocket connection error: %v. Checking for API_ONLY parameter...", err)
		logger.Log.Warnf("API_ONLY: %s.", API_ONLY)
		if !(API_ONLY == "true") {
			logger.Log.Debugf("API_ONLY in if block: %s.", API_ONLY)
			logger.Log.Errorf("WebSocket connection error: %v", err)
			return err
		}
	}

	helpers.CONN = conn
	return nil
}

func TestGetCrtshCTLogs(t *testing.T) {

	a := assert.New(t)
	input := "xn--mcrosoft-tkb.com"

	// Initiliaze WS connection to test
	initWsConnection()

	crtsh_response, err := GetDomainsFromCrtshCTLogs(input)
	if err != nil {
		a.Error(err)
	}

	fmt.Printf("\t[T] Converted domain: %v\n", crtsh_response)
	a.NotNil(crtsh_response, "Converted domain: %v", crtsh_response)
}
