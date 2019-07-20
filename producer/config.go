// Copyright 2019 Oliver Szabo
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
// limitations under the License

package producer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

// ReadProducerFromConfig create producer struct from producer configuration
func ReadProducerFromConfig(configFile string) (MeteringEventProducer, error) {
	cfg, err := ini.Load(configFile)
	if err != nil {
		return MeteringEventProducer{}, err
	}
	globalSettings := cfg.Section("global")
	timestampField := globalSettings.Key("timestampField").String()
	eventIDField := globalSettings.Key("eventIdField").String()
	eventInterval, _ := globalSettings.Key("eventInterval").Int()
	useLogFile, _ := globalSettings.Key("useLogFile").Bool()
	logFile := ""
	if useLogFile {
		logFile = globalSettings.Key("logFile").String()
	}

	fieldsSettings := cfg.Section("fields")
	fieldKeyValues := fieldsSettings.Keys()
	fields := make(logrus.Fields)
	for _, field := range fieldKeyValues {
		fields[field.Name()] = field.Value()
	}

	embeddedJSONFieldMapFiles := cfg.Section("embedded_json_fields:map")
	embeddedJSONFieldMapKeys := embeddedJSONFieldMapFiles.Keys()
	for _, field := range embeddedJSONFieldMapKeys {
		fileName := field.Value()
		jsonFile, err := os.Open(fileName)
		defer jsonFile.Close()
		if err != nil {
			fmt.Println(err)
			jsonFile.Close()
			os.Exit(1)
		}
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var result map[string]interface{}
		json.Unmarshal([]byte(byteValue), &result)
		fields[field.Name()] = result
	}

	embeddedJSONFieldArrayFiles := cfg.Section("embedded_json_fields:array")
	embeddedJSONFieldArrayKeys := embeddedJSONFieldArrayFiles.Keys()
	for _, field := range embeddedJSONFieldArrayKeys {
		fileName := field.Value()
		jsonFile, err := os.Open(fileName)
		defer jsonFile.Close()
		if err != nil {
			fmt.Println(err)
			jsonFile.Close()
			os.Exit(1)
		}
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var result []interface{}
		json.Unmarshal([]byte(byteValue), &result)
		fields[field.Name()] = result
	}

	return MeteringEventProducer{UseLogFile: useLogFile, LogFile: logFile, EventIDField: eventIDField,
		EventInerval: eventInterval, TimestampField: timestampField, Fields: fields}, nil
}
