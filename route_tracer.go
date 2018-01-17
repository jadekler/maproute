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
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func traceRoute(destinationIp string) []geo {
	fmt.Println("Tracing")

	waitForReadingToFinish := make(chan struct{})
	buf := bytes.Buffer{}
	stdoutInterceptor := interceptor{
		buf:          &buf,
		geos:         []geo{},
		noMoreOutput: false,
		done:         waitForReadingToFinish,
	}
	stderrInterceptor := bytes.Buffer{}

	cmd := exec.Command("traceroute", destinationIp)
	cmd.Stdout = &stdoutInterceptor
	cmd.Stderr = &stderrInterceptor

	go stdoutInterceptor.getIps()

	err := cmd.Run()
	if err != nil {
		fmt.Println(string(stderrInterceptor.Bytes()))
		panic(err)
	}

	stdoutInterceptor.noMoreOutput = true
	<-waitForReadingToFinish

	return stdoutInterceptor.geos
}

type interceptor struct {
	buf          *bytes.Buffer
	geos         []geo
	noMoreOutput bool
	done         chan (struct{})
}

func (i *interceptor) Write(p []byte) (n int, err error) {
	i.buf.Write(p)
	return len(p), nil
}

func (i *interceptor) getIps() {
	r := bufio.NewReader(i.buf)

	for {
		l, _, err := r.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				if i.noMoreOutput {
					i.done <- struct{}{}
					return
				}

				time.Sleep(time.Second)
				continue
			}

			panic(err)
		}

		line := string(l)
		ips := extractIps(line)

		if len(ips) == 0 {
			fmt.Println("* * *")
		} else {
			ip := ips[0]
			geo := getGeoForIp(ip)
			fmt.Printf("%s %s, %s (%v, %v)\n", ip, geo.City, geo.Country, geo.Lat, geo.Lng)

			i.geos = append(i.geos, geo)
		}
	}
}

func extractIps(traceRouteOut string) []string {
	res := ipRegex.FindAllString(traceRouteOut, -1)
	for i, _ := range res {
		res[i] = strings.Replace(res[i], "(", "", -1) // hack because i don't feel like figuring out regex heh
		res[i] = strings.Replace(res[i], ")", "", -1)
	}

	return res
}
