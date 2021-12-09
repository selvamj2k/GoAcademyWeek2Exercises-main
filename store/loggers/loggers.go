package loggers

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	StoreLogger *log.Logger
)

func SetupLoggers() {
	file, err := os.OpenFile("htaccess.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal()
	}
	storeLogfile, err := os.OpenFile("store.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal()
	}
	InfoLogger = log.New(file, "INFO", log.Ldate|log.Ltime)
	StoreLogger = log.New(storeLogfile, "INFO", log.Ldate|log.Ltime)
}
