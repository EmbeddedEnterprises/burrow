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

// The ProjectType describes whether this is a binary or a library project.
type ProjectType uint8

// The TYPE_... constants describe the type of a project.
const (
	TYPE_BIN ProjectType = iota
	TYPE_LIB
)

// Create creates a new burrow project.
func Create(context *cli.Context) error {
	if _, err := os.Stat("burrow.yaml"); err == nil {
		fmt.Println("Already a burrow project!")
		return cli.NewExitError("", burrow.EXIT_ACTION)
	}

	var err error
	projectType := TYPE_BIN
	projectName := "project"
	projectLicense := "MIT"
	projectDescription := "Burrow project"
	projectAuthors := []string{}

	for {
		fmt.Print("Is your project a binary (bin) or a library (lib)? ")
		reader := bufio.NewReader(os.Stdin)
		projectTypeStr, err := reader.ReadString('\n')
		projectTypeStr = projectTypeStr[:len(projectTypeStr)-1]
		if err == nil {
			if projectTypeStr == "bin" {
				projectType = TYPE_BIN
				break
			} else if projectTypeStr == "lib" {
				projectType = TYPE_LIB
				break
			}
		}
	}

	for {
		fmt.Print("What is the name of your project? ")
		reader := bufio.NewReader(os.Stdin)
		projectName, err = reader.ReadString('\n')
		projectName = projectName[:len(projectName)-1]
		if err == nil && projectName != "" {
			break
		}
	}

	for {
		fmt.Print("Which license (SPDX License or none) should your project use? ")
		reader := bufio.NewReader(os.Stdin)
		projectLicense, err = reader.ReadString('\n')
		projectLicense = projectLicense[:len(projectLicense)-1]

		if projectLicense == "none" {
			break
		} else if err == nil && projectLicense != "" {
			resp, err := http.Get(
				"https://raw.githubusercontent.com/spdx/license-list-data/master/text/" + projectLicense + ".txt",
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
		projectDescription, err = reader.ReadString('\n')
		projectDescription = projectDescription[:len(projectDescription)-1]
		if err == nil && projectDescription != "" {
			break
		}
	}

	for {
		fmt.Println("Please enter a comma-separated list of the authors of this project:")
		reader := bufio.NewReader(os.Stdin)
		project_authors_str, err := reader.ReadString('\n')
		project_authors_str = project_authors_str[:len(project_authors_str)-1]
		if err == nil && project_authors_str != "" {
			projectAuthors = strings.Split(project_authors_str, ",")
			break
		}
	}

	config := burrow.Configuration{}
	config.Name = projectName
	config.Version = "0.1.0"
	config.Description = projectDescription
	config.Authors = projectAuthors
	config.License = projectLicense
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
	config.Args.Git.Clone = ""
	ser, _ := yaml.Marshal(&config)

	_ = os.Mkdir("example", 0755)
	_ = ioutil.WriteFile("burrow.yaml", []byte(ser), 0644)
	_ = ioutil.WriteFile("README.md", []byte(fmt.Sprintf(readme, projectName, projectDescription)), 0644)
	_ = ioutil.WriteFile(".gitignore", []byte(gitignore), 0644)

	switch projectType {
	case TYPE_BIN:
		_ = ioutil.WriteFile("main.go", []byte(main), 0644)
	case TYPE_LIB:
		_ = ioutil.WriteFile("lib.go", []byte(fmt.Sprintf(lib, projectName)), 0644)
	}

	burrow.Exec("", "glide", "init")

	return nil
}
