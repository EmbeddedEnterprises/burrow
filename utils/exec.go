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
	"os/exec"

	"github.com/urfave/cli"
)

func Exec(target string, comm string, args ...string) error {
	cmd := exec.Command(comm, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = NewLogger(target, LOG_INFO)
	cmd.Stderr = NewLogger(target, LOG_WARN)

	if err := cmd.Run(); err != nil {
		Log(LOG_ERR, target, "Error running action: %v", err)
		return cli.NewExitError("", EXIT_ACTION)
	}
	return nil
}