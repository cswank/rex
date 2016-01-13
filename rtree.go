package rtree

import "net/http"

type handler func(w http.ResponseWriter, req *http.Request)

func New() *Tree {
	return &Tree{}
}

type Tree struct {
	f handler
}

func (r *Tree) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.f(w, req)
}

func (r *Tree) Get(path string, f handler) {
	r.f = f
}

func (r *Tree) Post(path string, f handler) {
	r.f = f
}

func (r *Tree) Put(path string, f handler) {
	r.f = f
}

func (r *Tree) Delete(path string, f handler) {
	r.f = f
}
