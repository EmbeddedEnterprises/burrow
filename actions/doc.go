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
	"os"
	"path/filepath"
	"strings"

	"github.com/EmbeddedEnterprises/burrow/utils"
	"github.com/mattn/go-shellwords"
	"github.com/urfave/cli"
)

// Doc hosts the go documentation on the current machine.
func Doc(context *cli.Context) error {
	burrow.LoadConfig()

	burrow.Log(burrow.LOG_INFO, "doc", "Hosting documentation")

	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		gopath := os.Getenv("GOPATH")
		resolve, err := os.Readlink(cwd)

		if err != nil {
			resolve = cwd
		}

		pkg := strings.TrimPrefix(resolve, gopath+"/src/")

		burrow.Log(burrow.LOG_INFO, "doc", "Documentation of the current package is available under:")
		burrow.Log(burrow.LOG_INFO, "doc", "    http://localhost:6060/pkg/%s", pkg)
	}

	args := []string{}
	args = append(args, "-http", ":6060", "-links", "-index")
	user_args, err := shellwords.Parse(burrow.Config.Args.Go.Doc)
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "doc", "Failed to read user arguments from config file: %s", err)
		return nil
	}
	args = append(args, user_args...)
	args = append(args, burrow.GetSecondLevelArgs()...)
	return burrow.Exec("doc", "godoc", args...)
}
