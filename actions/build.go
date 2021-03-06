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

// Package burrow contains all actions that can be executed as subcommands.
package burrow

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/EmbeddedEnterprises/burrow/utils"
	"github.com/mattn/go-shellwords"
	"github.com/urfave/cli"
)

// Build builds a burrow application to bin/.
func Build(context *cli.Context, useSecondLevelArgs bool) error {
	burrow.LoadConfig()

	outputs := []string{}
	sources := []string{}

	_, err := os.Stat("main.go")
	if err == nil {
		outputs = append(outputs, "./bin/"+burrow.Config.Name)
		sources = append(sources, "main.go")
	}

	_ = filepath.Walk("./example", func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") && !f.IsDir() {
			name := f.Name()
			outputs = append(outputs, "./bin/example/"+name[:len(name)-3])
			sources = append(sources, path)
		}
		return nil
	})

	if burrow.IsTargetUpToDate("build", outputs) && !context.Bool("force") {
		burrow.Log(burrow.LOG_INFO, "build", "Build is up-to-date")
		return nil
	}

	burrow.Log(burrow.LOG_INFO, "build", "Building project")

	_ = os.Mkdir("./bin", 0755)

	userArgs, err := shellwords.Parse(burrow.Config.Args.Go.Build)
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "build", "Failed to read user arguments from config file: %s", err)
		return err
	}
	buildArgs := burrow.GetSecondLevelArgs()

	deprecationArgs := make([][]string, 0)
	for i, output := range outputs {
		args := []string{}
		args = append(args, "build", "-o", output)
		args = append(args, userArgs...)

		if useSecondLevelArgs {
			args = append(args, buildArgs...)
		}

		args = append(args, sources[i])

		deprecationArgs = append(deprecationArgs, append([]string{"go"}, args...))

		if err = burrow.Exec("build", "go", args...); err != nil {
			return err
		}
	}

	if err == nil {
		burrow.UpdateTarget("build", outputs)
	}

	burrow.Deprecation("build", deprecationArgs...)

	return err
}
