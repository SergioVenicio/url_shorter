package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/SergioVenicio/url_shorter/database"
	"github.com/SergioVenicio/url_shorter/models"
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

	ctx := context.Background()
	rdb := database.GetClient()

	key := fmt.Sprintf("url:%s", newUrl.Id)
	value, _ := json.Marshal(newUrl)
	err = rdb.Set(ctx, key, value, 0).Err()
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

	ctx := context.Background()
	rdb := database.GetClient()
	result, err := rdb.Get(ctx, fmt.Sprintf("url:%s", id)).Result()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var url models.Url
	json.Unmarshal([]byte(result), &url)
	http.Redirect(w, r, url.Url, 303)

	result, _ = rdb.Get(ctx, fmt.Sprintf("access:%s", id)).Result()
	access, _ := strconv.Atoi(result)
	rdb.Set(ctx, fmt.Sprintf("access:%s", id), access+1, 0)
}

func GetUrl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paramId := vars["id"]
	if paramId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Id param is required!"))
		return
	}

	ctx := context.Background()
	rdb := database.GetClient()
	urlResult, _ := rdb.Get(ctx, fmt.Sprintf("url:%s", paramId)).Result()
	var url models.Url
	json.Unmarshal([]byte(urlResult), &url)

	accessResult, _ := rdb.Get(ctx, fmt.Sprintf("access:%s", paramId)).Result()
	accessCount, _ := strconv.Atoi(accessResult)
	access := models.Access{
		Url:    url,
		Access: accessCount,
	}

	json.NewEncoder(w).Encode(access)
}

func ListUrls(w http.ResponseWriter, r *http.Request) {
	var urls []models.Url
	ctx := context.Background()
	rdb := database.GetClient()
	var cursor uint64
	keys, cursor, err := rdb.Scan(
		ctx,
		cursor,
		"url:*",
		10,
	).Result()

	for _, k := range keys {
		var eachUrl models.Url
		value, _ := rdb.Get(ctx, k).Result()
		json.Unmarshal([]byte(value), &eachUrl)
		urls = append(urls, eachUrl)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(urls)
}
