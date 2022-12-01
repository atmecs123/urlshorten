package routes

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	Port    = "9000"
	UrlFile = "urls.json"
)

type Url struct {
	Id       string `json:"id"`
	LongUrl  string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}

// ShortenUrl is handler to shorten the long url and in response sends id,longUrl,shortUrl
func ShortenUrl(urlPath string) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		filepath := filepath.Join(urlPath, UrlFile)
		var reqUrl Url
		err := json.NewDecoder(r.Body).Decode(&reqUrl)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Unable to parse the request")
			return
		}
		if reqUrl.LongUrl == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Empty url in the request")
			return
		}
		if !govalidator.IsURL(reqUrl.LongUrl) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Not a valid url")
			return
		}
		_, err = http.Get(reqUrl.LongUrl)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Not a valid url. Please check the url again")
			return
		}
		var urls []Url
		_, err = os.Stat(filepath)
		// create file if not exists
		if os.IsNotExist(err) {
			var file, err = os.Create(filepath)
			if err != nil {
				log.Fatalf("Unable to create a file %v", err)
			}
			defer file.Close()
		}
		file, err := os.Stat(filepath)
		if err != nil {
			log.Printf("Unable to stat a file %v. \n  Creating the file path......", err)
		}
		if file.Size() != 0 {
			file, _ := ioutil.ReadFile(filepath)
			_ = json.Unmarshal([]byte(file), &urls)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Fatalf("Unable to unmarshal %v", err)
				return
			}
		}
		var respUrl Url
		for _, url := range urls {
			if reqUrl.LongUrl == url.LongUrl {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Short url already exists " + url.ShortUrl)
				return
			}
		}

		var id string
		id = uuid.New().String()[:6]
		shortUrl := "http://localhost:" + Port + "/" + id
		respUrl.Id = id
		respUrl.LongUrl = reqUrl.LongUrl
		respUrl.ShortUrl = shortUrl
		urls = append(urls, respUrl)
		data, err := json.MarshalIndent(urls, "", " ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("Unable to marshal the urls %v", err)
		}
		err = ioutil.WriteFile(filepath, data, 0644)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("Unable to write to the file %v", err)
		}
		json.NewEncoder(w).Encode(respUrl)

	}
	return http.HandlerFunc(handleFunc)
}

// ResolveUrl handler resolves the short or custom url to actual long url
func ResolveUrl(urlPath string) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		filepath := filepath.Join(urlPath, UrlFile)
		id := strings.TrimPrefix(r.URL.Path, "/")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("No short url id in path")
			return
		}
		var urls []Url
		file, err := ioutil.ReadFile(filepath)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Fatalf("Unable to read the urls file %v", err)
			return
		}
		err = json.Unmarshal([]byte(file), &urls)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Fatalf("Unable to unmarshal urls %v", err)
			return
		}

		var urlFound bool
		for _, url := range urls {
			if url.Id == id {
				urlFound = true
				http.Redirect(w, r, url.LongUrl, 302)
				return
			}
		}
		if !urlFound {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Invalid short url")
		}
	}
	return http.HandlerFunc(handleFunc)
}
