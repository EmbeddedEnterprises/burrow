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
	"github.com/urfave/cli"
)

// Clone clones a git repository
func Clone(context *cli.Context) error {
	options := context.Args()

	if len(options) != 1 {
		cli.ShowCommandHelp(context, "clone")
		return nil
	}

	url := options[0]

	burrow.Log(burrow.LOG_INFO, "clone", "Cloning git repository...")

	args := []string{}
	args = append(args, "clone", url)
	if err := burrow.Exec("clone", "git", args...); err != nil {
		return err
	}

	burrow.Deprecation("clone", append([]string{"git"}, args...))

	return nil
}
