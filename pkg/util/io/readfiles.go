// Copyright 2019-present Open Networking Foundation.
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

package io

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

// isURL returns whether the given string is a URL
func isURL(str string) bool {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	u, err := url.Parse(str)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

// downloadURL downloads the given URL to a []byte
func downloadURL(url string) ([]byte, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// GetData gets data from a file or URL
func GetData(pathOrURL string) (string, []byte, error) {
	var name string
	var data []byte
	if isURL(pathOrURL) {
		bytes, err := downloadURL(pathOrURL)
		if err != nil {
			return "", nil, err
		}
		u, err := url.Parse(pathOrURL)
		if err != nil {
			return "", nil, err
		}
		name = path.Base(u.Path)
		name = name[:len(name)-len(path.Ext(name))]
		data = bytes
	} else {
		file, err := os.Open(pathOrURL)
		if err != nil {
			return "", nil, err
		}
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return "", nil, err
		}
		name = path.Base(pathOrURL)
		name = name[:len(name)-len(path.Ext(name))]
		data = bytes
	}
	return name, data, nil
}
