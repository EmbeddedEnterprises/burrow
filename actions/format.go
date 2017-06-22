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
	"github.com/EmbeddedEnterprises/burrow/utils"
	"github.com/mattn/go-shellwords"
	"github.com/urfave/cli"
)

// Format formats the code of the current burrow project with gofmt.
func Format(context *cli.Context, useSecondLevelArgs bool) error {
	burrow.LoadConfig()

	outputs := []string{}

	if burrow.IsTargetUpToDate("format", outputs) && !context.Bool("force") {
		burrow.Log(burrow.LOG_INFO, "format", "Code formatting is up-to-date")
		return nil
	}

	burrow.Log(burrow.LOG_INFO, "format", "Formatting code")

	args := []string{}
	args = append(args, "-l", "-w")
	userArgs, err := shellwords.Parse(burrow.Config.Args.Go.Fmt)
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "format", "Failed to read user arguments from config file: %s", err)
		return nil
	}
	args = append(args, userArgs...)

	if useSecondLevelArgs {
		args = append(args, burrow.GetSecondLevelArgs()...)
	}

	args = append(args, burrow.GetCodefiles()...)
	err = burrow.Exec("format", "gofmt", args...)
	if err == nil {
		burrow.UpdateTarget("format", outputs)
	}
	return err
}
