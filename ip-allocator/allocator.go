// Copyright 2016 Open Networking Laboratory
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"net/http"
)

type Config struct {
	Port      int    `default:"4242"`
	Listen    string `default:"0.0.0.0"`
	Network   string `default:"10.0.0.0/24"`
	RangeLow  string `default:"10.0.0.2" envconfig:"RANGE_LOW"`
	RangeHigh string `default:"10.0.0.253" envconfig:"RANGE_HIGH"`
	LogLevel  string `default:"warning" envconfig:"LOG_LEVEL"`
	LogFormat string `default:"text" envconfig:"LOG_FORMAT"`
}

type Context struct {
	storage Storage
}

var log = logrus.New()

func main() {
	context := &Context{}

	config := Config{}
	err := envconfig.Process("ALLOCATE", &config)
	if err != nil {
		log.Fatalf("Unable to parse configuration options : %s", err)
	}

	switch config.LogFormat {
	case "json":
		log.Formatter = &logrus.JSONFormatter{}
	default:
		log.Formatter = &logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		}
	}

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		level = logrus.WarnLevel
	}
	log.Level = level

	log.Infof(`Configuration:
	    LISTEN:       %s
	    PORT:         %d
	    NETWORK:      %s
	    RANGE_LOW:    %s
	    RANGE_HIGH:   %s
	    LOG_LEVEL:    %s
	    LOG_FORMAT:   %s`,
		config.Listen, config.Port,
		config.Network, config.RangeLow, config.RangeHigh,
		config.LogLevel, config.LogFormat)

	context.storage = &MemoryStorage{}
	context.storage.Init(config.Network, config.RangeLow, config.RangeHigh)

	router := mux.NewRouter()
	router.HandleFunc("/allocations/{mac}", context.ReleaseAllocationHandler).Methods("DELETE")
	router.HandleFunc("/allocations/{mac}", context.AllocationHandler).Methods("GET")
	router.HandleFunc("/allocations/", context.ListAllocationsHandler).Methods("GET")
	router.HandleFunc("/addresses/{ip}", context.FreeAddressHandler).Methods("DELETE")
	http.Handle("/", router)

	http.ListenAndServe(fmt.Sprintf("%s:%d", config.Listen, config.Port), nil)
}
