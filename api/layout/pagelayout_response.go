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
