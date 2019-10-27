package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"net/textproto"
	"time"
)

// Request defines a requests with performance metrics
type Request struct {
	*http.Request
}

// NewRequest creates a new request with performance metrics under a specific context, that allows
// to cancel or timeout the request.
func NewRequest(method, url string, body io.Reader) (*Request, error) {
	req, err := http.NewRequest(method, url, body)

	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			fmt.Printf("GetConn - hostPort: %s\n\tTime: %s\n", hostPort, time.Now().String())
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("GotCon - connInfo: %v\n\tIdletime: %s\n\tTime: %s\n",
				connInfo, connInfo.IdleTime, time.Now().String())
		},
		PutIdleConn: func(err error) {
			fmt.Printf("PutIdleConn - err: %v\n\tTime: %s\n", err, time.Now().String())
		},
		GotFirstResponseByte: func() {
			fmt.Printf("GotFirstResponseByte - \n\tTime: %s\n", time.Now().String())
		},
		Got100Continue: func() {
			fmt.Printf("Got100Continue - \n\tTime: %s\n", time.Now().String())
		},
		Got1xxResponse: func(code int, header textproto.MIMEHeader) error {
			fmt.Printf("Got1xxResponse - code: %d\tMime: %v\n\tTime: %s\n", code, header, time.Now().String())
			return nil
		},
		DNSStart: func(dsi httptrace.DNSStartInfo) {
			fmt.Printf("DNSStart - DNSStartInfo host: %s\n\tTime: %s\n", dsi.Host, time.Now().String())
		},
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			fmt.Printf("DNSDone - DNSDoneInfo: %v\n\tTime: %s\n", ddi, time.Now().String())
		},
		ConnectStart: func(network, addr string) {
			fmt.Printf("ConnectStart - network: %s\taddr: %s\n\tTime: %s\n", network, addr, time.Now().String())
		},
		ConnectDone: func(network, addr string, err error) {
			fmt.Printf("ConnectDone - network: %s\taddr: %s\terror: %s\n\tTime: %s\n", network, addr, err, time.Now().String())
		},
		TLSHandshakeStart: func() {
			fmt.Printf("TLSHandshakeStart - \n\tTime: %s\n", time.Now().String())
		},
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			fmt.Printf("TLSHandshakeDone - ConnectionState completed: %t\terror: %s\n\tTime: %s\n", cs.HandshakeComplete, err, time.Now().String())
		},
		WroteHeaderField: func(key string, value []string) {
			fmt.Printf("WroteHeaderField - key: %s\tvalue: %v\n\tTime: %s\n", key, value, time.Now().String())
		},
		WroteHeaders: func() {
			fmt.Printf("WroteHeaders - \n\tTime: %s\n", time.Now().String())
		},
		Wait100Continue: func() {
			fmt.Printf("Wait100Continue - \n\tTime: %s\n", time.Now().String())
		},
		WroteRequest: func(wri httptrace.WroteRequestInfo) {
			fmt.Printf("WroteRequest - WroteRequestInfo: %v\n\tTime: %s\n", wri, time.Now().String())
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	return &Request{req}, err
}

func main() {
	req, err := NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	client := http.DefaultClient
	fmt.Printf("Start connection - \n\tTime: %s\n", time.Now().String())
	resp, err := client.Do(req.Request)
	fmt.Printf("End connection - \n\tTime: %s\n", time.Now().String())
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(bodyBytes[:20]))
	fmt.Println("\n-----------------")

}
