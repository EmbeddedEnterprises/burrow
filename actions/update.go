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
	"github.com/EmbeddedEnterprises/burrow/utils"
	"github.com/mattn/go-shellwords"
	"github.com/urfave/cli"
)

// Update gets all dependencies from the go.mod file and updates the go.sum file.
func Update(context *cli.Context, useSecondLevelArgs bool) error {
	burrow.LoadConfig()
	burrow.Log(burrow.LOG_INFO, "update", "Updating dependencies from go.mod config")
	args := []string{}
	args = append(args, "get", "-u")
	userArgs, err := shellwords.Parse(burrow.Config.Args.Go.Get)
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "update", "Failed to read user arguments from config file: %s", err)
		return err
	}
	args = append(args, userArgs...)

	if useSecondLevelArgs {
		args = append(args, burrow.GetSecondLevelArgs()...)
	}

	err = burrow.Exec("update", "go", args...)

	burrow.Deprecation("update", append([]string{"go"}, args...))

	return err
}
