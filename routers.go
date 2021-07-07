package main

import (
	"devportal/api/apispec"
	"devportal/api/doc"
	"devportal/api/layout"
	"devportal/api/product"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//Route struct
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes array of Route
type Routes []Route

//NewRouter : router method
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

//Index : Default route
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

//HealthCheck probe
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sample Tenant Server is up and running !!!"))
}

/*//Function to override default handler and validate for basic authentication. This will return 401 for incorrect auth
func auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if !validate(user, pass) {
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}

func validate(username, password string) bool {
	fmt.Printf("Username: %s, password: %s", username, password)

	if username == "test" && password == "test" {
		return true
	}

	return false
}*/

//Define Routes for tenant
var routes = Routes{
	Route{
		"Index",
		"GET",
		"/v2/",
		Index,
	},
	Route{
		"Swagger",
		"POST",
		"/v1/apis",
		apispec.GetAPISpecs,
	},

	Route{
		"Documentation",
		"GET",
		"/v1/docs",
		doc.GetDocumentServiceSelector,
	},

	Route{
		"Tenant Provider API",
		"GET",
		"/v1/products",
		product.GetProductsInfo,
	},

	Route{
		"Health Check API",
		"GET",
		"/healthcheck",
		HealthCheck,
	},

	Route{
		"Product Page Layout API",
		"GET",
		"/v1/layouts",
		layout.GetProductLayout,
	},
}
