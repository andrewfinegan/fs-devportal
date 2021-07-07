package layout

import (
	"bytes"
	"devportal/api/product"
	"devportal/config"
	"net/http"

	"github.com/ghodss/yaml"
)

//GetProductLayout : Get Layout YAML from Github for ProductPage and convert to JSON and Respond.
func GetProductLayout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Get the product layout file content
	//Use product layout file path
	response, err := http.Get(product.DevPortalConfig.GitHub.GitHubContentFullPath + product.DevPortalConfig.ContentPath.ProductLayoutFile)

	if err != nil {
		config.AddLogFields(config.Logger).Error(err)
		w.Write([]byte("Product Layout not found"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	respByte := buf.Bytes()
	defer response.Body.Close()

	//Convert YAML to JSON
	jsonResponse, err := yaml.YAMLToJSON(respByte)

	if err != nil {
		config.AddLogFields(config.Logger).Println("Yaml to json conversion failed")
		_, err = w.Write([]byte("Yaml to json conversion failed"))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		_, err = w.Write(jsonResponse)
		w.WriteHeader(http.StatusOK)
	}

}
