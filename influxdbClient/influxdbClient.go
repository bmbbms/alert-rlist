// Copyright 2015 Google Inc. All Rights Reserved.
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

package influxdbClient

import (
	"fmt"
	influxdb "github.com/influxdata/influxdb1-client/v2"
	"net/url"
	"time"
)

import (
	_ "github.com/influxdata/influxdb1-client"
)

func NewClient(c *influxdb.HTTPConfig) (influxdb.Client, error) {

	iConfig := &influxdb.HTTPConfig{
		Addr:     c.Addr,
		Username: c.Username,
		Password: c.Password,
	}
	client, err := influxdb.NewHTTPClient(*iConfig)

	if err != nil {
		return nil, err
	}
	if _, _, err := client.Ping(2 * time.Second); err != nil {
		return nil, fmt.Errorf("failed to ping InfluxDB server at %q - %v", c.Addr, err)
	}
	return client, nil
}

func BuildConfig(uri *url.URL) (*influxdb.HTTPConfig, error) {
	config := influxdb.HTTPConfig{
		Addr:               uri.Scheme + "://" + uri.Host,
		Username:           "admin",
		Password:           "admin#20220818",
		UserAgent:          "InfluxDBClient",
		Timeout:            5 * time.Second,
		InsecureSkipVerify: true,
		TLSConfig:          nil,
	}

	return &config, nil
}
