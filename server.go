package wslis

import (
	"io"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/yamux"
)

// Server will waits for connections to arrives and then calls EndpointHandler
// It can be use directly with the ListenAndServe functions or indirectly as an http.Handler
type Server struct {
	// Addr to listen (only needed if started with ListenAndServe)
	Addr string

	// Handler called when a new connection is called.
	Handler EndpointHandler
	s       *http.Server
}

// EndpointHandler responds to a new connection from a client
type EndpointHandler interface {
	HandleEndpoint(*Dialer)
}

// The EndpointHandlerFunc type is an adapter to allow the use of ordinary functions as Endpoint Handlers.
type EndpointHandlerFunc func(*Dialer)

// HandleEndpoint calls f(d)
func (f EndpointHandlerFunc) HandleEndpoint(d *Dialer) {
	f(d)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ListenAndServe is a convenient function that starts a server in the given
// address and send all new clients to handler
func ListenAndServe(addr string, handler EndpointHandler) error {
	s := &Server{
		Addr:    addr,
		Handler: handler,
	}
	return s.ListenAndServe()
}

// ListenAndServeTLS is the TLS version of ListenAndServe
func ListenAndServeTLS(addr string, certFile, keyFile string, handler EndpointHandler) error {
	s := &Server{
		Addr:    addr,
		Handler: handler,
	}
	return s.ListenAndServeTLS(certFile, keyFile)
}

// ListenAndServe will create a basic http server listening in the given address.
// If you need more control over the server, you can use the hole server as an http.Handler
func (s *Server) ListenAndServe() error {
	if s.s == nil {
		s.s = &http.Server{
			Addr:    s.Addr,
			Handler: s,
		}
	}
	return s.s.ListenAndServe()
}

// ListenAndServeTLS is the TLS version of ListenAndServe
func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	if s.s == nil {
		s.s = &http.Server{
			Addr:    s.Addr,
			Handler: s,
		}
	}
	return s.s.ListenAndServeTLS(certFile, keyFile)

}

// ServeHTTP establish the connection with the client and calls the EndpointHandler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	conn := newconn(wsconn)
	d, err := newDialer(conn)
	if err != nil {
		panic(err)
	}
	d.Header = r.Header

	s.Handler.HandleEndpoint(d)

}

// Dialer represent an active connection with a client, ready to be open
type Dialer struct {
	session *yamux.Session
	Header  http.Header
}

func newDialer(conn io.ReadWriteCloser) (*Dialer, error) {
	session, err := yamux.Client(conn, nil)
	return &Dialer{session, nil}, err
}

// Dial opens a new connection to the client
func (d *Dialer) Dial() (net.Conn, error) {
	return d.session.Open()
}
