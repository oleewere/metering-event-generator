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
package main

import (
	"fmt"
	"os"

	"github.com/oleewere/meteringp/producer"
	"github.com/urfave/cli"
)

// Version that will be generated during the build as a constant
var Version string

// GitRevString that will be generated during the build as a constant - represents git revision value
var GitRevString string

func main() {
	app := cli.NewApp()
	app.Name = "meteringp"
	app.Usage = "CLI tool for generate metering events"
	app.EnableBashCompletion = true
	app.UsageText = "meteringp [command options] [arguments...]"
	if len(Version) > 0 {
		app.Version = Version
	} else {
		app.Version = "0.1.0"
	}
	if len(GitRevString) > 0 {
		app.Version = app.Version + fmt.Sprintf(" (git short hash: %v)", GitRevString)
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config, c", Usage: "Producer configuration file location"},
	}
	app.Email = "oleewere@gmail.com"
	app.Author = "Oliver Mihaly Szabo"
	app.Copyright = "Copyright 2019 Oliver Mihaly Szabo"
	app.Action = func(c *cli.Context) error {
		configLocation := c.String("config")
		if len(configLocation) == 0 {
			fmt.Println("Parameter '--config' is required")
			os.Exit(1)
		}
		producer, err := producer.ReadProducerFromConfig(configLocation)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		producer.Run()

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
