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
	eventIDField := globalSettings.Key("eventIdField").String()
	eventInterval, _ := globalSettings.Key("eventInterval").Int()

	fieldsSettings := cfg.Section("fields")
	fieldKeyValues := fieldsSettings.Keys()
	fields := make(logrus.Fields)
	for _, field := range fieldKeyValues {
		fields[field.Name()] = field.Value()
	}

	return MeteringEventProducer{EventIDField: eventIDField, EventInerval: eventInterval, Fields: fields}, nil
}
