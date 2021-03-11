package layout

//PageLayoutResponse : Product page layout
type PageLayoutResponse []struct {
	Hero []struct {
		Image      string `json:"image,omitempty"`
		Descrption string `json:"descrption,omitempty"`
	} `json:"hero,omitempty"`
	JourneyCards []struct {
		JourneyCard []struct {
			Title       string `json:"title,omitempty"`
			Description string `json:"description,omitempty"`
			Link        []struct {
				Type string `json:"type,omitempty"`
				URI  string `json:"uri,omitempty"`
			} `json:"link,omitempty"`
			ContextualNavigationCards []struct {
				ContextualNavigationCard []struct {
					Title       string `json:"title,omitempty"`
					Description string `json:"description,omitempty"`
					Link        []struct {
						Type string `json:"type,omitempty"`
						URI  string `json:"uri,omitempty"`
					} `json:"link,omitempty"`
				} `json:"contextualNavigationCard"`
			} `json:"contextualNavigationCards,omitempty"`
		} `json:"journeyCard"`
	} `json:"journey Cards,omitempty"`
}
