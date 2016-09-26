package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var (
	serverOnce sync.Once
	server     *Server
)

type Server struct {
	sync.RWMutex
	counter int64
}

func (s *Server) Inc() {
	s.Lock()
	defer s.Unlock()
	s.counter++
}

func (s *Server) Counter() int64 {
	s.RLock()
	defer s.RUnlock()
	c := s.counter
	return c
}

func router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods("GET")
	r.HandleFunc("/counter", func(w http.ResponseWriter, r *http.Request) {
		server.Inc()
		w.Write([]byte(fmt.Sprintf("server.Counter() => %d", server.Counter())))
	}).Methods("GET")

	return r
}

func main() {
	http.Handle("/", router())
	http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil)
}

func init() {
	serverOnce.Do(func() {
		server = new(Server)
	})
}
