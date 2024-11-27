package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/http2"
)

func connectionUsingRandomClient(address string, numberOfBytes int, batchSize int, timeBetweenBatches int) error {
	fmt.Fprintf(os.Stdout, "%d bytes will be send to %s \n", numberOfBytes, address)

	client := http.Client{
		Transport: &http2.Transport{
			// So http2.Transport doesn't complain the URL scheme isn't 'https'
			AllowHTTP: true,
			// Pretend we are dialing a TLS endpoint.
			// Note, we ignore the passed tls.Config
			DialTLSContext: func(ctx context.Context, n, a string, _ *tls.Config) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, n, a)
			},
		},
	}

	r, w := io.Pipe()

	go func() {
		defer w.Close()
		for i := 0; i < numberOfBytes || numberOfBytes < 0; {
			size := batchSize
			if size+i > numberOfBytes && numberOfBytes > -1 {
				size = numberOfBytes - i
			}

			out := make([]byte, size)

			rand.Read(out)

			i += size

			fmt.Println("Sending data")
			_, _ = w.Write(out)

			time.Sleep(time.Duration(timeBetweenBatches) * time.Millisecond)
		}
	}()

	resp, err := client.Post(address, "", r)

	if err != nil {
		return err
	}

	bIn := make([]byte, 1024)

	received := 0
	n, err := resp.Body.Read(bIn)

	for ; err == nil; n, err = resp.Body.Read(bIn) {

		if n == 0 {
			continue
		}
		fmt.Println("Receiving data")
		received += n
	}

	fmt.Fprintf(os.Stdout, "%d bytes received from %s \n", received, address)

	if err != io.EOF {
		return err
	}

	return nil

}