package rux

import (
	"net/http"
	"strings"
)

var tree *Tree

type handler func(w http.ResponseWriter, req *http.Request)

func New() *Tree {
	tree = &Tree{
		methods: map[string]*node{
			"GET":    newNode(),
			"POST":   newNode(),
			"PUT":    newNode(),
			"DELETE": newNode(),
			"PATCH":  newNode(),
		},
	}
	return tree
}

type Tree struct {
	methods map[string]*node
}

func (t *Tree) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n, ok := t.methods[r.Method]
	if !ok {
		notFound(w)
	} else {
		n.handle(strings.Split(strings.Trim(r.URL.Path, "/"), "/"), w, r)
	}
}

func Vars(r *http.Request) map[string]string {
	pth := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	m := map[string]string{}
	tree.methods[r.Method].vars(pth, m)
	return m
}

func (t *Tree) Get(pth string, h handler) *node {
	return t.method("GET", pth, h)
}

func (t *Tree) Post(pth string, h handler) *node {
	return t.method("POST", pth, h)
}

func (t *Tree) Put(pth string, h handler) *node {
	return t.method("PUT", pth, h)
}

func (t *Tree) Delete(pth string, h handler) *node {
	return t.method("DELETE", pth, h)
}

func (t *Tree) Path(pth string, h handler) *node {
	return t.method("PATCH", pth, h)
}

func (t *Tree) method(name, pth string, h handler) *node {
	return t.methods[name].add(strings.Split(strings.Trim(pth, "/"), "/"), h)
}

type node struct {
	id       string
	key      string
	handler  handler
	resource bool
	children map[string]*node
}

func newNode() *node {
	return &node{
		children: map[string]*node{},
	}
}

func (n *node) vars(pth []string, m map[string]string) {
	if len(pth) == 0 && !n.resource {
		return
	}

	var c *node
	var ok bool
	l := len(pth)
	if n.resource && l != 0 {
		m[n.key] = pth[0]
		c, ok = n.children[n.id]
	} else if l != 0 {
		c, ok = n.children[pth[0]]
	}

	if ok && l != 0 {
		c.vars(pth[1:], m)
	}
}

func (n *node) add(pth []string, h handler) *node {
	if len(pth) == 0 {
		n.handler = h
		return n
	}
	n.resource, n.id, n.key = n.isResource(pth[0])
	var x string
	if n.resource {
		x = n.id
	} else {
		x = pth[0]
	}
	c, ok := n.children[x]
	if !ok {
		c = newNode()
		n.children[pth[0]] = c
	}
	return c.add(pth[1:], h)
}

func (n *node) isResource(s string) (bool, string, string) {
	var id string
	r := strings.Index(s, "{") == 0
	if r {
		id = s
	}
	return r, id, strings.Trim(id, "{}")
}

func (n *node) handle(pth []string, w http.ResponseWriter, r *http.Request) {
	if len(pth) == 0 {
		n.handler(w, r)
		return
	}
	var x string
	if n.resource {
		x = n.id
	} else {
		x = pth[0]
	}
	if c, ok := n.children[x]; ok {
		c.handle(pth[1:], w, r)
	} else {
		notFound(w)
	}
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}
