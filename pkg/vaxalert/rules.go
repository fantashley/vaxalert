package vaxalert

import (
	"errors"
	"time"

	"github.com/fantashley/vaxalert/pkg/vaxspotter"
	"github.com/golang/geo/s2"
)

const EarthRadiusMiles = 3958.8
const AppointmentsUnknown = -1

type AlertRule struct {
	StartDate        time.Time
	EndDate          time.Time
	Latitude         float64
	Longitude        float64
	MaxDistanceMiles int
	AppointmentType  vaxspotter.AppointmentType
	VaccineTypes     []vaxspotter.AppointmentVaccineType
}

func (a AlertRule) Validate() error {
	zeroTime := a.StartDate.IsZero() && a.EndDate.IsZero()
	zeroLatLng := a.Latitude == 0 && a.Longitude == 0
	zeroDist := a.MaxDistanceMiles == 0

	err := errors.New("invalid alert rule")
	switch {
	case zeroLatLng != zeroDist:
		return err
	case zeroTime && zeroLatLng && zeroDist:
		return err
	}

	return nil
}

func (a AlertRule) FilterAppointments(loc vaxspotter.Location) int {
	if a.MaxDistanceMiles != 0 {
		distance := a.getDistance(loc.Geometry.Coordinates)
		if a.MaxDistanceMiles < distance {
			return 0
		}
	}

	if loc.Properties.AppointmentsAvailable && (len(loc.Properties.Appointments) == 0) {
		return AppointmentsUnknown
	}

	if a.AppointmentType != "" {
		if _, ok := loc.Properties.AppointmentTypes[a.AppointmentType]; !ok {
			return 0
		}
	}
	if len(a.VaccineTypes) > 0 {
		found := false
		for _, vType := range a.VaccineTypes {
			if _, ok := loc.Properties.AppointmentVaccineTypes[vType]; ok {
				found = true
				break
			}
		}
		if !found {
			return 0
		}
	}

	apptCount := 0
	for _, appt := range loc.Properties.Appointments {
		if a.AppointmentType != "" {
			found := false
			for _, apptType := range appt.AppointmentTypes {
				if a.AppointmentType == apptType {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if len(a.VaccineTypes) > 0 {
			found := false
		apptLoop:
			for _, apptType := range appt.VaccineTypes {
				for _, vType := range a.VaccineTypes {
					if apptType == vType {
						found = true
						break apptLoop
					}
				}
			}
			if !found {
				continue
			}
		}
		if !a.evaluateTime(appt.Time) {
			continue
		}
		apptCount++
	}

	return apptCount
}

func (a AlertRule) getDistance(coords []float64) int {
	pointA := s2.LatLngFromDegrees(a.Latitude, a.Longitude)
	pointB := s2.LatLngFromDegrees(coords[1], coords[0])
	distance := pointA.Distance(pointB)
	return int(distance * EarthRadiusMiles)
}

func (a AlertRule) evaluateTime(apptTime time.Time) bool {
	if a.StartDate.IsZero() && a.EndDate.IsZero() {
		return true
	}
	startTime := time.Now()
	if !a.StartDate.IsZero() {
		startTime = a.StartDate
	}
	if apptTime.Before(startTime) {
		return false
	}
	endTime := time.Now().AddDate(0, 3, 0)
	if !a.EndDate.IsZero() {
		endTime = a.EndDate
	}
	if apptTime.After(endTime) {
		return false
	}
	return true
}
