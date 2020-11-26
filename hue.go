package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type LightParameters struct {
	// SatStart start value for saturation (0-254)
	SatStart int `json:"sat_start"`
	// SatEnd end value for saturation (0-254)
	SatEnd int `json:"sat_end"`
	// BriStart start value for brightness (0-254)
	BriStart int `json:"bri_start"`
	// BriEnd end value for brightness (0-254)
	BriEnd int `json:"bri_end"`
	// HueStart start value for hue (0-65535)
	HueStart int `json:"hue_start"`
	// HueEnd end value for hue (0-65535)
	HueEnd int `json:"hue_end"`
}

// Payload to send to hue bridge
type Payload struct {
	// On toggle if light is in or off
	On bool `json:"on"`
	// Sat saturation of color
	Sat int `json:"sat"`
	// Bri brightness
	Bri int `json:"bri"`
	// Hue color dimension
	Hue int `json:"hue"`
}

const (
	user      = "pxillTVPe-aY0rRrS3cNL9n4eWPl7v6Gd2vMqvEH"
	light     = 1
	bridgeIP  = "192.168.178.174"
	intervals = 60
)

func doRequest(payload Payload) {

	// marshal Payload to jsonPayload
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// set the HTTP method, url, and request body
	url := fmt.Sprintf("http://%s/api/%s/lights/%d/state", bridgeIP, user, light)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for jsonPayload
	req.Header.Set("Content-Type", "application/jsonPayload; charset=utf-8")
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}

}
