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

	fileLogger := MeteringEventFileLogger{}
	fileLoggerSettings := cfg.Section("log_file")
	fileLoggerEnabled, _ := fileLoggerSettings.Key("enabled").Bool()
	if fileLoggerEnabled {
		fileLogger.Enabled = true
		fileLogger.LogFile = fileLoggerSettings.Key("file").String()
		fileLogger.MaxAge, _ = fileLoggerSettings.Key("maxAge").Int()
		fileLogger.MaxBackups, _ = fileLoggerSettings.Key("maxBackups").Int()
		fileLogger.MaxSizeMB, _ = fileLoggerSettings.Key("maxSizeMB").Int()
		fileLogger.Compress, _ = fileLoggerSettings.Key("compress").Bool()
	}

	commandOutputFields := make(map[string]MeteringCommandDetails)

	commandOutputFieldsSection := cfg.Section("command_output_fields:text")
	commandOutputFieldsKeyValues := commandOutputFieldsSection.Keys()
	for _, commandOutputField := range commandOutputFieldsKeyValues {
		commandOutputFields[commandOutputField.Name()] = MeteringCommandDetails{Command: commandOutputField.Value()}
	}

	commandOutputFieldsJSONSection := cfg.Section("command_output_fields:json")
	commandOutputFieldsJSONKeyValues := commandOutputFieldsJSONSection.Keys()
	for _, commandOutputField := range commandOutputFieldsJSONKeyValues {
		commandOutputFields[commandOutputField.Name()] = MeteringCommandDetails{Command: commandOutputField.Value(), JSON: true}
	}

	fieldsSettings := cfg.Section("fields")
	fieldKeyValues := fieldsSettings.Keys()
	fields := make(logrus.Fields)
	for _, field := range fieldKeyValues {
		fields[field.Name()] = field.Value()
	}

	embeddedJSONFieldFiles := cfg.Section("embedded_json_fields")
	embeddedJSONFieldKeys := embeddedJSONFieldFiles.Keys()
	for _, field := range embeddedJSONFieldKeys {
		fileName := field.Value()
		jsonFile, err := os.Open(fileName)
		defer jsonFile.Close()
		if err != nil {
			fmt.Println(err)
			jsonFile.Close()
			os.Exit(1)
		}
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var result interface{}
		json.Unmarshal([]byte(byteValue), &result)
		fields[field.Name()] = result
	}

	return MeteringEventProducer{FileLogger: &fileLogger, EventIDField: eventIDField, EventInerval: eventInterval,
		TimestampField: timestampField, Fields: fields, FieldCommandPairs: commandOutputFields}, nil
}
