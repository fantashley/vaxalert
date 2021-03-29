package client

import "testing"

const vaxSpotterURL = "https://www.vaccinespotter.org/"

func TestGetStates(t *testing.T) {
	vaxClient, err := NewVaxClientV0(vaxSpotterURL, nil)
	if err != nil {
		t.Fatalf("Failed to get new VaxClientV0: %v", err)
	}

	states, err := vaxClient.GetStates()
	if err != nil {
		t.Errorf("Error getting providers: %v", err)
	}

	if len(states) == 0 {
		t.Fatal("No states received from GetStates()")
	}
}

func TestGetLocations(t *testing.T) {
	testStateCode := "MN"

	vaxClient, err := NewVaxClientV0(vaxSpotterURL, nil)
	if err != nil {
		t.Fatalf("Failed to get new VaxClientV0: %v", err)
	}

	locations, err := vaxClient.GetLocations(testStateCode)
	if err != nil {
		t.Errorf("Error getting locations: %v", err)
	}

	if len(locations.Locations) == 0 {
		t.Fatal("No locations received from GetLocations()")
	}
}

func TestLocationData(t *testing.T) {
	vaxClient, err := NewVaxClientV0(vaxSpotterURL, nil)
	if err != nil {
		t.Fatalf("Failed to get new VaxClientV0: %v", err)
	}

	states, err := vaxClient.GetStates()
	if err != nil {
		t.Errorf("Error getting providers: %v", err)
	}

	for _, state := range states {
		locations, err := vaxClient.GetLocations(state.Code)
		if err != nil {
			t.Errorf("Error getting locations: %v", err)
		}

		for _, loc := range locations.Locations {
			for _, appt := range loc.Properties.Appointments {
				if len(appt.AppointmentTypes) > 0 {
					t.Logf("%+v", appt.AppointmentTypes)
				}
			}
		}
	}
}
