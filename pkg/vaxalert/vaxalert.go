package vaxalert

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fantashley/vaxalert/pkg/vaxspotter"
	"github.com/hashicorp/go-multierror"
)

type VaxAlert struct {
	c         Config
	knownLocs map[int]int
}

func NewVaxAlert(c Config) (*VaxAlert, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("config failed validation: %w", err)
	}
	return &VaxAlert{
		c:         c,
		knownLocs: make(map[int]int),
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
			newLocs, allLocs := v.poll()
			if err := v.SendAlerts(ctx, newLocs, allLocs); err != nil {
				log.Printf("Error sending alerts: %v", err)
			}
		}
	}
}

func (v *VaxAlert) poll() (map[int]int, LocMap) {
	locations := make(LocMap)
	for _, state := range v.c.StateCodes {
		locs, err := v.c.VaxClient.GetLocations(state)
		if err != nil {
			log.Printf("failed to get locations in %s: %v", state, err)
			continue
		}
		for _, loc := range locs.Locations {
			locations[loc.Properties.ID] = loc
		}
	}

	currentLocs := make(map[int]int)
	for _, location := range locations {
		for _, rule := range v.c.Rules {
			apptCount := rule.FilterAppointments(location)
			if apptCount != 0 {
				currentLocs[location.Properties.ID] = apptCount
			}
		}
	}

	newLocs := make(map[int]int)
	for id, count := range currentLocs {
		if _, ok := v.knownLocs[id]; !ok {
			newLocs[id] = count
		}
	}

	v.knownLocs = currentLocs

	return newLocs, locations
}

func (v *VaxAlert) SendAlerts(ctx context.Context, newLocs map[int]int, allLocs LocMap) error {
	newCount := len(newLocs)
	if newCount == 0 {
		return nil
	}
	msg := formatMessage(newLocs, allLocs)
	var alertErr error
	for _, alerter := range v.c.Alerters {
		if err := alerter.Alert(ctx, msg); err != nil {
			alertErr = multierror.Append(alertErr, err)
		}
	}
	return alertErr
}

func formatMessage(newLocs map[int]int, allLocs LocMap) string {
	var sb strings.Builder
	for locID, apptCount := range newLocs {
		locObj := allLocs[locID]
		msg := fmt.Sprintf("%s in %s, %s %s has %d new appointments: %s\n",
			locObj.Properties.ProviderBrandName,
			locObj.Properties.City,
			locObj.Properties.State,
			locObj.Properties.PostalCode,
			apptCount,
			locObj.Properties.URL,
		)
		sb.WriteString(msg)
	}
	return sb.String()
}

type LocMap map[int]vaxspotter.Location
