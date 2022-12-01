package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	fmt.Println("dir is", dir)
	urlFilePath = filepath.Join(dir, "urls_test")
	err := os.MkdirAll(urlFilePath, 0744)
	if err != nil {
		log.Fatalf("Failed to create directory %v", err)
	}
}
func TestShortenUrl(t *testing.T) {
	var urlTests = []struct {
		url            string // input
		expectedResult string // string

	}{
		{"https://www.youtube.com/watch?v=OVBvOuxbpHA", ""},
		{"https://www.thepolyglotdeveloper.com/", ""},
		{"https://support.google/", "Not a valid url"},
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
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code. Expected: %d. got : %d.", http.StatusOK, status)
		}
	}
}

func TestResolveUrl(t *testing.T) {
	var urlTests = []struct {
		id             string // input
		expectedResult string // string

	}{
		{"32456", "Invalid short url passed"},
		{"", "No short url id passed"},
	}

	for _, urlId := range urlTests {
		req, _ := http.NewRequest("GET", "/"+urlId.id, nil)
		fmt.Println("my req is", req.URL)
		rr := httptest.NewRecorder()
		handler := http.Handler(ResolveUrl(urlFilePath))
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code. Expected: %d. got : %d.", http.StatusOK, status)
		}

		if urlId.expectedResult != rr.Body.String() {
			t.Errorf("Handler returned wrong response error. Expected %s and got %s", urlId.expectedResult, rr.Body.String())
		}
	}
}
