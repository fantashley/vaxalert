package vaxalert

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/fantashley/vaxalert/pkg/vaxspotter"
	"github.com/hashicorp/go-multierror"
)

type VaxAlert struct {
	c          Config
	knownAppts ApptMap
	knownLocs  LocMap
}

func NewVaxAlert(c Config) (*VaxAlert, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("config failed validation: %w", err)
	}
	return &VaxAlert{
		c:          c,
		knownAppts: make(ApptMap),
		knownLocs:  make(LocMap),
	}, nil
}

func (v *VaxAlert) Start(ctx context.Context) error {
	pollTicker := time.NewTicker(v.c.PollInterval)
	defer pollTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-pollTicker.C:
			newAppts, newLocs := v.poll()
			if err := v.SendAlerts(ctx, newAppts, newLocs); err != nil {
				log.Printf("Error sending alerts: %v", err)
			}
		}
	}
}

func (v *VaxAlert) poll() (ApptMap, LocMap) {
	var locations []vaxspotter.Location
	for _, state := range v.c.StateCodes {
		locs, err := v.c.VaxClient.GetLocations(state)
		if err != nil {
			log.Printf("failed to get locations in %s: %v", state, err)
			continue
		}
		locations = append(locations, locs.Locations...)
	}

	currentAppts := make(ApptMap)
	currentLocs := make(LocMap)
	for _, location := range locations {
		for _, rule := range v.c.Rules {
			matchingAppts, matchingLocs := rule.FilterAppointments(location)
			for ident, appt := range matchingAppts {
				currentAppts[ident] = appt
			}
			for id, loc := range matchingLocs {
				currentLocs[id] = loc
			}
		}
	}

	newAppts := make(ApptMap)
	for ident, appt := range currentAppts {
		if _, ok := v.knownAppts[ident]; !ok {
			newAppts[ident] = appt
		}
	}

	newLocs := make(LocMap)
	for id, loc := range currentLocs {
		if _, ok := v.knownLocs[id]; !ok {
			newLocs[id] = loc
		}
	}

	v.knownAppts = currentAppts
	v.knownLocs = currentLocs

	return newAppts, newLocs
}

func (v *VaxAlert) SendAlerts(ctx context.Context, newAppts ApptMap, newLocs LocMap) error {
	newCount := len(newAppts) + len(newLocs)
	if newCount == 0 {
		return nil
	}
	msg := fmt.Sprintf("%d new appointments found!", newCount)
	var alertErr error
	for _, alerter := range v.c.Alerters {
		if err := alerter.Alert(ctx, msg); err != nil {
			alertErr = multierror.Append(alertErr, err)
		}
	}
	return alertErr
}

type ApptMap map[ApptIdent]vaxspotter.Appointment
type LocMap map[int]vaxspotter.Location

type ApptIdent string

func getApptIdent(appt vaxspotter.Appointment, loc vaxspotter.Location) ApptIdent {
	return ApptIdent(appt.Time.String() + strconv.Itoa(loc.Properties.ID))
}
