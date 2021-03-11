package main

import (
	"devportal/api/apispec"
	"devportal/api/product"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIPGTenant(t *testing.T) {

	req, err := http.NewRequest("GET", "/v1/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(product.GetProductsInfo)
	handler.ServeHTTP(rr, req)

	//Check if API returned 200 OK
	checkResponseCode(t, http.StatusOK, rr.Code)

}

func TestIPGAPISpecs(t *testing.T) {

	req, err := http.NewRequest("POST", "/v1/apis", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(apispec.GetAPISpecs)
	handler.ServeHTTP(rr, req)

	//Check if API returned 200 OK
	checkResponseCode(t, http.StatusOK, rr.Code)

}

//Validates response code
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	} else {
		t.Log(expected, actual)
	}
}
