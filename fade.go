package main

import (
	"log"
	"time"
)

func fadeIn(settings Task) {
	sleepTime := time.Duration(settings.TimerDuration*60*1000/intervals) * time.Millisecond

	for i := 1; i <= intervals; i++ {

		lightParameters := settings.LightParameters

		payload := AmbientPayload{
			On:        true,
			Bri:       lightParameters.BriStart + updateValue(lightParameters.BriStart, lightParameters.BriEnd, i),
			ColorTemp: lightParameters.ColorTempStart + updateValue(lightParameters.ColorTempStart, lightParameters.ColorTempEnd, i),
		}
		log.Println(payload)
		doRequest(payload)
		time.Sleep(sleepTime)
	}
}

// calculate the update value for an iteration
func updateValue(start, end, iteration int) int {
	return int(float64(end-start) / float64(intervals) * float64(iteration))
}
