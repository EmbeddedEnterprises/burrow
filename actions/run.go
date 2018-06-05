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

// Run builds and runs the burrow application.
func Run(context *cli.Context, useSecondLevelArgs bool) error {
	burrow.LoadConfig()

	if err := Build(context, false); err != nil {
		return err
	}

	example := context.String("example")

	userArgs, err := shellwords.Parse(burrow.Config.Args.Run)
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "run", "Failed to read user arguments from config file: %s", err)
		return err
	}

	args := []string{}
	args = append(args, userArgs...)

	if useSecondLevelArgs {
		args = append(args, burrow.GetSecondLevelArgs()...)
	}

	if example == "" {
		burrow.Log(burrow.LOG_INFO, "run", "Running project")
		return burrow.Exec("", "./bin/"+burrow.Config.Name, args...)
	}

	burrow.Log(burrow.LOG_INFO, "run", "Running example %s", example)
	return burrow.Exec("", "./bin/example/"+example, args...)
}
