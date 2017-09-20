package main

import (
	"net"
	"net/http"
	"time"
	"io/ioutil"
	"io"
)

func test(client *http.Client) {
	println()
	res, err := client.Get("http://www.google.pl")
	if err != nil {
		return
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)
	for header, values := range res.Header {
		for _, value := range values {
			println(header, ":", value)
		}
	}
}

func main() {
	tr := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			println("Connecting to", network, addr)
			return (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial(network, addr)
		},
	}
	client := &http.Client{
		Transport: tr,
	}
	test(client)
	test(client)
}
