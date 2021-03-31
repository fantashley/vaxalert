# VaxAlert

## Thorough Instructions

### Installation

`go get github.com/fantashley/vaxalert/cmd/vaxalert`

### Configuration

#### Options

```bash
$ vaxalert -h
Usage of vaxalert:
  -alert-numbers string
        comma-separated numbers to send SMS alerts to when appointments are found
  -api-url string
        API URL of vaccine spotter (default "https://www.vaccinespotter.org/")
  -appt-end-date string
        end of date range for searching appointments
  -appt-start-date string
        start of date range for searching appointments
  -config string
        config file (optional)
  -latitude float
        latitude of coordinate to search around
  -longitude float
        longitude of coordinate to search around
  -max-distance int
        maximum distance in miles from coordinates to search for appointments
  -poll-interval duration
        new appointment polling frequency (default 5m0s)
  -second-dose
        only search for appointments for second dose
  -state-codes string
        comma-separated list of state codes
  -twilio-account-sid string
        Twilio account sid for SMS alerting
  -twilio-auth-token string
        Twilio auth token for SMS alerting
  -twilio-msg-sid string
        Twilio messaging sid
  -vaccine-type string
        type of vaccine to search for in appointments
```

#### config.json

```json
{
  "poll-interval": "30s",
  "state-codes": "MN",
  "twilio-account-sid": "<your sid here>",
  "twilio-auth-token": "<your token here>",
  "twilio-msg-sid": "<your msg sid here>",
  "alert-numbers": "+16516463003",
  "appt-start-date": "03/30/2021",
  "appt-end-date": "04/06/2021",
  "latitude": "44.9658649",
  "longitude": "-93.2354132",
  "max-distance": 50,
  "second-dose": false,
  "vaccine-type": "pfizer"
}
```

### Execution

`vaxalert -config config.json`
