package main

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const DEFAULT_MARKER = "## generated with cal_gen ##"
const DEFAULT_SEPARATOR = "\\n\\n"

type Event struct {
	From        time.Time
	To          time.Time
	Title       string
	Description string
}

type WebDAVConf struct {
	Enable       bool   `yaml:"Enable"`
	URL          string `yaml:"URL"`
	UrlObj       *url.URL
	CalendarName string `yaml:"CalendarName"`
	User         string `yaml:"User"`
	Pass         string `yaml:"Pass"`
}

type FileConf struct {
	Enable bool   `yaml:"Enable"`
	Name   string `yaml:"Name"`
}

type CalConf struct {
	Config struct {
		Timespan struct {
			From time.Time `yaml:"From"`
			To   time.Time `yaml:"To"`
		} `yaml:"Timespan"`
		TimeZone string     `yaml:"TimeZone"`
		WebDAV   WebDAVConf `yaml:"WebDAV"`
		File     FileConf   `yaml:"File"`
		Marker    string     `yaml:"Marker"`
	} `yaml:"Config"`
	Week struct {
		Sunday    []string `yaml:"Sunday"`
		Monday    []string `yaml:"Monday"`
		Tuesday   []string `yaml:"Tuesday"`
		Wednesday []string `yaml:"Wednesday"`
		Thursday  []string `yaml:"Thursday"`
		Friday    []string `yaml:"Friday"`
		Saturday  []string `yaml:"Saturday"`
	} `yaml:"Week"`
}

func main() {
	var cal *ics.Calendar

	dateZero := time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)

	yamlFile, err := os.Open("./cal.yaml")
	if err != nil {
		println("Error opening the file.")
		return
	}
	defer yamlFile.Close()

	decoder := yaml.NewDecoder(yamlFile)

	for {
		yaml := CalConf{}
		table := [7][]Event{}

		if err := decoder.Decode(&yaml); err != nil {
			break
		}

		fmt.Println("👀 Reading yaml file")

		if yaml.Config.Marker == "" {
			yaml.Config.Marker = DEFAULT_MARKER
		}

		yaml.Config.WebDAV.UrlObj, err = url.Parse(yaml.Config.WebDAV.URL)
		if err != nil {
			fmt.Println("Error parsing webdav url")
			return
		}

		t := reflect.TypeOf(yaml.Week)
		for i := 0; i < t.NumField(); i++ {
			for _, str := range reflect.ValueOf(yaml.Week).Field(i).Interface().([]string) {
				event := Event{}
				arr := strings.Split(str, "-")

				event.From, err = time.Parse("15:04", strings.Trim(arr[0], " "))
				if err != nil {
					fmt.Println("sth went wrong while parsing time")
					return
				}

				event.To, err = time.Parse("15:04", strings.Trim(arr[1], " "))
				if err != nil {
					fmt.Println("sth went wrong while parsing time")
					return
				}

				event.Title = strings.Trim(arr[2], " ")

				if len(arr) > 3 {
					event.Description = strings.Trim(arr[3], " ")
				}

				// fmt.Printf("Event: %v\n  day: %v\n  from: %v\n  till: %v\n", event.Title,
				// 	t.Field(i).Name, event.From.Format("15:04"), event.To.Format("15:04"))

				table[i] = append(table[i], event)
			}
		}

		// for _, day := range table {
		// 	for _, event := range day {
		// 		fmt.Printf("Event: %v\n  from: %v\n  till: %v\n", event.Title,
		// 			event.From.Format("15:04"), event.To.Format("15:04"))
		// 	}
		// }

		if yaml.Config.File.Enable {
			cal = ics.NewCalendar()
			cal.SetTimezoneId(yaml.Config.TimeZone)
			cal.SetProductId("cal_gen")
			cal.SetLastModified(time.Now())
		}

		if yaml.Config.WebDAV.Enable {
			clearCalendar(yaml.Config.WebDAV, yaml.Config.Timespan.From, yaml.Config.Timespan.To, yaml.Config.Marker)
		}

		if yaml.Config.WebDAV.Enable {
			fmt.Println("➕ Building timetable and adding events to caldav server")
		} else {
			fmt.Println("🎉 Building timetable for given time interval")
		}

		wg := sync.WaitGroup{}

		for currentDate := yaml.Config.Timespan.From; yaml.Config.Timespan.To.Compare(currentDate) > 0; currentDate = currentDate.Add(time.Hour * 24) {
			weekday := int(currentDate.Weekday())
			for _, eventData := range table[weekday] {
				event := ics.NewEvent(uuid.New().String())

				tzid := withTzid(yaml.Config.TimeZone)

				duration := eventData.From.Sub(dateZero)
				event.SetProperty(ics.ComponentPropertyDtStart, currentDate.Add(duration).Format("20060102T150405"), tzid)

				duration = eventData.To.Sub(dateZero)
				event.SetProperty(ics.ComponentPropertyDtEnd, currentDate.Add(duration).Format("20060102T150405"), tzid)

				event.SetSummary(eventData.Title)

				if eventData.Description != "" {
					event.SetDescription(eventData.Description + DEFAULT_SEPARATOR + yaml.Config.Marker)
				} else {
					event.SetDescription(yaml.Config.Marker)
				}

				if yaml.Config.WebDAV.Enable {
					wg.Add(1)
					go addEvent(yaml.Config.WebDAV, yaml.Config.TimeZone, event, &wg)
				}

				if yaml.Config.File.Enable {
					cal.AddVEvent(event)
				}
			}
		}

		if yaml.Config.File.Enable {
			icalFile, err := os.Create(yaml.Config.File.Name)
			if err != nil {
				fmt.Println("Error opening the file.")
				return
			}
			defer icalFile.Close()

			fmt.Println("🗓️ Creating calendar file")
			cal.SerializeTo(icalFile)
		}

		wg.Wait()

		fmt.Println("✅ Your calendar is ready now")
	}
}
