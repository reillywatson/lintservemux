package a

import (
	"net"
	"net/http"
)

type handler struct{}

func (h handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("cool"))
}

func Foo() {
	http.Handle("/a", &handler{}) // want "http.Handle uses DefaultServeMux"

	http.HandleFunc("/b", func(rw http.ResponseWriter, req *http.Request) { // want "http.HandleFunc uses DefaultServeMux"
		rw.Write([]byte("wooo"))
	})

	mux := http.NewServeMux()

	http.ListenAndServe(":6060", nil) // want "http.ListenAndServe should pass an http.Handler"
	http.ListenAndServe(":6061", mux)

	listener, _ := net.ListenTCP("ip", &net.TCPAddr{Port: 6060})
	http.ServeTLS(listener, nil, "", "") // want "http.ServeTLS should pass an http.Handler"
	http.ServeTLS(listener, mux, "", "")

	http.Serve(listener, nil) // want "http.Serve should pass an http.Handler"
	http.Serve(listener, mux)
}
