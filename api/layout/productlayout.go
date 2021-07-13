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
