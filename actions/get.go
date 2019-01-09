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

// Get installs a dependency in the vendor folder and adds it to the go.mod.
func Get(context *cli.Context, useSecondLevelArgs bool) error {
	burrow.LoadConfig()

	if len(context.Args()) != 1 {
		cli.ShowCommandHelp(context, "get")
		return nil
	}

	dep := context.Args()[0]
	burrow.Log(burrow.LOG_INFO, "get", "Adding new dependency %s", dep)

	args := []string{}
	args = append(args, "get")
	userArgs, err := shellwords.Parse(burrow.Config.Args.Go.Get)
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "get", "Failed to read user arguments from config file: %s", err)
		return err
	}
	args = append(args, userArgs...)

	if useSecondLevelArgs {
		args = append(args, burrow.GetSecondLevelArgs()...)
	}

	args = append(args, dep)

	err = burrow.Exec("get", "go", args...)

	burrow.Deprecation("get", append([]string{"go"}, args...))

	return err
}
