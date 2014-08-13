package main

import (
	"flag"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type Handler interface {
	Router(*mux.Router)
}

type Response struct {
	Success  bool
	Error    error
	Response interface{}
}

var router = mux.NewRouter()
var Listen = flag.String("listen", ":8080", "Listen addresss")

func Register(name string, handler Handler) {
	handler.Router(router.PathPrefix("/" + name + "/").Name(name).Subrouter())
}

func main() {
	flag.Parse()
	handler := handlers.CombinedLoggingHandler(os.Stderr, router)
	http.ListenAndServe(*Listen, handler)
}
