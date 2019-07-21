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

import log "github.com/sirupsen/logrus"

// MeteringEventProducer metering producer type which holds required configuration
type MeteringEventProducer struct {
	EventInerval      int
	Fields            log.Fields
	EventIDField      string
	TimestampField    string
	IDGeneratorFields []string
	FieldCommandPairs map[string]string
	FileLogger        *MeteringEventFileLogger
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
