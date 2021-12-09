package stores

import (
	"io/ioutil"
	pkgloggers "loggers"
	"net/http"
	"pkghttp"
	"strconv"
	"strings"
	"time"
)

const (
	storeURL = "/store"
)

func HandleStorePutMethod(resp http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.String(), storeURL) == true {

		//call put method
		pkgloggers.StoreLogger.Println("Put method called")
		defer req.Body.Close()
		//authenticate user
		storeUser, err := pkghttp.GetHTTPStoreUser(resp, req)
		if err != nil {
			return
		}
		//only let if authenticated
		//get store value from HTTP body

		storeValue := GetHTTPStoreValue(req)
		pkgloggers.StoreLogger.Printf("Value fetched %v", storeValue)
		//get store ID from HTTP query string
		storeKey := GetHTTPStoreKey(req)

		AddToStore(resp, storeKey, storeValue, storeUser)

	}
}

func HandleStoreGetMethod(resp http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()

	user, err := pkghttp.GetHTTPStoreUser(resp, req)

	if user != "" && err == nil {
		storeKey := GetHTTPStoreKey(req)
		Key, errInt := strconv.Atoi(storeKey)

		if errInt != nil {
			pkgloggers.StoreLogger.Println("Key not a valid number")
			resp.WriteHeader(http.StatusForbidden)
			return
		}
		s, b := FindStore(Key)

		if b == false {
			pkgloggers.StoreLogger.Println("Key not found")
			resp.WriteHeader(http.StatusNotFound)
			return
		} else {
			resp.WriteHeader(http.StatusOK)
			resp.Header().Set("Content-Type", "text/plain")
			resp.Write([]byte(s.Value))
		}
	}
}

func GetHTTPStoreValue(req *http.Request) string {
	pkgloggers.StoreLogger.Println("Get store value")
	bodyBytes, _ := ioutil.ReadAll(req.Body)
	v := string(bodyBytes)
	pkgloggers.StoreLogger.Printf("Store body bytes %v", string(bodyBytes))
	pkgloggers.StoreLogger.Printf("Store value fetched %v", v)
	return v
}

func GetHTTPStoreKey(req *http.Request) string {
	pkgloggers.StoreLogger.Println("Get store key")
	var keyVal []string
	keyVal = strings.SplitAfter(req.URL.String(), "/")
	v := keyVal[len(keyVal)-1]
	pkgloggers.StoreLogger.Println("Store key fetched", v)
	return v
}

func AddToStore(resp http.ResponseWriter, storeKey string, storeVal string, storeUser string) {

	pkgloggers.StoreLogger.Println("Add store function called")

	//get all store values here
	Key, errInt := strconv.Atoi(storeKey)

	if errInt != nil {
		resp.WriteHeader(http.StatusForbidden)
		return
	}

	pkgloggers.StoreLogger.Println("Before adding/updating to store")
	s, b := FindStore(Key)

	if b == false {
		store1 := NewStore(Key, storeVal, storeUser, 1, 1, time.Now())
		Stores = append(Stores, store1)
		pkgloggers.StoreLogger.Printf("Added store key:%v | value:%v", Key, storeVal)
	} else {
		if storeUser == s.User {
			s.Value = storeVal
			s.Age = time.Now()
			s.Reads = s.Reads + 1
			pkgloggers.StoreLogger.Printf("updated  store key:%v | value:%v", Key, storeVal)
		} else {
			resp.WriteHeader(http.StatusForbidden)
		}

	}
}
