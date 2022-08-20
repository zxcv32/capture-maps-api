package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zxcv32/capture-maps-api/internal/service"
)

func main() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.BindEnv("google-maps.api_key", "GOOGLE_MAPS_API_KEY")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	apiKey := viper.GetString("google-maps.api_key")
	if len(apiKey) == 0 {
		log.Fatalln("Google Maps API Key not set")
	}
	http.HandleFunc("/", Home)
	http.HandleFunc("/print", service.PrintHandler(apiKey))

	http.ListenAndServe(":"+viper.GetString("server.port"), nil)
}

func Home(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, viper.GetString("server.redirect"), http.StatusSeeOther)
}
