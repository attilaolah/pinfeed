package main

import (
    "io"
	"net/http"
	"os"
    "strings"
)

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
    io.Copy(w, res.Body)
}

func feedURL(username string) string {
	return "https://www.pinterest.com/" + username + "/feed.rss"
}

func username(path string) string {
    return strings.SplitN(path, "/", 3)[1]
}
