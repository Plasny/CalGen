# Cal Gen

Small go program for creating calendar files and modifing events on
caldav server according to a given schedule. I created it to ease
[timeblocking](https://todoist.com/productivity-methods/time-blocking)
preparation.

## Usage

To use this program you need to create a `cal.yaml` file. It should
contain data in format like below. Using it you can define whether to save
calendar file locally or to use a WebDAV server.

``` yaml
Config:
  TimeZone: YOUR_TIMEZONE   # eg.Europe/Warsaw
  Timespan:
    From: 2023-09-21
    To: 2023-11-01
  File:
    Enable: false 
    Name: cal.ical
  WebDAV:
    Enable: true 
    URL: https://yourCalDavServer.com/remote.php/dav/calendars/YOUR_USER
    CalendarName: YOUR_CALENDAR
    User: YOUR_USER
    Pass: YOUR_PASSWORD

Week:
  Monday:
    - 7:30-14:40 - some event - event description
    - 15:00-17:00 - some event
  
  Tuesday:
    - 8:25-12:55 - some event - event description
    - 13:30-14:15 - some event - event description
    - 15:00-17:00 - some event
  
  Wednesday:
    - 9:20-14:40 - some event
    - 15:00-17:00 - some event - event description
  
  Thursday:
    - 8:15-9:13 - some event - event description
    - 9:20-12:55 - some event
    - 13:30-14:15 - some event - event description
  
  Friday:
    - 7:00-12:55 - some event - event description
    - 17:00-21:00 - some event

  Saturday:
    - 9:00-12:00 - some event

  Sunday:
    - 9:00-10:00 - some event - event description
```

## Development && Building

```sh 
git clone https://github.com/Plasny/CalGen.git
cd CalGen
go build .
```

*If you have any questions reach me out via email pp.git@plasny.one*

