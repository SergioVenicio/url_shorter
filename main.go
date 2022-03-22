package main

import (
	"net/http"

	"github.com/SergioVenicio/url_shorter/controllers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	router := mux.NewRouter()

	router.HandleFunc("/url", controllers.CreateUrl).Methods("POST")
	router.HandleFunc("/url", controllers.ListUrls).Methods("GET")
	router.HandleFunc("/url/{id}", controllers.GetUrl).Methods("GET")
	router.HandleFunc("/{url}", controllers.Redirect).Methods("GET")

	http.ListenAndServe(":5000", router)
}
