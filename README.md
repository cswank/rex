Rex
===

The RESTful mux package.  If you have a RESTful
set or routes like:

    /
    /friends
    /friends/{friend}
    /friends/{friend}/hobbies/{hobby}
    /friends/{friend}/songs/{song}

Then you should be able to use Rex.  It doesn't match
routes based on regex and only provides enough functionality
to work with a RESTful api.  Rex can be used to serve the above
routes like this:

    package main

    import (
    	"log"
    	"net/http"

    	"github.com/cswank/rex"
    )

    func main() {
    	r := rex.New("example")
    	r.Get("/things", getThings)
        r.Post("/things", addThing)
    	r.Get("/things/{thing}", getThing)
        r.Delete("/things/{thing}", deleteThing)

    	http.Handle("/", r)
    	log.Println(http.ListenAndServe(":8888", r))
    }

    func getThings(w http.ResponseWriter, r *http.Request) {
    	w.Write([]byte("things"))
    }

    func addThing(w http.ResponseWriter, r *http.Request) {
        
    }

    func getThing(w http.ResponseWriter, r *http.Request) {
    	w.Write([]byte("thing"))
    }

    func deleteThing(w http.ResponseWriter, r *http.Request) {
    
    }

You can also serve files:

    package main

    import (
    	"log"
    	"net/http"

    	"github.com/cswank/rex"
    )

    func main() {
    	r := rex.New("example")
    	r.Get("/things", getThings)
        r.ServeFiles(http.FileServer(http.Dir("/www")))
    	http.Handle("/", r)
    	log.Println(http.ListenAndServe(":8888", r))
    }

    func getThings(w http.ResponseWriter, r *http.Request) {
    	w.Write([]byte("things"))
    }

The supported HTTP methods are Get, Post, Put, Delete, and Patch.

When you create a new Router you must pass in a name.  When getting
the vars from a request the name of the router must be passed in.


    func getThings(w http.ResponseWriter, r *http.Request) {
        vars := rex.Vars("example")
    	w.Write([]byte("things"))
    }

