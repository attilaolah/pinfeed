package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	origin  = "https://www.pinterest.com/"
	repoURL = "https://github.com/attilaolah/pinfeed"
)

var (
	thumb       = regexp.MustCompile("\\b(https?://[0-9a-z-]+.pinimg.com/)(\\d+x)(/[/0-9a-f]+.jpg)\\b")
	replacement = []byte("${1}1200x${3}")
	headers     = []string{
		// Cache control headers
		"Age",
		"Cache-Control",
		"Content-Type",
		"Date",
		"Etag",
		"Last-Modified",
		"Vary",
		// Pinterest-specific stuff
		"Pinterest-Breed",
		"Pinterest-Generated-By",
		"Pinterest-Version",
	}
)

func pinFeed(w http.ResponseWriter, r *http.Request) {
	// Home page:
	if r.URL.Path == "/" {
		http.Redirect(w, r, repoURL, http.StatusMovedPermanently)
		return
	}

	// Feed pages:
	req, err := http.NewRequest(r.Method, feedURL(r.URL.Path), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Pass along HTTP headers to Pinterest:
	for key, vals := range r.Header {
		for _, val := range vals {
			req.Header.Add(key, val)
		}
	}

	// Make an HTTP request:
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Copy white-listed headers to the response:
	for _, key := range headers {
		if val := res.Header.Get(key); val != "" {
			w.Header().Set(key, val)
		}
	}
	w.WriteHeader(res.StatusCode)

	// Write modified response:
	buf, err := replaceThumbs(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(buf)
}

func feedURL(path string) string {
	username, feed := userAndFeed(path)
	if feed == "" {
		feed = "feed"
	}
	return origin + "/" + username + "/" + feed + ".rss"
}

func userAndFeed(path string) (username, feed string) {
	path = strings.TrimSuffix(path, ".rss")
	parts := strings.SplitN(path, "/", 4)
	if len(parts) > 1 {
		username = parts[1]
	}
	if len(parts) > 2 {
		feed = parts[2]
	}
	return
}

func replaceThumbs(r io.Reader) (buf []byte, err error) {
	if buf, err = ioutil.ReadAll(r); err == nil {
		buf = thumb.ReplaceAll(buf, replacement)
	}
	return
}
