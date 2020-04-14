package wslis

import (
	"fmt"
	"net/http"
	"time"
)

func Example() {

	result := make(chan (string))

	// Server
	go ListenAndServe("localhost:9090", EndpointHandlerFunc(func(d *Dialer) {
		buf := make([]byte, 1024)

		id := d.Header.Get("DeviceId")
		if id != "123" {
			panic("We don't know you")
		}

		// We open a new connection to the client as soon as we get the connection
		conn, err := d.Dial()
		if err != nil {
			panic(err)
		}

		// We say hi
		conn.Write([]byte("hola"))

		// And we expect a hi back
		n, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		result <- string(buf[:n])
	}))

	// Client
	go func() {
		time.Sleep(time.Millisecond) // Lets give time to the server to start
		lis, err := DialAndListen("ws://localhost:9090", http.Header{"DeviceId": []string{"123"}})
		if err != nil {
			panic(err)
		}

		// Let's have an echo server
		conn, err := lis.Accept()
		if err != nil {
			panic(err)
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}
		conn.Write(buf[:n])
		conn.Close()
	}()

	fmt.Println(<-result)
	// Output: hola
}

func ExampleListenAndServeTLS() {

	result := make(chan (string))

	// Server
	go ListenAndServeTLS("localhost:9443", "ssl.crt", "ssl.key", EndpointHandlerFunc(func(d *Dialer) {
		buf := make([]byte, 1024)

		id := d.Header.Get("DeviceId")
		if id != "123" {
			panic("We don't know you")
		}

		// We open a new connection to the client as soon as we get the connection
		conn, err := d.Dial()
		if err != nil {
			panic(err)
		}

		// We say hi
		conn.Write([]byte("hola"))

		// And we expect a hi back
		n, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		result <- string(buf[:n])
	}))

	// Client
	go func() {
		time.Sleep(time.Second * 2) // Wait for the server to start
		lis, err := DialAndListen("wss://localhost:9443", http.Header{"DeviceId": []string{"123"}})
		if err != nil {
			panic(err)
		}

		// Let's have an echo server
		conn, err := lis.Accept()
		if err != nil {
			panic(err)
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}
		conn.Write(buf[:n])
		conn.Close()
	}()

	fmt.Println(<-result)
	// Output: hola
}
