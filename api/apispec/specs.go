package apispec

import (
	"bytes"
	"devportal/api/product"
	"net/http"
)

//GetAPISpecs : Returns API Specifications in YAML format
func GetAPISpecs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/yaml; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//Get API Spec/YAML File By Path from github repo
	response, err := http.Get(product.DevPortalConfig.GitHub.GitHubContentFullPath + product.DevPortalConfig.ContentPath.ApiSpecYamlFile)

	if err != nil {
		w.Write([]byte("Swagger file not found in github repo with name :: " + product.DevPortalConfig.ContentPath.ApiSpecYamlFile))
		w.WriteHeader(http.StatusBadRequest)
	} else {
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		responseByte := buf.Bytes()
		defer response.Body.Close()

		w.Write(responseByte)
		w.WriteHeader(http.StatusOK)
	}

}
