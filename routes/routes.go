package routes

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	Port = "9000"
)

type Url struct {
	Id       string `json:"id"`
	LongUrl  string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}

var UrlMap = make(map[string]Url)

// ShortenUrl is handler to shorten the long url and in response sends id,longUrl,shortUrl
func ShortenUrl(w http.ResponseWriter, r *http.Request) {
	var reqUrl Url
	err := json.NewDecoder(r.Body).Decode(&reqUrl)
	if err != nil {
		log.Fatal("Unable to decode the url request", err)
		return
	}
	if !govalidator.IsURL(reqUrl.LongUrl) {
		json.NewEncoder(w).Encode("Not a valid url")
		return
	}
	var respUrl Url
	shortUrl, ok := UrlMap[reqUrl.LongUrl]
	if ok {
		json.NewEncoder(w).Encode("Short url already exists " + shortUrl.ShortUrl)
		return
	} else {
		var id string
		id = uuid.New().String()[:6]
		newShort := "http://localhost:" + Port + "/" + id
		respUrl.Id = id
		respUrl.LongUrl = reqUrl.LongUrl
		respUrl.ShortUrl = newShort
		UrlMap[reqUrl.LongUrl] = respUrl
		json.NewEncoder(w).Encode(respUrl)
	}
}

// ResolveUrl handler resolves the short or custom url to actual long url
func ResolveUrl(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		json.NewEncoder(w).Encode("No short url id passed")
		return
	}
	var urlFound bool
	for long, short := range UrlMap {
		if short.Id == id {
			urlFound = true
			http.Redirect(w, r, long, 401)
			return
		}
	}
	if !urlFound {
		json.NewEncoder(w).Encode("Invalid short url passed")
	}
}
