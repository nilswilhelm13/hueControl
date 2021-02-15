package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Light struct {
	State struct {
		On        bool   `json:"on"`
		Bri       int    `json:"bri"`
		Ct        int    `json:"ct"`
		Alert     string `json:"alert"`
		Colormode string `json:"colormode"`
		Mode      string `json:"mode"`
		Reachable bool   `json:"reachable"`
	} `json:"state"`
	Swupdate struct {
		State       string `json:"state"`
		Lastinstall string `json:"lastinstall"`
	} `json:"swupdate"`
	Type             string `json:"type"`
	Name             string `json:"name"`
	Modelid          string `json:"modelid"`
	Manufacturername string `json:"manufacturername"`
	Productname      string `json:"productname"`
	Capabilities     struct {
		Certified bool `json:"certified"`
		Control   struct {
			Mindimlevel int `json:"mindimlevel"`
			Maxlumen    int `json:"maxlumen"`
			Ct          struct {
				Min int `json:"min"`
				Max int `json:"max"`
			} `json:"ct"`
		} `json:"control"`
		Streaming struct {
			Renderer bool `json:"renderer"`
			Proxy    bool `json:"proxy"`
		} `json:"streaming"`
	} `json:"capabilities"`
	Config struct {
		Archetype string `json:"archetype"`
		Function  string `json:"function"`
		Direction string `json:"direction"`
		Startup   struct {
			Mode       string `json:"mode"`
			Configured bool   `json:"configured"`
		} `json:"startup"`
	} `json:"config"`
	Uniqueid  string `json:"uniqueid"`
	Swversion string `json:"swversion"`
	Productid string `json:"productid"`
}

type ColorLightParameters struct {
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

type AmbientLightParameters struct {
	// BriStart start value for brightness (0-254)
	BriStart int `json:"bri_start"`
	// BriEnd end value for brightness (0-254)
	BriEnd int `json:"bri_end"`
	// ColorTempStart start value for hue (153-454)
	ColorTempStart int `json:"ct_start"`
	// ColorTempEnd end value for hue (153-454)
	ColorTempEnd int `json:"ct_end"`
}

// ColorPayload to send to hue bridge
type ColorPayload struct {
	// On toggle if light is in or off
	On bool `json:"on"`
	// Sat saturation of color
	Sat int `json:"sat"`
	// Bri brightness
	Bri int `json:"bri"`
	// Hue color dimension
	Hue int `json:"hue"`
}

// AmbientPayload to send to hue bridge
type AmbientPayload struct {
	// On toggle if light is in or off
	On bool `json:"on"`
	// Bri brightness
	Bri int `json:"bri"`
	// ColorTemp
	ColorTemp int `json:"ct"`
}

type Lights map[string]Light

const (
	user      = "WRDMUd1yfs1kcbdOzyhSlN727d02U1BWHBLMuwVP"
	light     = 4
	bridgeIP  = "192.168.178.174"
	intervals = 60
	groupId   = 3
)

func doRequest(payload AmbientPayload) {

	// marshal ColorPayload to jsonPayload
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// set the HTTP method, url, and request body
	url := fmt.Sprintf("http://%s/api/%s/groups/%d/action", bridgeIP, user, groupId)
	log.Printf("URL: %s", url)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		panic(err)
	}
	// set the request header Content-Type for jsonPayload
	req.Header.Set("Content-Type", "application/jsonPayload; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	log.Printf("Status Code: %d", resp.StatusCode)

}

func getLights() (Lights, error) {
	url := fmt.Sprintf("http://%s/api/%s/lights", bridgeIP, user)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// set the request header Content-Type for jsonPayload
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	ids := topLevelKeys(body)
	lights, err := parseLights(body, ids)
	if err != nil {
		return nil, err
	}
	return lights, nil
}

func parseLights(responseBody []byte, ids []string) (Lights, error) {
	var body map[string]interface{}
	err := json.Unmarshal(responseBody, &body)
	if err != nil {
		return nil, err
	}
	lights := make(Lights)
	for _, id := range ids {
		var light Light
		lightJSON, err := json.Marshal(body[id])
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(lightJSON, &light)
		if err != nil {
			return nil, err
		}
		lights[id] = light
	}
	return lights, nil
}

func topLevelKeys(j []byte) []string {
	// a map container to decode the JSON structure into
	c := make(map[string]json.RawMessage)

	// unmarschal JSON
	e := json.Unmarshal(j, &c)

	// panic on error
	if e != nil {
		panic(e)
	}

	// a string slice to hold the keys
	k := make([]string, len(c))

	// iteration counter
	i := 0

	// copy c's keys into k
	for s, _ := range c {
		k[i] = s
		i++
	}

	return k
}
