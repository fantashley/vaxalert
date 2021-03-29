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
}

func NewVaxAlert(c Config) (VaxAlert, error) {
	if err := c.Validate(); err != nil {
		return VaxAlert{}, fmt.Errorf("config failed validation: %w", err)
	}
	return VaxAlert{
		c:          c,
		knownAppts: make(ApptMap),
	}, nil
}

func (v VaxAlert) Start(ctx context.Context) error {
	pollTicker := time.NewTicker(v.c.PollInterval)
	defer pollTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-pollTicker.C:
			newAppts := v.poll()
			if err := v.SendAlerts(ctx, newAppts); err != nil {
				log.Printf("Error sending alerts: %v", err)
			}
		}
	}
}

func (v VaxAlert) poll() ApptMap {
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
	for _, location := range locations {
		for _, rule := range v.c.Rules {
			for ident, appt := range rule.FilterAppointments(location) {
				currentAppts[ident] = appt
			}
		}
	}

	newAppts := make(ApptMap)
	for ident, appt := range currentAppts {
		if _, ok := v.knownAppts[ident]; !ok {
			newAppts[ident] = appt
		}
	}

	v.knownAppts = currentAppts

	return newAppts
}

func (v VaxAlert) SendAlerts(ctx context.Context, newAppts ApptMap) error {
	var alertErr error
	msg := fmt.Sprintf("%d new appointments found!", len(newAppts))
	for _, alerter := range v.c.Alerters {
		if err := alerter.Alert(ctx, msg); err != nil {
			alertErr = multierror.Append(alertErr, err)
		}
	}
	return alertErr
}

type ApptMap map[ApptIdent]vaxspotter.Appointment

type ApptIdent string

func getApptIdent(appt vaxspotter.Appointment, loc vaxspotter.Location) ApptIdent {
	return ApptIdent(appt.Time.String() + strconv.Itoa(loc.Properties.ID))
}
