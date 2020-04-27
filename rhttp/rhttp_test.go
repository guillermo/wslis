package rhttp

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/guillermo/wslis"
)

func ExampleTLS() {
	result := make(chan (string))

	// Server
	go wslis.ListenAndServe("localhost:9094", ClientHandlerFunc(func(header http.Header, client *http.Client) {

		client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		buf := make([]byte, 1024)
		id := header.Get("DeviceId")
		if id != "123" {
			panic("We don't know you")
		}

		res, err := client.Get("https://asdf/patatas")
		if err != nil {
			panic(err)
		}
		n, err := res.Body.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n < 2 {
			panic(n)
		}

		result <- string(buf[:n])

	}))

	// Client
	go func() {
		time.Sleep(time.Millisecond * 10) // Lets give time to the server to start
		lis, err := wslis.DialAndListen("ws://localhost:9094", http.Header{"DeviceId": []string{"123"}})
		if err != nil {
			panic(err)
		}

		err = http.ServeTLS(lis, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			_, err := w.Write([]byte(r.URL.Path[1:]))
			if err != nil {
				panic(err)
			}

		}), "ssl.crt", "ssl.key")
		if err != nil {
			panic(err)
		}

	}()

	fmt.Println(<-result)
	// Output: patatas
}

func Example() {

	result := make(chan (string))

	// Server
	go wslis.ListenAndServe("localhost:9090", ClientHandlerFunc(func(header http.Header, client *http.Client) {

		buf := make([]byte, 1024)
		id := header.Get("DeviceId")
		if id != "123" {
			panic("We don't know you")
		}

		res, err := client.Get("http://asdf/hola")
		if err != nil {
			panic(err)
		}
		n, err := res.Body.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n < 2 {
			panic(n)
		}

		result <- string(buf[:n])

	}))

	// Client
	go func() {
		time.Sleep(time.Millisecond) // Lets give time to the server to start
		lis, err := wslis.DialAndListen("ws://localhost:9090", http.Header{"DeviceId": []string{"123"}})
		if err != nil {
			panic(err)
		}

		err = http.Serve(lis, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			_, err := w.Write([]byte(r.URL.Path[1:]))
			if err != nil {
				panic(err)
			}

		}))
		if err != nil {
			panic(err)
		}

	}()

	fmt.Println(<-result)
	// Output: hola
}
