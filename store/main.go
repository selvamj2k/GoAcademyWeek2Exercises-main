package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)
var (
	InfoLogger *log.Logger
	StoreLogger *log.Logger
)

var logins = map[string]string{
	"user_a":     "passwordA",
	"user_b":     "passwordB",
	"user_c":     "passwordC",
	"admin": "Password1",
}

type store struct {
	key int
	value string
	user string
	writes int
	reads int
	age time.Time
}

func (s *store) String() string {
	return fmt.Sprintf("Store key %v, value %v, user %v", s.key, s.value, s.user)
}

func NewStore(key int,value string, owner string, writes int, read int, age time.Time) *store {
	store:= store{key,value,owner,writes,read,age}
	return &store
}

var stores []*store


const (
	ConnHost = "localhost"
	root="/"
	pingURL="/ping"
	storeURL="/store"
)

func main(){

	setupLoggers()
	InfoLogger.Println("Starting server")

	fmt.Println("Starting Server")
	http.HandleFunc(root,serve)
	http.ListenAndServe(":8000", nil)
	fmt.Println("Server available")

}

func setupLoggers(){
	file, err := os.OpenFile("htaccess.log",os.O_APPEND|os.O_CREATE|os.O_WRONLY,0666)

	if err!=nil{
		log.Fatal()
	}

	storeLogfile, err := os.OpenFile("store.log",os.O_APPEND|os.O_CREATE|os.O_WRONLY,0666)

	if err!=nil{
		log.Fatal()
	}

	InfoLogger=log.New(file,"INFO", log.Ldate| log.Ltime)
	StoreLogger=log.New(storeLogfile,"INFO", log.Ldate| log.Ltime)
}


func serve(resp http.ResponseWriter, req *http.Request){

	InfoLogger.Println("IP:", req.RemoteAddr, "HTTP Method:",req.Method, "URL:", req.URL.String()  )

	StoreLogger.Println("Request URL", req.URL.String())
	if req.URL.String()==pingURL {
		ping(resp,req)
	}

	go HandlePutMethod(resp,req )

}

func HandlePutMethod(resp http.ResponseWriter,req *http.Request){
	if strings.Contains(req.URL.String(),storeURL)==true {
		switch req.Method {
		case http.MethodPut:

			//call put method
			StoreLogger.Println("Put method called")

			defer req.Body.Close()

			//authenticate user
			storeUserChan, err:=GetHTTPStoreUser(resp,req)

			//only let if authenticated
			//get store value from HTTP body
			//Value:=GetHTTPStoreValue(req)
			storeValueChan := GetHTTPStoreValue(req)

			//get store ID from HTTP query string
			storeKeyChan:=GetHTTPStoreKey(req)


			AddToStore(resp, storeKeyChan, storeValueChan, storeUserChan, err)

		}
	}
}

func GetHTTPStoreUser(resp http.ResponseWriter,req *http.Request) (chan string, chan error) {
	StoreLogger.Println("Get store user")
	userChan:=make(chan string)
	errChan:=make(chan error)

	go func()  {
		username, password, ok := req.BasicAuth()
		if !ok {
			StoreLogger.Println("Error parsing basic auth", ok)
			resp.WriteHeader(http.StatusUnauthorized) // 401
			errChan<- errors.New("user statusUnauthorized")
		}
		expectedPassword, ok := logins[username]
		if !ok {
			StoreLogger.Println("Error Unknown username: %s", username)
			resp.WriteHeader(http.StatusUnauthorized)
			errChan<- errors.New("user statusUnauthorized")
		}
		if password != expectedPassword {
			StoreLogger.Println("Password is incorrect: %s\n", username)
			resp.WriteHeader(http.StatusUnauthorized) // 401
			errChan<- errors.New("user statusUnauthorized")
		}
		StoreLogger.Printf("logged in user is %v", username)
		userChan<-username
		errChan<- nil
		}()

	return userChan, errChan
}

func GetHTTPStoreValue(req *http.Request) chan string {
	StoreLogger.Println("Get store value")
	valueChan:=make(chan string)

	go func() {
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		var v string
		json.Unmarshal(bodyBytes, &v )
		StoreLogger.Println("Store value fetched", v)
		valueChan<-v
	}()

	return valueChan
}

func GetHTTPStoreKey(req *http.Request) chan string {
	StoreLogger.Println("Get store key")

	idChan:=make(chan string)

	go func(){
		var keyVal []string
		keyVal=strings.SplitAfter(req.URL.String(),"/")

		v :=  keyVal[len(keyVal)-1]
		StoreLogger.Println("Store key fetched", v)
		idChan<-v
	}()

	return idChan
}

func AddToStore(resp http.ResponseWriter, storeKey chan string, storeVal chan string, storeUser chan string,  storeErr chan error) {

	//Key, err:=strconv.Atoi(<-storeKey)
	StoreLogger.Println("In method add store")

	//get all store values here
	Value := <-storeVal
	User := <-storeUser
	Key, errInt := strconv.Atoi(<-storeKey)
	err := <-storeErr

	if errInt != nil {
		resp.WriteHeader(http.StatusForbidden)
	}

	if err == nil {
		StoreLogger.Println("Before adding/updating to store")
		s, b := FindStore(Key)

		if b == false {
			store1 := NewStore(Key, Value, User, 1, 1, time.Now())
			stores = append(stores, store1)
			StoreLogger.Printf("Added store key:%v | value:%v", Key, Value)
		}else{
			if User==s.user{
				s.value=Value
				s.age=time.Now()
				s.reads=s.reads+1
				StoreLogger.Printf("updated  store key:%v | value:%v", Key, Value)
			}else{
				resp.WriteHeader(http.StatusForbidden)
			}
		}

	}
}



func FindStore(Key int) (*store, bool) {
	for i:=0; i<= len(stores)-1; i++{
		store:=stores[0]
		if store.key==Key {
			StoreLogger.Printf("Store found for the passed in key %v. update store", Key)
			return store, true
		}
	}
	StoreLogger.Printf("Store not found for the passed in key %v. Add new", Key)
	return _, false
}

func ping(resp http.ResponseWriter,req *http.Request){
	switch req.Method {
	case http.MethodGet:
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "text/plain")
		jasonData, err := json.Marshal([]string{"pong"})
		if err != nil {
			http.Error(resp, "Error", http.StatusInternalServerError)
			return
		}
		resp.Write(jasonData)
		return
	}



