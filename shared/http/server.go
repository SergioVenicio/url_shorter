package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/SergioVenicio/url_shorter/useCases/url/controllers"
	"github.com/gorilla/mux"
)

func RunServer() {
	router := mux.NewRouter()

	router.HandleFunc("/url", controllers.CreateUrl).Methods("POST")
	router.HandleFunc("/url", controllers.ListUrls).Methods("GET")
	router.HandleFunc("/url/{id}", controllers.GetUrl).Methods("GET")
	router.HandleFunc("/{url}", controllers.Redirect).Methods("GET")

	server := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	fmt.Println("Running server on", server)
	http.ListenAndServe(server, router)
}
