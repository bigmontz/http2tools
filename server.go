package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func startEchoServer(listeningAddress string) error {
	http2server := &http2.Server{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(os.Stdout, "Hello, %v, http: %v\n", r.URL.Path, r.TLS == nil)
		defer r.Body.Close()

		b := make([]byte, 2)

		n, err := r.Body.Read(b)

		for ; err == nil; n, err = r.Body.Read(b) {
			if n == 0 {
				continue
			}

			if !writeResponse(w, b, n) {
				return
			}
		}

		if err != io.EOF {
			fmt.Fprintf(os.Stderr, "Oops. An error while reading stream '%s'\n", err)
		}
	})

	httpServer := &http.Server{
		Addr:    listeningAddress,
		Handler: h2c.NewHandler(handler, http2server),
	}

	if err := http2.ConfigureServer(httpServer, http2server); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Server starting to listening at %s\n", listeningAddress)

	return httpServer.ListenAndServe()
}

func writeResponse(w http.ResponseWriter, b []byte, n int) bool {
	if _, wErr := w.Write(b[0:n]); wErr != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while writing stream '%s'\n", wErr)
		return false
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	return true
}
