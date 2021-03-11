package product

//TenantProviderResponse : Tenant Provider API Response
type TenantProviderResponse struct {
	Title            string `json:"title"`
	Name             string `json:"name"`
	Brand            string `json:"brand"`
	BrandLogoURL     string `json:"brandLogoURL"`
	Solution         string `json:"solution"`
	BrandDescription string `json:"brandDescription"`
	Categories       []struct {
		Name          string `json:"name"`
		Value         string `json:"value"`
		Subcategories []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"subcategories"`
	} `json:"categories"`
	Product struct {
		Featured         bool   `json:"featured"`
		LogoURL          string `json:"logoURL"`
		Description      string `json:"description"`
		APISpecification string `json:"apiSpecification"`
		Layout           string `json:"layout"`
		Documentation    string `json:"documentation"`
	} `json:"product"`
}
