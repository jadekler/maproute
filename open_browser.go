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
	"os/exec"
	"fmt"
)

func openFileInBrowser(path string) {
	err := tryMacOpen(path)
	if err == nil {
		return
	} else {
		fmt.Println("Tried mac's open cli directly, no dice (no idea if we're running on a mac or not, btw)")
		fmt.Println(err)
	}

	err = tryChrome(path)
	if err == nil {
		return
	} else {
		fmt.Println("Tried chrome directly, no dice")
		fmt.Println(err)
	}

	err = tryFirefox(path)
	if err == nil {
		return
	} else {
		fmt.Println("Tried firefox directly, no dice")
		fmt.Println(err)
	}

	fmt.Println("Couldn't find a way to open an html file!")
	fmt.Printf("You can open the file directly at %s", path)
}

func tryChrome(path string) error {
	cmd := "chrome"
	args := []string{path}
	_, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return err
	}

	return nil
}

func tryFirefox(path string) error {
	cmd := "firefox"
	args := []string{path}
	_, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return err
	}

	return nil
}

func tryMacOpen(path string) error {
	cmd := "open"
	args := []string{path}
	_, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return err
	}

	return nil
}
