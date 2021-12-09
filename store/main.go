package main

import (
	"fmt"
	pkgloggers "loggers"
	"net/http"
	"os"
	pkgstores "stores"
	"strconv"
)

func main() {

	pkgloggers.SetupStoreLogger()
	pkgloggers.StoreLogger.Println("Read start up parameters")
	port := GetPort(os.Args)

	pkgloggers.StoreLogger.Println("Starting server")
	fmt.Println("Starting Server")

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/store/", serveStore)
	http.HandleFunc("/list", pkgstores.HandleStoreListMethod)

	//check port for binding
	pkgloggers.StoreLogger.Printf("Try connecting to the passed in port %v", port)

	//_, err := net.Dial("tcp", ConnHost+":"+  strconv.Itoa(v))
	//_, err := http.Get("http://" + ConnHost + ":" + strconv.Itoa(v))
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		pkgloggers.StoreLogger.Printf("Error in connecting to the passed in port %v", err)
		pkgloggers.StoreLogger.Printf("Unable to bind port %v. Exit code 2", port)
		os.Exit(2)
	}

}

func GetPort(args []string) string {
	pkgloggers.StoreLogger.Printf("Length of start up parameters: %v ", len(args))

	if len(args) >= 3 {
		v, err := strconv.Atoi(args[2])
		if err != nil || v <= 0 {
			pkgloggers.StoreLogger.Println("Start up parameter not valid. Exit code 1")
			os.Exit(1)
		} else {
			return args[2]
		}
	} else {
		pkgloggers.StoreLogger.Println("Start up parameter missing. Exit code 1")
		os.Exit(1)
	}

	return args[2]
}

func serveStore(resp http.ResponseWriter, req *http.Request) {

	pkgloggers.SetupInfoLogger(req)
	pkgloggers.InfoLogger.Println("Serve called")
	pkgloggers.StoreLogger.Println("Request URL", req.URL.String())

	switch req.Method {
	case http.MethodPut:
		pkgloggers.InfoLogger.Println("Store Put called")
		pkgstores.HandleStorePutMethod(resp, req)
	case http.MethodGet:
		pkgloggers.InfoLogger.Println("Store Get called")
		pkgstores.HandleStoreGetMethod(resp, req)
	case http.MethodDelete:
		pkgloggers.InfoLogger.Println("Store Delete called")
		pkgstores.HandleStoreDeleteMethod(resp, req)
	}
}

func ping(resp http.ResponseWriter, req *http.Request) {
	pkgloggers.InfoLogger.Println("Ping called")
	switch req.Method {
	case http.MethodGet:
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "text/plain")
		pong := "pong"
		resp.Write([]byte(pong))
		return
	}
}
