package main

import (
	"fmt"
	"io"
	"net/http"
)

const (
	proxyPort   = 8000
	servicePort = 80
)

type Proxy struct{}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// forward the request to the new aad endpoint
	resp, err := p.sendRequest(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	p.writeResponse(w, resp)
}

func main() {
	http.ListenAndServe(fmt.Sprintf(":%d", proxyPort), &Proxy{})
}

func (p *Proxy) sendRequest(req *http.Request) (*http.Response, error) {
	proxyURL := fmt.Sprintf("http://127.0.0.1:%d%s", servicePort, req.RequestURI)
	// Create an HTTP client and a proxy request based on the original request.
	httpClient := http.Client{}
	proxyReq, err := http.NewRequest(req.Method, proxyURL, req.Body)

	resp, err := httpClient.Do(proxyReq)
	return resp, err
}

func (p *Proxy) writeResponse(w http.ResponseWriter, res *http.Response) {
	for name, values := range res.Header {
		w.Header()[name] = values
	}
	// Set a special header to notify that the proxy actually serviced the request.
	w.Header().Set("Server", "my-proxy")
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
	res.Body.Close()
}
