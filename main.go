// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var ipRegex = regexp.MustCompile(`\([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\)`)

func main() {
	if len(os.Args) <= 1 {
		panic("Please provide a destination IP")
	}

	dest := os.Args[1]
	if dest == "" {
		panic("Please provide a destination IP")
	}

	mapsApiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if mapsApiKey == "" {
		panic("Missing GOOGLE_MAPS_API_KEY environment variable. Snag one at https://developers.google.com/maps/documentation/javascript/get-api-key")
	}

	ips := traceRoute(dest)
	fmt.Println("Getting geos")

	geos := []geo{}
	for _, ip := range ips {
		ipGeo := getGeoForIp(ip)
		if ipGeo.Lat == 0 && ipGeo.Lng == 0 {
			fmt.Printf("Throwing out ip %s because its coords resolve to (0,0)\n", ip)
			continue // throw out the 0,0s
		}

		geos = append(geos, ipGeo)
	}

	if len(geos) == 0 {
		panic(fmt.Sprintf("Didn't get any location data on those IPs (%s) - try a different destination IP!", ips))
	}

	htmlFilePath := createHtml(formatGeos(geos), mapsApiKey)

	openFileInBrowser(htmlFilePath)
}

func extractIps(traceRouteOut string) []string {
	res := ipRegex.FindAllString(traceRouteOut, -1)
	for i, _ := range res {
		res[i] = strings.Replace(res[i], "(", "", -1) // hack because i don't feel like figuring out regex heh
		res[i] = strings.Replace(res[i], ")", "", -1)
	}

	return res
}

type geo struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lon"`
}

func getGeoForIp(ip string) geo {
	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	if err != nil {
		panic(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var r geo
	err = json.Unmarshal(bodyBytes, &r)
	if err != nil {
		panic(err)
	}

	return r
}

func formatGeos(geos []geo) (s string) {
	for _, g := range geos {
		s += fmt.Sprintf(`{lat: %v, lng: %v},`, g.Lat, g.Lng)
	}

	return s[:len(s)-1]
}

func createHtml(geos, mapsApiKey string) string {
	f, err := ioutil.TempFile("", "maproute")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	fileContents := template
	fileContents = strings.Replace(fileContents, "API_KEY_HERE", mapsApiKey, -1)
	fileContents = strings.Replace(fileContents, "GEOS_HERE", geos, -1)

	f.Write([]byte(fileContents))

	err = os.Rename(f.Name(), f.Name()+".html")
	if err != nil {
		panic(err)
	}

	return f.Name() + ".html"
}
