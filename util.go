package main

import ics "github.com/arran4/golang-ical"

func withTzid(timezone string) ics.PropertyParameter {
	return &ics.KeyValues{
		Key:   string(ics.ParameterTzid),
		Value: []string{timezone},
	}
}

// func getTzone() ics.VTimezone {
//     tzone := ics.VTimezone{}
//     tzone.SetProperty(ics.ComponentProperty(ics.PropertyTzid), "Europe/Warsaw")
//
//     tzone_day := ics.Daylight{}
//     tzone_day.SetProperty(ics.ComponentProperty(ics.PropertyTzoffsetfrom), "+0100")
//     tzone_day.SetProperty(ics.ComponentProperty(ics.PropertyTzoffsetto), "+0200")
//     tzone_day.SetProperty(ics.ComponentProperty(ics.PropertyTzname), "CEST")
//     tzone_day.SetProperty(ics.ComponentPropertyDtStart, "19700329T020000")
//     tzone_day.SetProperty(ics.ComponentPropertyRrule, "FREQ=YEARLY;BYMONTH=3;BYDAY=-1SU")
//     tzone.Components = append(tzone.Components, &tzone_day)
//
//     tzone_std := ics.Standard{}
//     tzone_std.SetProperty(ics.ComponentProperty(ics.PropertyTzoffsetfrom), "+0100")
//     tzone_std.SetProperty(ics.ComponentProperty(ics.PropertyTzoffsetto), "+0200")
//     tzone_std.SetProperty(ics.ComponentProperty(ics.PropertyTzname), "CET")
//     tzone_std.SetProperty(ics.ComponentPropertyDtStart, "19701025T030000")
//     tzone_std.SetProperty(ics.ComponentPropertyRrule, "FREQ=YEARLY;BYMONTH=10;BYDAY=-1SU")
//     tzone.Components = append(tzone.Components, &tzone_std)
//
//     return tzone
// }
