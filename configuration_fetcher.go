package main

import (
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
)

type Calendars struct {
	Calendars []CalendarOutput
}

type CalendarOutput struct {
	CalendarId string
	OutputName string
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func CalendarsFromFile(file string) (*Calendars, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &Calendars{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}
