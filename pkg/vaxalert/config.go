package vaxalert

import (
	"errors"
	"fmt"
	"time"

	"github.com/fantashley/vaxalert/pkg/client"
)

type Config struct {
	VaxClient    *client.VaxClientV0
	PollInterval time.Duration
	Rules        []AlertRule
	StateCodes   []string
}

func (c Config) Validate() error {
	if c.VaxClient == nil {
		return errors.New("vaxclient is nil")
	}
	for i, rule := range c.Rules {
		if err := rule.Validate(); err != nil {
			return fmt.Errorf("alert rule %d failed validation: %w", i, err)
		}
	}
	for _, code := range c.StateCodes {
		if len(code) != 2 {
			return fmt.Errorf("invalid state code %q", code)
		}
	}

	return nil
}
