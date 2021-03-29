package vaxspotter

import "time"

type Locations struct {
	Type      string     `json:"type"`
	Locations []Location `json:"features"`
	Metadata  struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		StoreCount  int    `json:"store_count"`
		BoundingBox struct {
			Type        string        `json:"type"`
			Coordinates [][][]float64 `json:"coordinates"`
		} `json:"bounding_box"`
		ProviderBrands          []ProviderBrand `json:"provider_brands"`
		ProviderBrandCount      int             `json:"provider_brand_count"`
		AppointmentsLastFetched time.Time       `json:"appointments_last_fetched"`
	} `json:"metadata"`
}

type Location struct {
	Type     string `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		ID                               int                             `json:"id"`
		URL                              string                          `json:"url"`
		City                             string                          `json:"city"`
		Name                             string                          `json:"name"`
		State                            string                          `json:"state"`
		Address                          string                          `json:"address"`
		Provider                         string                          `json:"provider"`
		TimeZone                         string                          `json:"time_zone"`
		PostalCode                       string                          `json:"postal_code"`
		Appointments                     []Appointment                   `json:"appointments"`
		ProviderBrand                    string                          `json:"provider_brand"`
		CarriesVaccine                   bool                            `json:"carries_vaccine"`
		AppointmentTypes                 map[AppointmentType]bool        `json:"appointment_types"`
		ProviderBrandID                  int                             `json:"provider_brand_id"`
		ProviderBrandName                string                          `json:"provider_brand_name"`
		ProviderLocationID               string                          `json:"provider_location_id"`
		AppointmentsAvailable            bool                            `json:"appointments_available"`
		AppointmentVaccineTypes          map[AppointmentVaccineType]bool `json:"appointment_vaccine_types"`
		AppointmentsLastFetched          time.Time                       `json:"appointments_last_fetched"`
		AppointmentsAvailableAllDoses    bool                            `json:"appointments_available_all_doses"`
		AppointmentsAvailable2NdDoseOnly bool                            `json:"appointments_available_2nd_dose_only"`
	} `json:"properties"`
}

type Appointment struct {
	Time             time.Time                `json:"time"`
	Type             string                   `json:"type"`
	VaccineTypes     []AppointmentVaccineType `json:"vaccine_types"`
	AppointmentTypes []AppointmentType        `json:"appointment_types"`
}

type AppointmentType string

const (
	AppointmentTypeUnknown        AppointmentType = "unknown"
	AppointmentTypeSecondDoseOnly AppointmentType = "2nd_dose_only"
	AppointmentTypeAllDoses       AppointmentType = "all_doses"
)

type AppointmentVaccineType string

const (
	AppointmentVaccineTypeUnknown AppointmentVaccineType = "unknown"
	AppointmentVaccineTypePfizer  AppointmentVaccineType = "pfizer"
	AppointmentVaccineTypeModerna AppointmentVaccineType = "moderna"
	AppointmentVaccineTypeJJ      AppointmentVaccineType = "jj"
)
