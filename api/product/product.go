package product

import (
	"bytes"
	"devportal/config"
	"net/http"
)

var DevPortalConfig config.Config

// GetProductsInfo : Returns Product Information Response containing details of tenant product
func GetProductsInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the product layout file content
	response, err := http.Get(DevPortalConfig.GitHub.GitHubContentFullPath + DevPortalConfig.ContentPath.TenantProviderApiFile) //Fetch tenant Provider API Json
	if err != nil {
		config.AddLogFields(config.Logger).Error(err)
		w.Write([]byte("Tenant Provider API not found "))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Body)
	respByte := buf.Bytes()
	defer response.Body.Close()

	_, err = w.Write([]byte(respByte))
	w.WriteHeader(http.StatusOK)
}
