package main

import (
	"fmt"
	"strings"
	"time"
)

var dayToString = map[time.Weekday]string{
	time.Monday:    "MON",
	time.Tuesday:   "TUE",
	time.Wednesday: "WED",
	time.Thursday:  "THU",
	time.Friday:    "FRI",
	time.Saturday:  "SAT",
	time.Sunday:    "SUN",
}

func cronExpression(settings Task) string {
	expression := fmt.Sprintf("%d %d * * %s", settings.Minutes, settings.Hours, concatWeekDays(settings.Days))
	fmt.Println(expression)
	return expression
}

func concatWeekDays(days []time.Weekday) string {
	var result strings.Builder
	//result.WriteString("\"")
	for i, day := range days {
		result.WriteString(dayToString[day])
		if i != len(days)-1 {
			result.WriteString(",")
		}
	}
	//result.WriteString("\"")
	return result.String()
}
