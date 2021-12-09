package main

import (
	"fmt"
	pkgloggers "loggers"
	"net/http"
	"os"
	pkgstores "stores"
	"strconv"
)

const (
	root     = "/"
	pingURL  = "/ping"
	ConnHost = "localhost"
)

var logins = map[string]string{
	"user_a": "passwordA",
	"user_b": "passwordB",
	"user_c": "passwordC",
	"admin":  "Password1",
}

func main() {

	pkgloggers.SetupLoggers()
	pkgloggers.InfoLogger.Println("Read start up parameters")
	pkgloggers.InfoLogger.Printf("Length of start up parameters: %v ", len(os.Args))

	if len(os.Args) >= 3 {
		v, err := strconv.Atoi(os.Args[2])
		if err != nil || v <= 0 {
			pkgloggers.InfoLogger.Println("Start up parameter not valid. Exit code 1")
			os.Exit(1)
		} else {
			//check port for binding
			pkgloggers.InfoLogger.Printf("Port number fetched %v", v)
			pkgloggers.InfoLogger.Printf("Try connecting to the passed in port %v", v)

			//_, err := net.Dial("tcp", ConnHost+":"+  strconv.Itoa(v))
			_, err := http.Get("http://" + ConnHost + ":" + strconv.Itoa(v))

			pkgloggers.InfoLogger.Printf("Error in connecting to the passed in port %v", err)
			if err != nil {
				pkgloggers.InfoLogger.Printf("Unable to bind port %v. Exit code 2", v)
				os.Exit(2)
			}
		}
	} else {
		pkgloggers.InfoLogger.Println("Start up parameter missing. Exit code 1")
		os.Exit(1)
	}

	pkgloggers.InfoLogger.Println("Starting server")
	fmt.Println("Starting Server")

	http.HandleFunc(root, serve)
	//http.ListenAndServe(":8000", nil)
	http.ListenAndServe(":"+os.Args[3], nil)

	fmt.Println("Server available")

}

func serve(resp http.ResponseWriter, req *http.Request) {

	pkgloggers.InfoLogger.Println("IP:", req.RemoteAddr, "HTTP Method:", req.Method, "URL:", req.URL.String())

	pkgloggers.StoreLogger.Println("Request URL", req.URL.String())
	if req.URL.String() == pingURL {
		ping(resp, req)
	}

	switch req.Method {
	case http.MethodPut:
		pkgstores.HandleStorePutMethod(resp, req)
	case http.MethodGet:
		pkgstores.HandleStoreGetMethod(resp, req)
	}

}

func ping(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "text/plain")
		//jasonData, err := json.Marshal([]string{"pong"})
		//if err != nil {
		//	http.Error(resp, "Error", http.StatusInternalServerError)
		//	return
		//}
		//resp.Write(jasonData)
		pong := "pong"
		resp.Write([]byte(pong))
		return
	}
}
