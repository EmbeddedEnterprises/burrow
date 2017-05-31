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

package main

import (
	"os"

	"github.com/EmbeddedEnterprises/burrow/actions"
	"github.com/urfave/cli"
)

// TODO: Add -f, --force to all cached target to disable build cache
// TODO: Add -- to every target to enable argument passing to underlying program

func main() {
	app := cli.NewApp()
	app.Name = "burrow"
	app.Usage = "build glide managed go programs"
	app.Version = "0.0.1"
	app.Action = func(context *cli.Context) error {
		return cli.ShowAppHelp(context)
	}
	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"create"},
			Flags:   []cli.Flag{},
			Usage:   "Create a new burrow project.",
			Action:  burrow.Create,
		},
		{
			Name:    "get",
			Aliases: []string{},
			Flags:   []cli.Flag{},
			Usage:   "Install a dependency in the vendor folder and add it to the glide yaml",
			Action:  burrow.Get,
		},
		{
			Name:    "fetch",
			Aliases: []string{"ensure", "f", "e"},
			Flags:   []cli.Flag{},
			Usage:   "Get all dependencies from the lock file to reproduce a build",
			Action:  burrow.Fetch,
		},
		{
			Name:    "update",
			Aliases: []string{"u", "up"},
			Flags:   []cli.Flag{},
			Usage:   "Update all dependencies from the yaml file and update the lock file",
			Action:  burrow.Update,
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Flags:   []cli.Flag{},
			Usage:   "Build and run the application",
			Action:  burrow.Run,
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Flags:   []cli.Flag{},
			Usage:   "Run all existing tests of the application",
			Action:  burrow.Test,
		},
		{
			Name:    "build",
			Aliases: []string{"b"},
			Flags:   []cli.Flag{},
			Usage:   "Build the application",
			Action:  burrow.Build,
		},
		{
			Name:    "install",
			Aliases: []string{"i", "in", "inst"},
			Flags:   []cli.Flag{},
			Usage:   "Install the application in the GOPATH",
			Action:  burrow.Install,
		},
		{
			Name:    "uninstall",
			Aliases: []string{"un", "uninst"},
			Flags:   []cli.Flag{},
			Usage:   "Uninstall the application from the GOPATH",
			Action:  burrow.Uninstall,
		},
		{
			Name:    "package",
			Aliases: []string{"pack"},
			Flags:   []cli.Flag{},
			Usage:   "Create a .tar.gz containing the binary",
			Action:  burrow.Package,
		},
		{
			Name:    "publish",
			Aliases: []string{"pub"},
			Flags:   []cli.Flag{},
			Usage:   "Publish the current version by building a package and setting a version tag in git",
			Action:  burrow.Publish,
		},
		{
			Name:    "clean",
			Aliases: []string{},
			Flags:   []cli.Flag{},
			Usage:   "Clean the project from any build artifacts",
			Action:  burrow.Clean,
		},
		{
			Name:    "doc",
			Aliases: []string{},
			Flags:   []cli.Flag{},
			Usage:   "Generate the godoc documentation for this project",
			Action:  burrow.Doc,
		},
		{
			Name:    "format",
			Aliases: []string{"fmt"},
			Flags:   []cli.Flag{},
			Usage:   "Format the code of this project with gofmt",
			Action:  burrow.Format,
		},
		{
			Name:    "check",
			Aliases: []string{},
			Flags:   []cli.Flag{},
			Usage:   "Check the code with go vet",
			Action:  burrow.Check,
		},
	}

	app.Run(os.Args)
}
