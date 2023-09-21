package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	ics "github.com/arran4/golang-ical"
)

type EventURLsXML struct {
	URLs []string `xml:"response>href"`
}

func deleteEvent(url1, url2 string, conf WebDAVConf, wg *sync.WaitGroup) {
	defer wg.Done()

	url := url1 + url2
	client := http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Error creating new request:", err)
		return
	}

	req.SetBasicAuth(conf.User, conf.Pass)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		fmt.Printf("Request failed with status code %d\n", res.StatusCode)
		return
	}
}

func clearCalendar(conf WebDAVConf, start time.Time, end time.Time) error {
	fmt.Println("👀 Retriving old events from server")

	url := conf.URL + "/" + conf.CalendarName

	startStr := start.Format("20060102T150405")
	endStr := end.Format("20060102T150405")
	body := fmt.Sprintf(`<?xml version='1.0' encoding='utf-8' ?>
    <C:calendar-query xmlns:C='urn:ietf:params:xml:ns:caldav'>
        <D:prop xmlns:D='DAV:'>
            <D:getetag />
        </D:prop>
        <C:filter>
            <C:comp-filter name='VCALENDAR'>
                <C:comp-filter name='VEVENT'>
                    <C:time-range start='%v' end='%v'/>
                </C:comp-filter>
            </C:comp-filter>
        </C:filter>
    </C:calendar-query>
    `, startStr, endStr)

	req, err := http.NewRequest("REPORT", url, strings.NewReader(body))
	if err != nil {
		fmt.Println("Error creating new request:", err)
		return err
	}

	req.SetBasicAuth(conf.User, conf.Pass)
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("Depth", "1")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body\n")
		return err
	}

	if res.StatusCode != http.StatusMultiStatus {
		fmt.Printf("Request failed with status code %d\n", res.StatusCode)
		return err
	}

	fmt.Println("📝 Parsing data")
	urls := EventURLsXML{}
	err = xml.Unmarshal(resBody, &urls)
	if err != nil {
		fmt.Println("Error parsing xml")
		return err
	}

	fmt.Println("❌ Removing old events")
	url = conf.UrlObj.Scheme + "://" + conf.UrlObj.Host
	wg := sync.WaitGroup{}
	for _, url2 := range urls.URLs {
		wg.Add(1)
		go deleteEvent(url, url2, conf, &wg)
	}

	wg.Wait()
	return nil
}

func addEvent(conf WebDAVConf, timezone string, event *ics.VEvent, wg *sync.WaitGroup) error {
	defer wg.Done()
	// tzone := getTzone()

	cal := ics.NewCalendar()

	cal.AddVEvent(event)
	cal.SetProductId("CalGen Script")
	cal.SetCalscale("GREGORIAN")
	// cal.Components = append(cal.Components, &tzone)

	// fmt.Println(cal.Serialize())
	url := conf.URL + "/" + conf.CalendarName + "/" + event.Id() + ".ics"
	req, err := http.NewRequest("PUT", url, strings.NewReader(cal.Serialize()))
	if err != nil {
		fmt.Println("Error creating new request:", err)
		return err
	}

	req.SetBasicAuth(conf.User, conf.Pass)
	req.Header.Add("Content-Type", "text/calendar")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		fmt.Printf("Request failed with status code %d\n", res.StatusCode)
		return err
	}

	return nil
}
