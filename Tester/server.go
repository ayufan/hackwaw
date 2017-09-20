package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func forwardTweets(rw http.ResponseWriter, req *http.Request) {
	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			url, _ := url.Parse(*apiServer)
			req.Host = url.Host
			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
			req.URL.Path = "/tweets"
			log.Println("PROXY:", req.URL)
		},
	}
	proxy.ServeHTTP(rw, req)
}

func forwardPushes(rw http.ResponseWriter, req *http.Request) {
	emptyText(rw, req)
	//	proxy := httputil.ReverseProxy{
	//		Director: func(req *http.Request) {
	//			url, _ := url.Parse(*apiServer)
	//			req.Host = url.Host
	//			req.URL.Scheme = url.Scheme
	//			req.URL.Host = url.Host
	//			req.URL.Path = "/push"
	//			log.Println("PROXY:", req.URL)
	//		},
	//	}
	//	proxy.ServeHTTP(rw, req)
}

func validateTweetRequest(handler http.HandlerFunc, tweetReceiver func(tweet Tweet)) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			badRequest(rw, req, "not GET")
			return
		}

		_, err := time.Parse(time.RFC3339, req.FormValue("from"))
		if err != nil {
			badRequest(rw, req, "invalid `from`", err.Error())
			return
		}

		_, err = time.Parse(time.RFC3339, req.FormValue("to"))
		if err != nil {
			badRequest(rw, req, "invalid `to`", err.Error())
			return
		}

		handler(rw, req)

		brw, ok := rw.(bytesWriter)

		if ok && brw.Status() == http.StatusOK && tweetReceiver != nil {
			var tweets []Tweet
			json.Unmarshal(brw.Bytes(), &tweets)
			for _, tweet := range tweets {
				tweetReceiver(tweet)
			}
		}
	}
}

func validatePushRequest(handler http.HandlerFunc, pushReceiver func(push Slack) error) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			badRequest(rw, req, "not POST")
			return
		}

		mt, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
		if err != nil {
			badRequest(rw, req, "error parsing Content-Type:", err.Error())
			return
		}
		if mt != "application/json" {
			badRequest(rw, req, "not application/json")
			return
		}

		data, err := ioutil.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			badRequest(rw, req, "read body:", err.Error())
			return
		}

		var push Slack
		err = json.Unmarshal(data, &push)
		if err != nil {
			badRequest(rw, req, "decode body:", err.Error())
			return
		}

		if pushReceiver != nil {
			err = pushReceiver(push)
			if err != nil {
				serverError2(rw, req, "push receiver:", err.Error())
				return
			}
		}

		req.Body = ioutil.NopCloser(bytes.NewReader(data))
		handler(rw, req)
	}
}

func slowDown(duration time.Duration, handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		time.Sleep(duration)
		handler(rw, req)
	}
}

func badRequest(rw http.ResponseWriter, req *http.Request, message ...string) {
	rw.Header().Set("X-HTTP-Status", "400 Bad request "+strings.Join(message, ","))
	http.Error(rw, "400 Bad request "+strings.Join(message, ","), 400)
}

func serverError2(rw http.ResponseWriter, req *http.Request, message ...string) {
	rw.Header().Set("X-HTTP-Status", "500 Server error "+strings.Join(message, ","))
	http.Error(rw, "500 Server error "+strings.Join(message, ","), 500)
}

func serverError(rw http.ResponseWriter, req *http.Request) {
	serverError2(rw, req, "generic error")
}

func invalidJson(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		badRequest(rw, req, "not GET")
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	fmt.Fprint(rw, "''this''can't''be''valid''json!'")
}

func jsonNoContentType(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		badRequest(rw, req, "not GET")
		return
	}

	fmt.Fprint(rw, "[]")
}

func emptyJson(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		badRequest(rw, req, "not GET")
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	fmt.Fprint(rw, "[]")
}

func emptyText(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/plain")
	fmt.Fprint(rw, "ok")
}

func validateMethod(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			badRequest(rw, req, "not "+method)
			return
		}
		handler(rw, req)
	}
}

func jsonResponse(response interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		data, err := json.Marshal(response)
		if err != nil {
			serverError(rw, req)
			return
		}
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(data)
	}
}

func startServer(tweets, push http.HandlerFunc) *httptest.Server {
	mux := http.NewServeMux()
	mux.Handle("/tweets", tweets)
	mux.Handle("/push", push)

	l, err := net.Listen("tcp", *serverListen)
	if err != nil {
		panic(fmt.Sprintf("startServer: failed to listen on %v: %v", *serverListen, err))
	}

	loggingHandler := func(rw http.ResponseWriter, req *http.Request) {
		lrw := newLoggingResponseWriter(rw)
		defer lrw.Log(req)
		mux.ServeHTTP(lrw, req)
	}

	server := &httptest.Server{
		Listener: l,
		Config:   &http.Server{Handler: http.HandlerFunc(loggingHandler)},
	}
	server.Start()
	return server
}
