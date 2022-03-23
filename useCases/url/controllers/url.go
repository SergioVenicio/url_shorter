package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/SergioVenicio/url_shorter/useCases/url/models"
	"github.com/SergioVenicio/url_shorter/useCases/url/repositories"
	"github.com/gorilla/mux"
)

func CreateUrl(w http.ResponseWriter, r *http.Request) {
	var newUrl models.Url
	json.NewDecoder(r.Body).Decode(&newUrl)

	err := newUrl.SetShorted()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = repositories.Save(&newUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(newUrl)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paramUrl := vars["url"]
	if paramUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("URL param is required!"))
		return
	}

	urlParts := strings.Split(paramUrl, "/")
	id := urlParts[len(urlParts)-1]
	url, err := repositories.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.Url, 303)
	repositories.IncrementAccess(id)
}

func GetUrl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Id param is required!"))
		return
	}

	url, err := repositories.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	access, _ := repositories.GetAccess(id)
	access.Url = url
	json.NewEncoder(w).Encode(access)
}

func ListUrls(w http.ResponseWriter, r *http.Request) {
	qOffset := r.URL.Query().Get("offset")
	qLimit := r.URL.Query().Get("limit")
	offset, _ := strconv.Atoi(qOffset)
	limit, _ := strconv.Atoi(qLimit)
	urls, err := repositories.List(int64(offset), int64(limit))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(urls)
}
