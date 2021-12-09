package pkghttp

import (
	"errors"
	pkgloggers "loggers"
	"net/http"
)

var logins = map[string]string{
	"user_a": "passwordA",
	"user_b": "passwordB",
	"user_c": "passwordC",
	"admin":  "Password1",
}

func GetHTTPStoreUser(resp http.ResponseWriter, req *http.Request) (string, error) {

	pkgloggers.StoreLogger.Println("Get store user")

	username, password, ok := req.BasicAuth()
	if !ok {
		pkgloggers.StoreLogger.Println("Error parsing basic auth", ok)
		resp.WriteHeader(http.StatusUnauthorized) // 401
		return "", errors.New("user statusUnauthorized")
	}
	expectedPassword, ok := logins[username]
	if !ok {
		pkgloggers.StoreLogger.Println("Error Unknown username: %s", username)
		resp.WriteHeader(http.StatusUnauthorized)
		return "", errors.New("user statusUnauthorized")
	}
	if password != expectedPassword {
		pkgloggers.StoreLogger.Println("Password is incorrect: %s\n", username)
		resp.WriteHeader(http.StatusUnauthorized) // 401
		return "", errors.New("user statusUnauthorized")
	}
	pkgloggers.StoreLogger.Printf("logged in user is %v", username)

	return username, nil
}

// Handlers
