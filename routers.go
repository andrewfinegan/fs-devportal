// ################################################################################
// Copyright © 2021-2022 Fiserv, Inc. or its affiliates. 
// Fiserv is a trademark of Fiserv, Inc., 
// registered or used in the United States and foreign countries, 
// and may or may not be registered in your country.  
// All trademarks, service marks, 
// and trade names referenced in this 
// material are the property of their 
// respective owners. This work, including its contents 
// and programming, is confidential and its use 
// is strictly limited. This work is furnished only 
// for use by duly authorized licensees of Fiserv, Inc. 
// or its affiliates, and their designated agents 
// or employees responsible for installation or 
// operation of the products. Any other use, 
// duplication, or dissemination without the 
// prior written consent of Fiserv, Inc. 
// or its affiliates is strictly prohibited. 
// Except as specified by the agreement under 
// which the materials are furnished, Fiserv, Inc. 
// and its affiliates do not accept any liabilities 
// with respect to the information contained herein 
// and are not responsible for any direct, indirect, 
// special, consequential or exemplary damages 
// resulting from the use of this information. 
// No warranties, either express or implied, 
// are granted or extended by this work or 
// the delivery of this work
// ################################################################################

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
