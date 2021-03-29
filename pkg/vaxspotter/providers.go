package vaxspotter

import (
	"time"
)

type ProviderBrand struct {
	ID                      int                 `json:"id"`
	Key                     string              `json:"key"`
	URL                     string              `json:"url"`
	Name                    string              `json:"name"`
	Status                  ProviderBrandStatus `json:"status"`
	ProviderID              string              `json:"provider_id"`
	LocationCount           int                 `json:"location_count"`
	AppointmentsLastFetched time.Time           `json:"appointments_last_fetched"`
}

type ProviderBrandStatus string

const (
	ProviderBrandStatusActive   ProviderBrandStatus = "active"
	ProviderBrandStatusInactive ProviderBrandStatus = "inactive"
	ProviderBrandStatusUnknown  ProviderBrandStatus = "unknown"
)
