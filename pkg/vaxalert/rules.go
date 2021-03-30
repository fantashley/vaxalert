package vaxalert

import (
	"errors"
	"time"

	"github.com/fantashley/vaxalert/pkg/vaxspotter"
	"github.com/golang/geo/s2"
)

const EarthRadiusMiles = 3958.8

type AlertRule struct {
	StartDate        time.Time
	EndDate          time.Time
	Latitude         float64
	Longitude        float64
	MaxDistanceMiles int
	AppointmentType  vaxspotter.AppointmentType
	VaccineType      vaxspotter.AppointmentVaccineType
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

func (a AlertRule) FilterAppointments(loc vaxspotter.Location) ApptMap {
	matchingAppts := make(ApptMap)

	if a.MaxDistanceMiles != 0 {
		distance := a.getDistance(loc.Geometry.Coordinates)
		if a.MaxDistanceMiles < distance {
			return nil
		}
	}

	if a.AppointmentType != "" {
		if _, ok := loc.Properties.AppointmentTypes[a.AppointmentType]; !ok {
			return nil
		}
	}
	if a.VaccineType != "" {
		if _, ok := loc.Properties.AppointmentVaccineTypes[a.VaccineType]; !ok {
			return nil
		}
	}

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
		if a.VaccineType != "" {
			found := false
			for _, apptType := range appt.VaccineTypes {
				if a.VaccineType == apptType {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if !a.evaluateTime(appt.Time) {
			continue
		}
		matchingAppts[getApptIdent(appt, loc)] = appt
	}

	return matchingAppts
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
