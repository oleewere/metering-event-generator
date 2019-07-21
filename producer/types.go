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
	"github.com/oleewere/go-buffered-processor/processor"
	log "github.com/sirupsen/logrus"
)

// MeteringEventProducer metering producer type which holds required configuration
type MeteringEventProducer struct {
	EventInerval      int
	Fields            log.Fields
	EventIDField      string
	TimestampField    string
	IDGeneratorFields []string
	FieldCommandPairs map[string]MeteringCommandDetails
	FileLogger        *MeteringEventFileLogger
	BufferedProcessor *MeteringEventBufferedProcessor
}

// MeteringEventFileLogger holds file logger details
type MeteringEventFileLogger struct {
	Enabled    bool
	LogFile    string
	MaxSizeMB  int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// MeteringCommandDetails holds command details that is used to gather a specific field
type MeteringCommandDetails struct {
	Command    string
	JSONFormat bool
}

// MeteringEventBufferedProcessor holds buffer and data processor for publishing events
type MeteringEventBufferedProcessor struct {
	ProcessorCommand string
	BatchContext     *processor.BatchContext
}
