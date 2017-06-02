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

package burrow

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/EmbeddedEnterprises/burrow/utils"
	"github.com/go-yaml/yaml"
	"github.com/urfave/cli"
)

const main = `package main

func main() {
}
`

const lib = `package %s
`

const readme = `# %s

%s

---
`

const gitignore = `*~
bin
package
vendor
`

type ProjectType uint8

const (
	TYPE_BIN ProjectType = iota
	TYPE_LIB
)

func Create(context *cli.Context) error {
	if _, err := os.Stat("burrow.yaml"); err == nil {
		fmt.Println("Already a burrow project!")
		return cli.NewExitError("", burrow.EXIT_ACTION)
	}

	var err error
	project_type := TYPE_BIN
	project_name := "project"
	project_license := "MIT"
	project_description := "Burrow project"
	project_authors := []string{}

	for {
		fmt.Print("Is your project a binary (bin) or a library (lib)? ")
		reader := bufio.NewReader(os.Stdin)
		project_type_str, err := reader.ReadString('\n')
		project_type_str = project_type_str[:len(project_type_str)-1]
		if err == nil {
			if project_type_str == "bin" {
				project_type = TYPE_BIN
				break
			} else if project_type_str == "lib" {
				project_type = TYPE_LIB
				break
			}
		}
	}

	for {
		fmt.Print("What is the name of your project? ")
		reader := bufio.NewReader(os.Stdin)
		project_name, err = reader.ReadString('\n')
		project_name = project_name[:len(project_name)-1]
		if err == nil && project_name != "" {
			break
		}
	}

	for {
		fmt.Print("Which license (SPDX License or none) should your project use? ")
		reader := bufio.NewReader(os.Stdin)
		project_license, err = reader.ReadString('\n')
		project_license = project_license[:len(project_license)-1]

		if project_license == "none" {
			break
		} else if err == nil && project_license != "" {
			resp, err := http.Get(
				"https://raw.githubusercontent.com/spdx/license-list-data/master/text/" + project_license + ".txt",
			)
			if err != nil {
				continue
			}
			defer resp.Body.Close()

			license_bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			license := string(license_bytes)
			if license == "404: Not Found\n" {
				continue
			}

			err = ioutil.WriteFile("LICENSE", license_bytes, 0644)
			if err != nil {
				continue
			}

			break
		}
	}

	for {
		fmt.Println("Please enter a description of your project:")
		reader := bufio.NewReader(os.Stdin)
		project_description, err = reader.ReadString('\n')
		project_description = project_description[:len(project_description)-1]
		if err == nil && project_description != "" {
			break
		}
	}

	for {
		fmt.Println("Please enter a comma-separated list of the authors of this project:")
		reader := bufio.NewReader(os.Stdin)
		project_authors_str, err := reader.ReadString('\n')
		project_authors_str = project_authors_str[:len(project_authors_str)-1]
		if err == nil && project_authors_str != "" {
			project_authors = strings.Split(project_authors_str, ",")
			break
		}
	}

	config := burrow.Configuration{}
	config.Name = project_name
	config.Version = "0.1.0"
	config.Description = project_description
	config.Authors = project_authors
	config.License = project_license
	config.Package.Include = []string{}
	config.Args.Run = ""
	config.Args.Go.Test = ""
	config.Args.Go.Build = ""
	config.Args.Go.Doc = ""
	config.Args.Go.Vet = ""
	config.Args.Go.Fmt = ""
	config.Args.Glide.Install = ""
	config.Args.Glide.Update = ""
	config.Args.Glide.Get = ""
	config.Args.Git.Tag = "-s -m 'Update version'"
	ser, _ := yaml.Marshal(&config)

	_ = os.Mkdir("example", 0755)
	_ = ioutil.WriteFile("burrow.yaml", []byte(ser), 0644)
	_ = ioutil.WriteFile("README.md", []byte(fmt.Sprintf(readme, project_name, project_description)), 0644)
	_ = ioutil.WriteFile(".gitignore", []byte(gitignore), 0644)

	switch project_type {
	case TYPE_BIN:
		_ = ioutil.WriteFile("main.go", []byte(main), 0644)
	case TYPE_LIB:
		_ = ioutil.WriteFile("lib.go", []byte(fmt.Sprintf(lib, project_name)), 0644)
	}

	burrow.Exec("", "glide", "init")

	return nil
}
