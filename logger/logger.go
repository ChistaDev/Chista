package logger

import (
	"io"
	logging "log"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

var (
	Log       *logrus.Logger // share will all packages
	Dump_Mode bool
)

// Load the .env file and get the content's of key. Return the content's of the key.
func GoDotEnvVariable(key string) string {
	return os.Getenv(key)
}

// Initiliaze the package before main function
func init() {
	// The file needs to exist prior
	f, err := os.OpenFile("chista.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		// Use go's logger, while we configure Logrus
		logging.Fatalf("error opening file: %v", err)
	}

	// Configure Logrus
	Log = logrus.New()
	//Log.Formatter = &logrus.JSONFormatter{}
	Log.Formatter = &logrus.TextFormatter{}
	Log.SetReportCaller(true)
	mw := io.MultiWriter(os.Stdout, f)
	Log.SetOutput(mw)

	// Set the log level as Debug to produce useful info to client.
	Log.Level = logrus.TraceLevel

	dump_mode := GoDotEnvVariable("DUMP_MODE")
	dump_mode_bool, _ := strconv.ParseBool(dump_mode)
	Dump_Mode = dump_mode_bool
	Log.Infoln("DUMP_MODE=", dump_mode_bool)
}
