package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Expects something like this:
// $ traceroute 162.144.216.152
// traceroute to 162.144.216.152 (162.144.216.152), 64 hops max, 52 byte packets
// 1  100.76.175.252 (100.76.175.252)  2.504 ms  2.090 ms  1.607 ms
// 2  us-sfo-spe-bb1-xe-4-2-0-0.n.corp.google.com (172.31.151.10)  1.475 ms  1.854 ms  1.870 ms
// 3  us-sfo-spe-fw1-reth0-503.n.corp.google.com (172.31.151.177)  2.428 ms  2.255 ms  2.114 ms
// 4  104.132.11.193 (104.132.11.193)  2.369 ms  2.415 ms  2.137 ms
// 5  pr01-xe-4-2-1-511.pao03.net.google.com (72.14.210.179)  10.286 ms  3.437 ms  3.325 ms
// 6  pr02-ae2.sjc07.net.google.com (108.170.242.231)  4.359 ms
// pr02-ae6.sjc07.net.google.com (108.170.243.7)  4.735 ms  5.891 ms
// 7  te2-0-0d0.cir1.sanjose2-ca.us.xo.net (216.156.84.29)  5.824 ms  4.533 ms  4.543 ms
// 8  216.156.16.194.ptr.us.xo.net (216.156.16.194)  21.200 ms  21.806 ms  21.348 ms
// 9  207.88.12.150.ptr.us.xo.net (207.88.12.150)  21.335 ms  22.661 ms  22.585 ms
// 10  207.88.12.159.ptr.us.xo.net (207.88.12.159)  21.850 ms  21.555 ms  21.211 ms
// 11  te-4-1-0.rar3.miami-fl.us.xo.net (207.88.12.161)  21.259 ms  22.038 ms  22.972 ms
// 12  216.156.16.29.ptr.us.xo.net (216.156.16.29)  25.705 ms  36.554 ms  22.068 ms
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
			continue // throw out the 0,0s
		}

		geos = append(geos, ipGeo)
	}

	htmlFilePath := createHtml(formatGeos(geos), mapsApiKey)

	openFileInBrowser(htmlFilePath)
}

func traceRoute(destinationIp string) []string {
	fmt.Println("Tracing")

	cmd := "traceroute"
	args := []string{destinationIp}
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		panic(err)
	}

	return extractIps(string(out))
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

	err = os.Rename(f.Name(), f.Name() + ".html")
	if err != nil {
		panic(err)
	}

	return f.Name() + ".html"
}