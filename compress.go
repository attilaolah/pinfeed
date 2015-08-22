package main

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"strings"
)

type compressor struct {
	http.ResponseWriter
	w interface {
		io.Writer
		io.Closer
	}
}

// Compress enables gzip and deflate compression for outgoing requests.
func compress(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", strings.Trim(w.Header().Get("Vary")+",Accept-Encoding", ","))

		for _, enc := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
			enc = strings.TrimSpace(enc)
			if enc != "gzip" && enc != "deflate" {
				continue
			}
			w.Header().Set("Content-Encoding", enc)
			c := compressor{ResponseWriter: w}
			switch enc {
			case "gzip":
				c.w = gzip.NewWriter(w)
			case "deflate":
				c.w = zlib.NewWriter(w)
			}
			defer c.w.Close()
			w = c
			break
		}
		next(w, r)
	}
}

// Write calls io.Writer.Write().
func (c compressor) Write(b []byte) (int, error) {
	return c.w.Write(b)
}
