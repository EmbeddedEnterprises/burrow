// -*- mode: go; tab-width: 4; -*-

/* burrow - a go build system that uses glide for dependency management.
 *
 * Copyright (C) 2017  EmbeddedEnterprises
 *     Fin Christensen <christensen.fin@gmail.com>,
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

// This package contains utilities used for the burrow build system.
package burrow

import (
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

// The Configuration struct describes the layout of the burrow.yaml.
type Configuration struct {
	Name        string
	Version     string
	Description string
	Authors     []string
	License     string
	Package     struct {
		Include []string
	}
	Args struct {
		Run string
		Go  struct {
			Test  string
			Build string
			Doc   string
			Vet   string
			Fmt   string
		}
		Glide struct {
			Install string
			Update  string
			Get     string
		}
		Git struct {
			Tag   string
			Clone string
		}
	}
}

// Config is the global instance of the Configuration struct and contains the parsed data of the
// burrow.yaml.
var Config Configuration = Configuration{}

// The isConfigLoaded flag describes whether the burrow.yaml has already been parsed into the
// global Config variable.
var isConfigLoaded bool = false

// LoadConfig loads the content of the burrow.yaml into the global Config variable and sets the
// isConfigLoaded flag on success.
func LoadConfig() {
	if isConfigLoaded {
		return
	}

	data, err := ioutil.ReadFile("burrow.yaml")

	if err != nil {
		Log(LOG_ERR, "burrow", "Not a burrow project!")
		os.Exit(EXIT_CONFIG)
	}

	err = yaml.Unmarshal(data, &Config)

	if err != nil {
		Log(LOG_ERR, "burrow", "Failed to read burrow config: %v", err)
		os.Exit(EXIT_CONFIG)
	}

	isConfigLoaded = true
}
