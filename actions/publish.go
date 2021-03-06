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

// Publish builds the application, packages the application and creates a new version tag in git.
func Publish(context *cli.Context, useSecondLevelArgs bool) error {
	burrow.LoadConfig()
	if err := Package(context); err != nil {
		return err
	}
	burrow.Log(burrow.LOG_INFO, "publish", "Publishing new version tag in git")

	err := burrow.Exec("publish", "git", "diff-index", "--quiet", "HEAD", "--")
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "publish", "You have unstaged changes, commit them to proceed!")
		return cli.NewExitError("", burrow.EXIT_ACTION)
	}

	args := []string{}
	args = append(args, "tag", "-f")
	userArgs, err := shellwords.Parse(burrow.Config.Args.Git.Tag)
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "publish", "Failed to read user arguments from config file: %s", err)
		return err
	}
	args = append(args, userArgs...)

	if useSecondLevelArgs {
		args = append(args, burrow.GetSecondLevelArgs()...)
	}

	args = append(args, "v"+burrow.Config.Version)
	err = burrow.Exec("publish", "git", args...)

	burrow.Deprecation("publish", append([]string{"git"}, args...))

	return err
}
