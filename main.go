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
	"io/ioutil"
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

	geos := traceRoute(dest)

	filteredGeos := []geo{}
	for _, g := range geos {
		if g.Lat == 0 && g.Lng == 0 {
			continue
		}
		filteredGeos = append(filteredGeos, g)
	}

	if len(filteredGeos) == 0 {
		panic("Didn't get any location data on those IPs - try a different destination IP!")
	}

	htmlFilePath := createHtml(formatGeos(filteredGeos), mapsApiKey)
	openFileInBrowser(htmlFilePath)
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
