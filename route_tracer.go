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
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func traceRoute(destinationIp string) []string {
	fmt.Println("Tracing")

	stdoutInterceptor := interceptor{output: []byte{}}
	stderrInterceptor := bytes.Buffer{}

	cmd := exec.Command("traceroute", destinationIp)
	cmd.Stdout = &stdoutInterceptor
	cmd.Stderr = &stderrInterceptor

	err := cmd.Run()
	if err != nil {
		fmt.Println(string(stderrInterceptor.Bytes()))
		panic(err)
	}

	ips := extractIps(string(stdoutInterceptor.output))

	return ips
}

type interceptor struct {
	output []byte
}

func (i *interceptor) Write(p []byte) (n int, err error) {
	i.output = append(i.output, p...)

	return os.Stdout.Write(p)
}
