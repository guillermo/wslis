package rhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/guillermo/wslis"
)

type ClientHandlerFunc func(header http.Header, client *http.Client)

func (ch ClientHandlerFunc) HandleEndpoint(d *wslis.Dialer) {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: nil,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return d.Dial()
			},
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	ch(d.Header, client)

}
