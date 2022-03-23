package main

import (
	"github.com/SergioVenicio/url_shorter/shared/http"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	http.RunServer()
}
