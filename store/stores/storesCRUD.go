package stores

import (
	"encoding/json"
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

func HandleStoreDeleteMethod(resp http.ResponseWriter, req *http.Request) {

	if strings.Contains(req.URL.String(), storeURL) == true {

		//call put method
		pkgloggers.StoreLogger.Println("Delete method called")

		//get user
		storeUser, err := pkghttp.GetHTTPStoreUser(resp, req)
		if err != nil {
			return
		}

		//get store ID from HTTP query string
		storeKey := GetHTTPStoreKey(resp, req)

		s, b := FindStore(storeKey)

		//if key not found return
		if b == false {
			pkgloggers.StoreLogger.Println("Key not found")
			resp.WriteHeader(http.StatusNotFound)
			return
		}

		//different user from the one who originally created cannot delete
		if s.User != storeUser {
			pkgloggers.StoreLogger.Printf("different user (%v) from the one who originally (%v) created cannot delete", s.User, storeUser)
			resp.WriteHeader(http.StatusForbidden)
			return
		}
		pkgloggers.StoreLogger.Printf("User trying to delete: logged in user %v", s.User)

		//we are here so ready to delete
		i, b := indexOf(storeKey)

		pkgloggers.StoreLogger.Printf("Index of store item trying to delete is this :%v ", i)
		if b == true {
			//s, err := RemoveElement(s, 2)
			Stores = RemoveIndex(i)

		}

	}
}

func HandleStorePutMethod(resp http.ResponseWriter, req *http.Request) {

	if strings.Contains(req.URL.String(), storeURL) == true {

		//call put method
		pkgloggers.StoreLogger.Println("Put method called")

		//get user
		storeUser, err := pkghttp.GetHTTPStoreUser(resp, req)
		if err != nil {
			return
		}

		//get store value from HTTP body
		storeValue := GetHTTPStoreValue(req)
		pkgloggers.StoreLogger.Printf("Value fetched %v", storeValue)

		//get store ID from HTTP query string
		storeKey := GetHTTPStoreKey(resp, req)

		//add to store
		AddToStore(resp, storeKey, storeValue, storeUser)

	}
}

func HandleStoreListMethod(resp http.ResponseWriter, _ *http.Request) {

	storesInfo := make([]StoreInfo, len(Stores))

	for i, v := range Stores {
		storesInfo[i].Key = strconv.Itoa(v.Key)
		storesInfo[i].User = v.User
	}

	jasonData, err := json.Marshal(storesInfo)
	if err != nil {
		//error
		resp.WriteHeader(http.StatusInternalServerError)
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write(jasonData)

}

func HandleStoreGetMethod(resp http.ResponseWriter, req *http.Request) {

	user, err := pkghttp.GetHTTPStoreUser(resp, req)

	if user != "" && err == nil {
		storeKey := GetHTTPStoreKey(resp, req)
		s, b := FindStore(storeKey)

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

func GetHTTPStoreKey(resp http.ResponseWriter, req *http.Request) int {
	pkgloggers.StoreLogger.Println("Get store key")
	var keyVal []string
	keyVal = strings.SplitAfter(req.URL.String(), "/")
	v := keyVal[len(keyVal)-1]
	pkgloggers.StoreLogger.Println("Store key fetched", v)

	Key, errInt := strconv.Atoi(v)

	//if invalid number return
	if errInt != nil {
		pkgloggers.StoreLogger.Println("Key not a valid number")
		resp.WriteHeader(http.StatusForbidden)
		return 0
	}

	return Key
}

func AddToStore(resp http.ResponseWriter, storeKey int, storeVal string, storeUser string) {

	pkgloggers.StoreLogger.Println("Before adding/updating to store")
	s, b := FindStore(storeKey)

	if b == false {
		store1 := NewStore(storeKey, storeVal, storeUser, 1, 1, time.Now())
		Stores = append(Stores, store1)
		pkgloggers.StoreLogger.Printf("Added store key:%v | value:%v", storeKey, storeVal)
	} else {
		if storeUser == s.User {
			s.Value = storeVal
			s.Age = time.Now()
			s.Reads = s.Reads + 1
			pkgloggers.StoreLogger.Printf("updated  store key:%v | value:%v", storeKey, storeVal)
		} else {
			resp.WriteHeader(http.StatusForbidden)
		}

	}
}
