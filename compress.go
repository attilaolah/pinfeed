package main

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"strings"
	"sync"
)

var pools = map[string]sync.Pool{
	"gzip": {
		New: func() interface{} {
			return gzip.NewWriter(nil)
		},
	},
	"deflate": {
		New: func() interface{} {
			return zlib.NewWriter(nil)
		},
	},
}

type compressor struct {
	http.ResponseWriter
	encoder encoder
}

type encoder interface {
	io.Writer
	Reset(io.Writer)
	Flush() error
}

// compress enables gzip and deflate compression for outgoing requests.
func compress(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", strings.Trim(w.Header().Get("Vary")+",Accept-Encoding", ","))

		for _, enc := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
			enc = strings.TrimSpace(enc)

			pool, ok := pools[enc]
			if !ok { // no such encoding
				continue
			}

			w.Header().Set("Content-Encoding", enc)

			zw := pool.Get().(encoder)
			defer pool.Put(zw)
			defer zw.Flush()
			zw.Reset(w)

			w = &compressor{
				ResponseWriter: w,
				encoder:        zw,
			}
			break
		}
		next(w, r)
	}
}

// Write calls io.Writer.Write().
func (c *compressor) Write(b []byte) (int, error) {
	return c.encoder.Write(b)
}

func decodeBody(r *http.Response) error {
	// Decode the response:
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		body, err := gzip.NewReader(r.Body)
		if err != nil {
			return err
		}
		r.Body = body
	case "deflate":
		body, err := zlib.NewReader(r.Body)
		if err != nil {
			return err
		}
		r.Body = body
	}

	return nil
}
