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
	"bytes"
	"os/exec"
	"strings"
)

// RunLocalCommand run local system command
func RunLocalCommand(command string, arg ...string) (string, string, error) {
	outStr, errStr := "", ""
	cmd := exec.Command(command, arg...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr = string(stdout.Bytes()), string(stderr.Bytes())
	outStr = strings.TrimRightFunc(outStr, func(c rune) bool {
		return c == '\r' || c == '\n'
	})
	return outStr, errStr, err
}
