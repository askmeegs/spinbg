// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var counter struct {
	sync.RWMutex
	m map[string]int
}
var wg sync.WaitGroup

func MakeRequest(url string, i int) {
	defer wg.Done()
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	counter.Lock()

	b := string(body)
	b = strings.TrimSpace(b)
	if _, ok := counter.m[b]; !ok {
		counter.m[b] = 1
	} else {
		counter.m[b] += 1
	}
	counter.Unlock()
	return
}

func main() {
	total := 100

	counter = struct {
		sync.RWMutex
		m map[string]int
	}{m: make(map[string]int)}

	url := "http://35.238.38.69/"

	wg.Add(total)

	for i := 0; i < total; i++ {
		go MakeRequest(url, i)
	}
	wg.Wait()

	counter.RLock()

	fmt.Println("-------------------------------------------")
	for k, v := range counter.m {
		fmt.Printf("%s\t\t %d%%\n", k, v)
	}
	fmt.Println("-------------------------------------------")

	counter.RUnlock()

}
