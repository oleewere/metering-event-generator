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
	"os"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// MeteringJSONFormatter type for custom json formatter
type MeteringJSONFormatter struct {
}

// Format json formatter for logrus entries
func (f *MeteringJSONFormatter) Format(entry *log.Entry) ([]byte, error) {
	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}

// MeteringEventProducer metering producer type which holds required configuration
type MeteringEventProducer struct {
	UseLogFile        bool
	LogFile           string
	EventInerval      int
	Fields            log.Fields
	EventIDField      string
	TimestampField    string
	IDGeneratorFields []string
	FieldCommandPairs map[string]string
}

// Run start metering event producer
func (p *MeteringEventProducer) Run() {
	log.SetFormatter(new(MeteringJSONFormatter))
	if p.UseLogFile {
		logFile, err := os.OpenFile(p.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			os.Exit(1)
		}
		log.SetOutput(logFile)
	} else {
		log.SetOutput(os.Stdout)
	}
	for {
		fields := p.Fields
		fields[p.EventIDField] = uuid.NewV4()
		fields[p.TimestampField] = time.Now().Unix()
		if len(p.FieldCommandPairs) > 0 {
			for field, command := range p.FieldCommandPairs {
				splitted := strings.Split(command, " ")
				var output string
				var err error
				if len(splitted) == 1 {
					output, _, err = RunLocalCommand(splitted[0])
				} else {
					output, _, err = RunLocalCommand(splitted[0], splitted[1:]...)
				}
				if err == nil {
					fields[field] = output
				}
			}
		}
		log.WithFields(fields).Info()
		duration := int64(p.EventInerval) * int64(time.Second)
		time.Sleep(time.Duration(duration))
	}
}
