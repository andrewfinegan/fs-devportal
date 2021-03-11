package apispec

//SpecsFilterRequest : Filter Request
type SpecsFilterRequest struct {
	Categories []struct {
		Name          string `json:"name"`
		Subcategories []struct {
			Name string `json:"name"`
		} `json:"subcategories"`
	} `json:"categories"`
}
