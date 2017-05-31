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

func Doc(context *cli.Context) error {
	burrow.LoadConfig()

	outputs := []string{}

	if burrow.IsTargetUpToDate("doc", outputs) {
		burrow.Log(burrow.LOG_INFO, "doc", "Documentation is up-to-date")
		return nil
	}

	burrow.Log(burrow.LOG_INFO, "doc", "Building documentation")

	args := []string{}
	args = append(args, "doc")
	user_args, err := shellwords.Parse(burrow.Config.Args.Go.Doc)
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "doc", "Failed to read user arguments from config file: %s", err)
		return nil
	}
	args = append(args, user_args...)
	err = burrow.Exec("go", args...)
	if err == nil {
		burrow.UpdateTarget("doc", outputs)
	}
	return err
}
