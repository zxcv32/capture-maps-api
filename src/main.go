package main

import (
	"net/http"
	"zxcv32/capture-maps-api/src/service"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/print", service.HandleRequest)

	http.ListenAndServe(":8090", nil)
}

func home(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "http://zxcv32.com", 301)
}
