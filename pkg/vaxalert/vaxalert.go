package vaxalert

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fantashley/vaxalert/pkg/vaxspotter"
)

type VaxAlert struct {
	c          Config
	knownAppts map[time.Time]vaxspotter.Appointment
}

func NewVaxAlert(c Config) (VaxAlert, error) {
	if err := c.Validate(); err != nil {
		return VaxAlert{}, fmt.Errorf("config failed validation: %w", err)
	}
	return VaxAlert{
		c:          c,
		knownAppts: make(map[time.Time]vaxspotter.Appointment),
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

		}
	}
}

func (v VaxAlert) poll() {
	var locations []vaxspotter.Location
	for _, state := range v.c.StateCodes {
		locs, err := v.c.VaxClient.GetLocations(state)
		if err != nil {
			log.Printf("failed to get locations in %s: %v", state, err)
		}
	}
}
