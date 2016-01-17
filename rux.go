package rux

import (
	"net/http"
	"strings"
)

type handler func(w http.ResponseWriter, req *http.Request)

func New() *Tree {
	return &Tree{
		root: newNode(),
	}
}

type Tree struct {
	root *node
}

func (t *Tree) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.root.handle(strings.Split(strings.Trim(r.URL.Path, "/"), "/"), w, r)
}

func (t *Tree) Get(pth string, h handler) *node {
	return t.root.add(strings.Split(strings.Trim(pth, "/"), "/"), h)
}

func (t *Tree) Post(pth string, h handler) {
}

func (t *Tree) Put(pth string, h handler) {
}

func (t *Tree) Delete(pth string, h handler) {
}

type node struct {
	handler  handler
	resource bool
	children map[string]*node
}

func newNode() *node {
	return &node{
		children: map[string]*node{},
	}
}

func (n *node) add(pth []string, h handler) *node {
	if len(pth) == 0 {
		n.handler = h
		return n
	}
	c, ok := n.children[pth[0]]
	if !ok {
		c = newNode()
		n.children[pth[0]] = c
	}
	return c.add(pth[1:], h)
}

func (n *node) isResource(s string) bool {
	return strings.Index(s, "{") == 0
}

func (n *node) handle(pth []string, w http.ResponseWriter, r *http.Request) {
	if len(pth) == 0 || (len(pth) == 1 && n.isResource(pth[0])) {
		n.handler(w, r)
		return
	}
	if n.resource {
		n.handle(pth[1:], w, r)
	} else if c, ok := n.children[pth[0]]; ok {
		c.handle(pth[1:], w, r)
	} else {
		notFound(w)
	}
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}
