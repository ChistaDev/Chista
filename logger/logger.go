package logger

import (
	"io"
	logging "log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	Log       *logrus.Logger // share will all packages
	Dump_Mode bool
)

func GoDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		Log.Fatalf("Error while openning ENV file.")
	}

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
	Log.Formatter = &logrus.TextFormatter{}
	Log.SetReportCaller(true)
	mw := io.MultiWriter(os.Stdout, f)
	Log.SetOutput(mw)

	// Set the log level from .env file
	debug_mode := GoDotEnvVariable("DEBUG_MODE")
	debug_mode_bool, _ := strconv.ParseBool(debug_mode)
	if debug_mode_bool == true {
		Log.Level = logrus.DebugLevel
		Log.Infoln("DEBUG_MODE=" + debug_mode)
	} else {
		Log.Level = logrus.InfoLevel
		Log.Infoln("DEBUG_MODE=" + debug_mode)
	}

	dump_mode := GoDotEnvVariable("DUMP_MODE")
	dump_mode_bool, _ := strconv.ParseBool(dump_mode)
	Dump_Mode = dump_mode_bool
	Log.Infoln("DUMP_MODE=", dump_mode_bool)
}
