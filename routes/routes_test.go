package routes

import (
	"bytes"
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var urlFilePath string = ""

func init() {
	dir, _ := homedir.Dir()
	urlFilePath = filepath.Join(dir, "urls_test")
	err := os.MkdirAll(urlFilePath, 0744)
	if err != nil {
		log.Fatalf("Failed to create directory %v", err)
	}
}
func TestShortenUrl(t *testing.T) {
	var urlTests = []struct {
		url        string
		statusCode int
	}{
		{"https://www.youtube.com/watch?v=OVBvOuxbpHA", http.StatusOK},
		{"https://support.google/", http.StatusBadRequest},
		{"https://www.thepolyglotdeveloper.com/", http.StatusOK},
		{"https://support.goog[]p[]/", http.StatusBadRequest},
	}

	for _, url := range urlTests {
		var myurl Url
		myurl.LongUrl = url.url
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(myurl)
		if err != nil {
			log.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/shorten", &buf)
		if err != nil {
			t.Errorf("Error creating a new request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.Handler(ShortenUrl(urlFilePath))
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != url.statusCode {
			t.Errorf("Handler returned wrong status code for url %v Expected: %d. got : %d.", url.url, http.StatusOK, status)
		}
	}
}

func TestResolveUrl(t *testing.T) {
	var urlTests = []struct {
		id         string
		statusCode int
	}{
		{"32456", http.StatusBadRequest},
		{"", http.StatusBadRequest},
	}
	defer func() {
		err := os.Remove(filepath.Join(urlFilePath, UrlFile))
		if err != nil {
			log.Fatal(err)

		}
	}()
	for _, urlId := range urlTests {
		req, _ := http.NewRequest("GET", "/"+urlId.id, nil)
		rr := httptest.NewRecorder()
		handler := http.Handler(ResolveUrl(urlFilePath))
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != urlId.statusCode {
			t.Errorf("Handler returned wrong status code. Expected: %d. got : %d.", http.StatusOK, status)
		}
	}
}
