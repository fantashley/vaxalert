package vaxspotter

import "time"

type State struct {
	Code                    string          `json:"code"`
	Name                    string          `json:"name"`
	StoreCount              string          `json:"store_count"`
	ProviderBrandCount      string          `json:"provider_brand_count"`
	AppointmentsLastFetched time.Time       `json:"appointments_last_fetched"`
	ProviderBrands          []ProviderBrand `json:"provider_brands"`
}
