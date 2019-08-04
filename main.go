package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type store struct {
	data map[string]string
	m    sync.RWMutex
}

var (
	addr = flag.String("addr", ":8081", "http service address")
	s    = store{
		data: map[string]string{},
		m:    sync.RWMutex{},
	}
)

func main() {
	flag.Parse()

	r := httprouter.New()

	r.GET("/entry/:key", show)
	r.GET("/list", show)
	r.PUT("/entry/:key/:value", update)

	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	k := p.ByName("key")
	if k == "" {
		s.m.RLock()
		fmt.Fprintf(w, "Read list: %v", s.data)
		s.m.RUnlock()
		return
	}
	s.m.RLock()
	fmt.Fprintf(w, "Read entry: s.data[%s] = %s", k, s.data[k])
	s.m.RUnlock()
}

func update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	k := p.ByName("key")
	v := p.ByName("value")
	s.m.Lock()
	s.data[k] = v
	s.m.Unlock()
	fmt.Fprintf(w, "Updated: s.data[%s] = %s", k, v)
}
