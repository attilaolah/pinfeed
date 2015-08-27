package main

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"strings"
	"sync"
)

var pools = struct {
	gzip, deflate sync.Pool
}{
	gzip: sync.Pool{
		New: func() interface{} {
			return gzip.NewWriter(nil)
		},
	},
	deflate: sync.Pool{
		New: func() interface{} {
			return zlib.NewWriter(nil)
		},
	},
}

type compressor struct {
	http.ResponseWriter
	w encoder
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
			if enc != "gzip" && enc != "deflate" {
				continue
			}
			w.Header().Set("Content-Encoding", enc)
			c := compressor{ResponseWriter: w}
			switch enc {
			case "gzip":
				c.w = pools.gzip.Get().(encoder)
				defer pools.gzip.Put(c.w)
			case "deflate":
				c.w = pools.deflate.Get().(encoder)
				defer pools.deflate.Put(c.w)
			}
			c.w.Reset(w)
			defer c.w.Flush()
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
