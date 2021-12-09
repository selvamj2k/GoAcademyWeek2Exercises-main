package loggers

import (
	"log"
	"net/http"
	"os"
)

var (
	InfoLogger  *log.Logger
	StoreLogger *log.Logger
)

func SetupStoreLogger() {

	storeLogfile, err1 := os.OpenFile("store.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err1 != nil {
		log.Fatal()
	}

	StoreLogger = log.New(storeLogfile, "INFO", log.Ldate|log.Ltime)
}

func SetupInfoLogger(req *http.Request) {
	file, err := os.OpenFile("htaccess.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal()
	}

	InfoLogger = log.New(file, "INFO "+req.RemoteAddr+":"+req.Method+":"+req.URL.String(), log.Ldate)
}
