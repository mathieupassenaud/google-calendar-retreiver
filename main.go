package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

var output map[string]string

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func main() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/calendar-go-quickstart.json
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}

	output = make(map[string]string)
	calendars, err := CalendarsFromFile("configuration.json")
	if err != nil {
		log.Fatalf("Unable to retrieve configuration file. %v", err)
	}

	for _, calendar := range calendars.Calendars {
		t := time.Now().Format(time.RFC3339)
		events, err := srv.Events.List(calendar.CalendarId).ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(1).OrderBy("startTime").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events. %v", err)
		}
		if len(events.Items) > 0 {
			for _, i := range events.Items {
				var when string
				var to string
				// If the DateTime is an empty string the Event is an all-day Event.
				// So only Date is available.

				if i.Start.DateTime != "" {
					when = i.Start.DateTime
					to = i.End.DateTime

					whenTime, err := time.Parse(time.RFC3339, when)
					if err != nil {
						log.Fatalf("unable to parse : ", err)
					}
					toTime, err1 := time.Parse(time.RFC3339, to)
					if err1 != nil {
						log.Fatalf("unable to parse : ", err1)
					}
					// check if time.Now is in event
					if inTimeSpan(whenTime, toTime, time.Now()) {
						output[calendar.OutputName] = i.Summary
					}
				} else {
					when = i.Start.Date
					// check if this is a current event
					if strings.Compare(when, time.Now().Local().Format("2000-12-30")) == 0 {
						output[calendar.OutputName] = i.Summary
					}
				}
			}
		}
	}
	outputString, _ := json.Marshal(output)
	fmt.Printf("%s", outputString)
}
