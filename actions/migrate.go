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

// Package burrow contains all actions that can be executed as subcommands.
package burrow

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/EmbeddedEnterprises/burrow/utils"
	"github.com/urfave/cli"
)

func Migrate(context *cli.Context) error {
	target := "migrate"
	gopath := os.Getenv("GOPATH")

	if gopath == "" {
		burrow.Log(
			burrow.LOG_ERR,
			target,
			"No GOPATH environment variable set. This is needed to migrate an old GOPATH",
		)
		burrow.Log(
			burrow.LOG_ERR,
			target,
			"style project to a new 'go mod' project.",
		)
		return cli.NewExitError("", burrow.EXIT_ACTION)
	}

	gopath, err := filepath.Abs(gopath)
	if err != nil {
		burrow.Log(
			burrow.LOG_ERR,
			target,
			"Cannot resolve GOPATH to an absolute path!",
		)
		return cli.NewExitError("", burrow.EXIT_ACTION)
	}

	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		burrow.Log(
			burrow.LOG_ERR,
			target,
			"Cannot resolve current working directory to an absolute path.",
		)
		return cli.NewExitError("", burrow.EXIT_ACTION)
	}

	if strings.HasPrefix(cwd, gopath) {
		burrow.Log(
			burrow.LOG_ERR,
			target,
			"Cannot migrate into GOPATH. Either provide the go uri of your project",
		)
		burrow.Log(
			burrow.LOG_ERR,
			target,
			"(my.domain/user/project) or access the burrow project via a symlink from",
		)
		burrow.Log(
			burrow.LOG_ERR,
			target,
			"outside the GOPATH.",
		)
		return cli.NewExitError("", burrow.EXIT_ACTION)
	}

	args := context.Args()
	_, err = os.Stat("burrow.yaml")
	inBurrowProject := err == nil

	var source string
	var destination string
	if inBurrowProject {
		if len(args) > 0 {
			burrow.Log(
				burrow.LOG_ERR,
				target,
				"No positional arguments allowed when burrow.yaml available in current",
			)
			burrow.Log(
				burrow.LOG_ERR,
				target,
				"directory.",
			)
			return cli.NewExitError("", burrow.EXIT_ACTION)
		}

		source, err = filepath.EvalSymlinks(cwd)
		if err != nil {
			// current directory is burrow project but does not symlink into GOPATH
			source = cwd
		}
		destination = cwd

		os.Chdir(filepath.Dir(cwd))
	} else {
		if len(args) != 1 {
			burrow.Log(
				burrow.LOG_ERR,
				target,
				"Please provide exactly one positional argument to migrate containing the",
			)
			burrow.Log(
				burrow.LOG_ERR,
				target,
				"go uri for the project you want to migrate, when your current working",
			)
			burrow.Log(
				burrow.LOG_ERR,
				target,
				"directory is not a burrow project!",
			)
			burrow.Log(
				burrow.LOG_ERR,
				target,
				"Example: burrow migrate github.com/myuser/myproject",
			)
			return cli.NewExitError("", burrow.EXIT_ACTION)
		}
		source = filepath.Join(gopath, args[0])

		if _, err := os.Stat(source); err != nil {
			burrow.Log(
				burrow.LOG_ERR,
				target,
				"Cannot find specified project in your GOPATH!",
			)
			return cli.NewExitError("", burrow.EXIT_ACTION)
		}

		destination = filepath.Join(cwd, filepath.Base(source))
	}

	if err := os.Remove(destination); inBurrowProject && err != nil {
		burrow.Log(
			burrow.LOG_ERR,
			target,
			"Failed to remove symlink of burrow project!",
		)
		return cli.NewExitError("", burrow.EXIT_ACTION)
	}

	if _, err := os.Stat("/path/to/whatever"); !os.IsNotExist(err) {
		burrow.Log(
			burrow.LOG_ERR,
			target,
			"Destination path already exists!",
		)
		return cli.NewExitError("", burrow.EXIT_ACTION)
	}

	if err := os.Rename(source, destination); err != nil {
		burrow.Log(
			burrow.LOG_ERR,
			target,
			"Failed to move go project from inside GOPATH to new destination!",
		)
		return cli.NewExitError("", burrow.EXIT_ACTION)
	}

	if inBurrowProject {
		os.Chdir(destination)
	}

	goURI, err := filepath.Rel(filepath.Join(gopath, "src"), source)
	if err != nil {
		burrow.Log(
			burrow.LOG_ERR,
			target,
			"Failed to get go uri from source path!",
		)
		return cli.NewExitError("", burrow.EXIT_ACTION)
	}

	os.RemoveAll(filepath.Join(destination, "vendor"))
	burrow.Exec(target, "go", "mod", "init", goURI)
	os.Remove(filepath.Join(destination, "glide.yaml"))
	os.Remove(filepath.Join(destination, "glide.lock"))

	return nil
}
