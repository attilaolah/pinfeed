package main

import (
    "io"
    "io/ioutil"
	"net/http"
	"os"
    "regexp"
    "strings"
)

var thumb = regexp.MustCompile("https?://media-cache-[0-9a-z]+.pinimg.com/192x/[0-9a-f]{2}/[0-9a-f]{2}/[0-9a-f]{2}/[0-9a-f]{32}.jpg")

func main() {
	http.HandleFunc("/", pinfeed)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		panic(err)
	}
}

func pinfeed(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
        // TODO: some home page would be nice…
		w.WriteHeader(http.StatusNoContent)
		return
	}
    res, err := http.Get(feedURL(username(r.URL.Path)))
    if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        return
    }
    defer res.Body.Close()
    for key, vals := range res.Header {
        for _, val := range vals {
            w.Header().Add(key, val)
        }
    }
    w.WriteHeader(res.StatusCode)
    buf, err := replaceAllThumbs(res.Body)
    if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.Write(buf)
}

func feedURL(username string) string {
	return "https://www.pinterest.com/" + username + "/feed.rss"
}

func username(path string) string {
    return strings.SplitN(path, "/", 3)[1]
}

func replaceAllThumbs(r io.Reader) (buf []byte, err error) {
    if buf, err = ioutil.ReadAll(r); err != nil {
        return
    }
    return
}
