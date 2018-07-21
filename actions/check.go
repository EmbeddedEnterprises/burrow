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

	"github.com/EmbeddedEnterprises/burrow/utils"
	"github.com/mattn/go-shellwords"
	"github.com/urfave/cli"
)

// Check checks the code of a burrow project with go vet.
func Check(context *cli.Context, useSecondLevelArgs bool) error {
	burrow.LoadConfig()

	outputs := []string{}

	if burrow.IsTargetUpToDate("check", outputs) && !context.Bool("force") {
		burrow.Log(burrow.LOG_INFO, "check", "Code has already been checked")
		return nil
	}

	burrow.Log(burrow.LOG_INFO, "check", "Checking code")

	args := []string{}
	args = append(args, "vet")
	userArgs, err := shellwords.Parse(burrow.Config.Args.Go.Vet)
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "check", "Failed to read user arguments from config file: %s", err)
		return err
	}
	args = append(args, userArgs...)

	if useSecondLevelArgs {
		args = append(args, burrow.GetSecondLevelArgs()...)
	}
	wd, err := os.Getwd()
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "check", "Failed to get working directory: %s", err)
		return err
	}
	err = burrow.ExecDir("check", wd, "go", args...)
	if err == nil {
		burrow.UpdateTarget("check", outputs)
	}

	return err
}
