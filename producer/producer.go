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

	lumberjack "github.com/natefinch/lumberjack"
	"github.com/oleewere/go-buffered-processor/processor"
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

// Process run local process commands on gathered metering events
func (p *MeteringEventProducer) Process(batchContext *processor.BatchContext) error {
	splitted := strings.Split(p.BufferedProcessor.ProcessorCommand, " ")
	var err error
	if len(splitted) == 1 {
		_, _, err = RunLocalCommand(splitted[0])
	} else {
		_, _, err = RunLocalCommand(splitted[0], splitted[1:]...)
	}
	return err
}

// HandleError handle errors during time based buffer processing (it is not used by this generator)
func (p *MeteringEventProducer) HandleError(batchContext *processor.BatchContext, err error) {
}

// Run start metering event producer
func (p *MeteringEventProducer) Run() {
	log.SetFormatter(new(MeteringJSONFormatter))
	fileLogger := p.FileLogger
	if fileLogger != nil && fileLogger.Enabled {
		lumberjackLogger := &lumberjack.Logger{
			Filename:   fileLogger.LogFile,
			MaxSize:    fileLogger.MaxSizeMB,
			MaxBackups: fileLogger.MaxBackups,
			MaxAge:     fileLogger.MaxAge,
			Compress:   fileLogger.Compress,
		}
		log.SetOutput(lumberjackLogger)
	} else {
		log.SetOutput(os.Stdout)
	}
	for {
		fields := p.Fields
		fields[p.EventIDField] = uuid.NewV4()
		fields[p.TimestampField] = time.Now().Unix()
		if len(p.FieldCommandPairs) > 0 {
			for field, commandDetails := range p.FieldCommandPairs {
				splitted := strings.Split(commandDetails.Command, " ")
				var output string
				var err error
				if len(splitted) == 1 {
					output, _, err = RunLocalCommand(splitted[0])
				} else {
					output, _, err = RunLocalCommand(splitted[0], splitted[1:]...)
				}
				if err == nil {
					if commandDetails.JSON {
						var jsonResult interface{}
						json.Unmarshal([]byte(output), &jsonResult)
						fields[field] = jsonResult
					} else {
						fields[field] = output
					}
				}
			}
		}
		log.WithFields(fields).Info()
		if p.BufferedProcessor != nil {
			batchContext := p.BufferedProcessor.BatchContext
			processor.ProcessData(fields, batchContext, p)
		}
		duration := int64(p.EventInerval) * int64(time.Second)
		time.Sleep(time.Duration(duration))
	}
}
