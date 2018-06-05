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
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/EmbeddedEnterprises/burrow/utils"
	"github.com/urfave/cli"
)

func askProjectPath() string {
	for {
		fmt.Print("Please specify a base location for your project (e.g. github.com/myuser): ")
		reader := bufio.NewReader(os.Stdin)
		projectPath, err := reader.ReadString('\n')
		projectPath = projectPath[:len(projectPath)-1]
		if err == nil && projectPath != "" {
			return projectPath
		}
	}
}

// New creates a new burrow project in the GOPATH.
func New(context *cli.Context) error {
	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "new", "Failed to get current directory!")
		return err
	}

	gopath, err := filepath.Abs(os.Getenv("GOPATH"))
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "new", "Failed to get GOPATH!")
		return err
	}

	destination := "."
	if !strings.HasPrefix(cwd, gopath) {
		projectPath := askProjectPath()
		destination = path.Join(gopath, "src", projectPath)
	}

	project := NewProject()
	project.Location = path.Join(destination, project.Name)
	err = project.Dump()
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "new", "Failed to write project to destination!")
		return err
	}

	err = os.Symlink(project.Location, project.Name)
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "new", "Failed to create symlink!")
		return err
	}

	return nil
}
