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
  OtherCals:
    - other.ical

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

## Config options

#### File

a group of options related to saving calendar to file on disk.

```yaml
Config:
  File:
    Enable: false 
    Name: cal.ical
```

- **Enable** - true if you want to save data to the file or false if you don't. Defaults to false.
- **Name** - name of the calendar file

#### WebDAV

a group of options needed to connect to the WebDAV server

```yaml
Config:
  WebDAV:
    Enable: true 
    URL: https://yourCalDavServer.com/remote.php/dav/calendars/YOUR_USER
    CalendarName: YOUR_CALENDAR
    User: YOUR_USER
    Pass: YOUR_PASSWORD
```

- **Enable** - true if you want to save data on the remote server or false if you don't. Defaults to false.
- **URL** - URL to caldav server
- **CalendarName** - name of the calendar to which cal_gen will add events
- **User** - your username
- **Pass** - your password in plain text (needs to be changed to sth more secure)

#### OtherCals 

optional list of paths to other ics/ical files which should be added to the timetable

```yaml
Config:
  OtherCals:
  - cal1.ical
  - cal2.ical
```

## Development && Building

```sh 
git clone https://github.com/Plasny/CalGen.git
cd CalGen
go build .
```

*If you have any questions reach me out via email pp.git@plasny.one*

