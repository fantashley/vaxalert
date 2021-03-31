package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fantashley/vaxalert/pkg/client"
	"github.com/fantashley/vaxalert/pkg/vaxalert"
	"github.com/fantashley/vaxalert/pkg/vaxspotter"
	"github.com/peterbourgon/ff"
)

func main() {
	fs := flag.NewFlagSet("vaxalert", flag.ExitOnError)
	var (
		_            = fs.String("config", "", "config file (optional)")
		apiURL       = fs.String("api-url", "https://www.vaccinespotter.org/", "API URL of vaccine spotter")
		pollInterval = fs.Duration("poll-interval", 5*time.Minute, "new appointment polling frequency")
		stateCodes   = fs.String("state-codes", "", "comma-separated list of state codes to search in")

		twilioAccountSid   = fs.String("twilio-account-sid", "", "Twilio account sid for SMS alerting")
		twilioAuthToken    = fs.String("twilio-auth-token", "", "Twilio auth token for SMS alerting")
		twilioMessagingSid = fs.String("twilio-msg-sid", "", "Twilio messaging sid")
		alertNumbers       = fs.String("alert-numbers", "", "comma-separated numbers to send SMS alerts to when appointments are found")

		apptStartDate = fs.String("appt-start-date", "", "start of date range for searching appointments")
		apptEndDate   = fs.String("appt-end-date", "", "end of date range for searching appointments")
		latitude      = fs.Float64("latitude", 0, "latitude of coordinate to search around")
		longitude     = fs.Float64("longitude", 0, "longitude of coordinate to search around")
		maxDistance   = fs.Int("max-distance", 0, "maximum distance in miles from coordinates to search for appointments")

		secondDoseOnly = fs.Bool("second-dose", false, "only search for appointments for second dose")
		vaccineType    = fs.String("vaccine-type", "", "(not required) type of vaccine to search for in appointments")
	)

	ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("VAX"),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.JSONParser),
	)

	alerter := vaxalert.NewTwilioAlerter(
		*twilioAccountSid,
		*twilioAuthToken,
		*twilioMessagingSid,
		&http.Client{},
		strings.Split(*alertNumbers, ","),
	)

	var err error

	startTime := time.Time{}
	if *apptStartDate != "" {
		startTime, err = time.Parse("01/02/2006", *apptStartDate)
		if err != nil {
			log.Panicf("Failed to parse start time: %v", err)
		}
	}

	endTime := time.Time{}
	if *apptEndDate != "" {
		endTime, err = time.Parse("01/02/2006", *apptEndDate)
		if err != nil {
			log.Panicf("Failed to parse end time: %v", err)
		}
		endTime = endTime.Add(24 * time.Hour) // inclusive of end day
	}

	var apptType vaxspotter.AppointmentType
	if *secondDoseOnly {
		apptType = vaxspotter.AppointmentTypeSecondDoseOnly
	}

	rule := vaxalert.AlertRule{
		StartDate:        startTime,
		EndDate:          endTime,
		Latitude:         *latitude,
		Longitude:        *longitude,
		MaxDistanceMiles: *maxDistance,
		AppointmentType:  apptType,
		VaccineType:      vaxspotter.AppointmentVaccineType(*vaccineType),
	}

	vaxClient, err := client.NewVaxClientV0(*apiURL, &http.Client{})
	if err != nil {
		log.Panicf("Failed to get vaccine client: %v", err)
	}

	alertCfg := vaxalert.Config{
		VaxClient:    vaxClient,
		PollInterval: *pollInterval,
		Rules:        []vaxalert.AlertRule{rule},
		StateCodes:   strings.Split(*stateCodes, ","),
		Alerters:     []vaxalert.Alerter{alerter},
	}

	alertObj, err := vaxalert.NewVaxAlert(alertCfg)
	if err != nil {
		log.Panicf("Failed to initialize vaccine alerter: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		cancel()
	}()

	if err = alertObj.Start(ctx); err != nil {
		log.Panicf("VaxAlert returned an error: %v", err)
	}
}
