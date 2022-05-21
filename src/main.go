package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"zxcv32/capture-maps-api/src/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Warningln("No .env file found")
	}
	apiKey := os.Getenv("API_KEY")
	if len(apiKey) == 0 {
		log.Fatalln("Google Maps API Key not set")
	}
	http.HandleFunc("/", home)
	http.HandleFunc("/print", service.PrintHandler(apiKey))

	http.ListenAndServe(":8090", nil)
}

func home(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "http://zxcv32.com", 301)
}
