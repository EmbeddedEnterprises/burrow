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
	"path"
	"path/filepath"
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

func askProjectType() ProjectType {
	for {
		fmt.Print("Is your project a binary (bin) or a library (lib)? ")
		reader := bufio.NewReader(os.Stdin)
		projectTypeStr, err := reader.ReadString('\n')
		projectTypeStr = projectTypeStr[:len(projectTypeStr)-1]
		if err == nil {
			if projectTypeStr == "bin" {
				return TYPE_BIN
			} else if projectTypeStr == "lib" {
				return TYPE_LIB
			}
		}
	}
}

func askProjectName(projectType ProjectType) string {
	for {
		fmt.Print("What is the name of your project? \n")

		isLib := projectType == TYPE_LIB

		if isLib {
			fmt.Print("It should be short and clear and should not contain dashes. \n")
			fmt.Print("More Information: https://blog.golang.org/package-names \n")
		}

		reader := bufio.NewReader(os.Stdin)
		projectName, err := reader.ReadString('\n')
		projectName = projectName[:len(projectName)-1]

		// It should not be possible to create libs, which contains dash in their name.
		// This will result in an invalid package name and glide will go nuts.

		isDasherized := strings.Contains(projectName, "-")

		if err == nil && projectName != "" && (!isLib || !isDasherized) {
			return projectName
		}
	}
}

func askProjectLicense() string {
	var projectLicense string
	for {
		fmt.Print("Which license (SPDX License or none) should your project use? ")
		reader := bufio.NewReader(os.Stdin)
		projectLicense, err := reader.ReadString('\n')
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

			licenseBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			license := string(licenseBytes)
			if license == "404: Not Found\n" {
				continue
			}

			err = ioutil.WriteFile("LICENSE", licenseBytes, 0644)
			if err != nil {
				continue
			}

			break
		}
	}

	return projectLicense
}

func askProjectDescription() string {
	for {
		fmt.Println("Please enter a description of your project:")
		reader := bufio.NewReader(os.Stdin)
		projectDescription, err := reader.ReadString('\n')
		projectDescription = projectDescription[:len(projectDescription)-1]
		if err == nil && projectDescription != "" {
			return projectDescription
		}
	}
}

func askProjectAuthors() []string {
	for {
		fmt.Println("Please enter a comma-separated list of the authors of this project:")
		reader := bufio.NewReader(os.Stdin)
		projectAuthorsStr, err := reader.ReadString('\n')
		projectAuthorsStr = projectAuthorsStr[:len(projectAuthorsStr)-1]
		if err == nil && projectAuthorsStr != "" {
			return strings.Split(projectAuthorsStr, ",")
		}
	}
}

// Project describes all inputs that are needed to create a new project.
type Project struct {
	Location    string
	Type        ProjectType
	Name        string
	License     string
	Description string
	Authors     []string
}

// NewProject creates a new project from user input from stdin.
func NewProject() *Project {
	projectType := askProjectType()
	projectName := askProjectName(projectType)
	projectLicense := askProjectLicense()
	projectDescription := askProjectDescription()
	projectAuthors := askProjectAuthors()

	return &Project{
		Location:    path.Dir("."),
		Type:        projectType,
		Name:        projectName,
		License:     projectLicense,
		Description: projectDescription,
		Authors:     projectAuthors,
	}
}

// Dump writes all files of a project to the specified location (Project.Location).
func (p *Project) Dump() error {
	config := burrow.Configuration{}
	config.Name = p.Name
	config.Version = "0.1.0"
	config.Description = p.Description
	config.Authors = p.Authors
	config.License = p.License
	config.Package.Include = []string{}
	config.Args.Run = ""
	config.Args.Go.Test = ""
	config.Args.Go.Build = ""
	config.Args.Go.Doc = ""
	config.Args.Go.Vet = ""
	config.Args.Go.Fmt = "-s"
	config.Args.Glide.Install = ""
	config.Args.Glide.Update = ""
	config.Args.Glide.Get = ""
	config.Args.Git.Tag = "-s -m 'Update version'"
	config.Args.Git.Clone = ""
	ser, err := yaml.Marshal(&config)
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "project", "Failed to serialize config file: %s", err)
		return err
	}

	if _, err = os.Stat(p.Location); os.IsNotExist(err) {
		if err = os.MkdirAll(p.Location, 0755); err != nil {
			burrow.Log(burrow.LOG_ERR, "project", "Failed to create project directory!")
			return err
		}
	}

	if err = os.MkdirAll(path.Join(p.Location, "example"), 0755); err != nil {
		burrow.Log(burrow.LOG_ERR, "project", "Failed to create example directory!")
		return err
	}
	if err = ioutil.WriteFile(path.Join(p.Location, "burrow.yaml"), []byte(ser), 0644); err != nil {
		burrow.Log(burrow.LOG_ERR, "project", "Failed to write configuration!")
		return err
	}
	if err = ioutil.WriteFile(
		path.Join(p.Location, "README.md"),
		[]byte(fmt.Sprintf(readme, p.Name, p.Description)),
		0644,
	); err != nil {
		burrow.Log(burrow.LOG_ERR, "project", "Failed to write README!")
		return err
	}
	if err = ioutil.WriteFile(path.Join(p.Location, ".gitignore"), []byte(gitignore), 0644); err != nil {
		burrow.Log(burrow.LOG_ERR, "project", "Failed to write gitignore!")
		return err
	}

	switch p.Type {
	case TYPE_BIN:
		if err = ioutil.WriteFile(path.Join(p.Location, "main.go"), []byte(main), 0644); err != nil {
			burrow.Log(burrow.LOG_ERR, "project", "Failed to write main.go!")
			return err
		}
	case TYPE_LIB:
		if err = ioutil.WriteFile(
			path.Join(p.Location, "lib.go"),
			[]byte(fmt.Sprintf(lib, p.Name)),
			0644,
		); err != nil {
			burrow.Log(burrow.LOG_ERR, "project", "Failed to write lib.go!")
			return err
		}
	}

	return burrow.ExecDir("", p.Location, "glide", "init")
}

// Create initializes a directory as a burrow project.
func Create(context *cli.Context) error {
	if _, err := os.Stat("burrow.yaml"); err == nil {
		fmt.Println("Already a burrow project!")
		return cli.NewExitError("", burrow.EXIT_ACTION)
	}

	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "init", "Failed to get current directory!")
		return err
	}

	gopath, err := filepath.Abs(os.Getenv("GOPATH"))
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "init", "Failed to get GOPATH!")
		return err
	}

	if !strings.HasPrefix(cwd, gopath) {
		burrow.Log(burrow.LOG_WARN, "init", "Initializing project outside of GOPATH!")
	}

	project := NewProject()
	return project.Dump()
}
