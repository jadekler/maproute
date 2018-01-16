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
