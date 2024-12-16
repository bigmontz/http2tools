package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func startEchoServer(listeningAddress string) error {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(os.Stdout, "Hello, %v, http: %v\n", r.URL.Path, r.TLS == nil)
		defer r.Body.Close()

		b := make([]byte, 1024)

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
	return startServer(listeningAddress, handlerFunc)
}

func startTcpProxyServer(listeningAddress string, targetAddress string, useTls bool) error {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(os.Stdout, "Hello, %v, http: %v\n", r.URL.Path, r.TLS == nil)
		w.WriteHeader(200)
		defer r.Body.Close()

		conn, err := dial(targetAddress, useTls)

		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "%s", err)
			fmt.Println("Error:", err)
			return
		}
		defer conn.Close()

		go func() {
			b2 := make([]byte, 1024)
			n2, err2 := conn.Read(b2)

			for ; err2 == nil; n2, err2 = conn.Read(b2) {
				if n2 == 0 {
					continue
				}
				if !writeResponse(w, b2, n2) {
					return
				}
			}
		}()

		b := make([]byte, 1024)

		fmt.Println("reading")

		n, err := r.Body.Read(b)

		for ; err == nil; n, err = r.Body.Read(b) {
			fmt.Fprintf(os.Stdout, "Bytes received %d \n", n)
			if n == 0 {
				continue
			}

			fmt.Fprintf(os.Stdout, "Bytes received %d \n", n)

			if _, err = conn.Write(b[:n]); err != nil {
				fmt.Fprintf(os.Stderr, "Oops. An error while writing to tcp server '%s'\n", err)
				return
			}
		}

		if err != io.EOF {
			fmt.Fprintf(os.Stderr, "Oops. An error while reading stream '%s'\n", err)
		}
	})
	return startServer(listeningAddress, handlerFunc)
}

func dial(targetAddress string, useTls bool) (net.Conn, error) {
	if useTls {
		return tls.Dial("tcp", targetAddress, &tls.Config{InsecureSkipVerify: true})
	}

	return net.Dial("tcp", targetAddress)
}

func startServer(listeningAddress string, handlerFunc http.HandlerFunc) error {
	http2server := &http2.Server{}

	httpServer := &http.Server{
		Addr:    listeningAddress,
		Handler: h2c.NewHandler(handlerFunc, http2server),
	}

	if err := http2.ConfigureServer(httpServer, http2server); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Server starting to listening at %s\n", listeningAddress)

	return httpServer.ListenAndServe()
}

func startTcpEchoServer(listeningAddress string) error {
	listener, err := net.Listen("tcp", listeningAddress)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Fprintln(os.Stdout, "Server is listening at", listeningAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			continue
		}

		go echoHandlerTcp(conn)
	}
}

func echoHandlerTcp(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error:", err)
			}
			return
		}

		_, err = conn.Write(buffer[:n])
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error:", err)
			}
			return
		}
	}

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
