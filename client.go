package wslis

import (
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/yamux"
)

// A Client calls a server through websocket and waits for connections.
type Client struct {
	// URL of the server to connect.
	// It should be ws:// or wss:// protocol
	URL string

	// Headers to send to the server while opening a connection
	Header http.Header

	// HandshakeTimeout specifies the duration for the handshake to complete.
	HandshakeTimeout time.Duration
}

// DialAndListen does a websocket
func (c *Client) DialAndListen() (net.Listener, error) {

	dialer := &websocket.Dialer{
		HandshakeTimeout: c.HandshakeTimeout,
	}

	wsconn, _, err := dialer.Dial(c.URL, c.Header)
	if err != nil {
		return nil, err
	}

	conn := newconn(wsconn)

	return yamux.Server(conn, nil)
}

// DialAndListen creates a client and connects to the given url
func DialAndListen(url string, header http.Header) (net.Listener, error) {
	c := &Client{URL: url, Header: header}
	return c.DialAndListen()

}
