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
)

type geo struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lon"`
	City    string  `json:"city"`
	Country string  `json:"country"`
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

	if r.Lat == 0 && r.Lng == 0 {
		r.Country = "Unknown Country"
		r.City = "Unknown City"
	}

	return r
}

func formatGeos(geos []geo) (s string) {
	for _, g := range geos {
		s += fmt.Sprintf(`{lat: %v, lng: %v},`, g.Lat, g.Lng)
	}

	return s[:len(s)-1]
}
