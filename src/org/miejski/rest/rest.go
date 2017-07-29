package rest

import (
	"net/http"
	"fmt"
)

func methodCheck(method string, handler func(w http.ResponseWriter, request *http.Request)) func(w http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		if method == request.Method {
			handler(w,request)
		} else {
			error_msg := fmt.Sprintf("Method not supported: %s, required: %s", request.Method, method)
			http.Error(w, error_msg, http.StatusMethodNotAllowed)
		}
	}
}

func POST(fn func(w http.ResponseWriter, request *http.Request)) func(w http.ResponseWriter, request *http.Request) {
	return methodCheck("POST", fn )
}

func GET(fn func(w http.ResponseWriter, request *http.Request)) func(w http.ResponseWriter, request *http.Request) {
	return methodCheck("GET", fn )
}

