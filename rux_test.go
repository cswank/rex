package rux_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cswank/rux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gadgets", func() {
	var (
		r   *rux.Tree
		w   *httptest.ResponseRecorder
		req *http.Request
	)

	BeforeEach(func() {
		root := func(ww http.ResponseWriter, rr *http.Request) {
			ww.Write([]byte("root"))
		}

		pals := func(ww http.ResponseWriter, rr *http.Request) {
			ww.Write([]byte("pals"))
		}

		r = rux.New()
		r.Get("/", root)
		r.Get("/pals", pals)

		w = httptest.NewRecorder()
	})

	AfterEach(func() {
	})

	It("gets root", func() {
		var err error
		req, err = http.NewRequest("GET", "/", nil)
		Expect(err).To(BeNil())
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.String()).To(Equal("root"))
	})

	It("gets the first collection", func() {
		var err error
		req, err = http.NewRequest("GET", "/pals", nil)
		Expect(err).To(BeNil())
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.String()).To(Equal("pals"))
	})
})
