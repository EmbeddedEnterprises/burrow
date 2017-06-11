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

// This package contains a go build system that used glide for dependency management.
package main

import (
	"os"

	"github.com/EmbeddedEnterprises/burrow/actions"
	"github.com/urfave/cli"
)

// The main function is the entry point of burrow. This should only contain cli configuration.
func main() {
	force_flag := cli.BoolFlag{
		Name:  "force, f",
		Usage: "Forces this action to be run, even if cached data is available",
	}
	example_flag := cli.StringFlag{
		Name:  "example, e",
		Usage: "Run an example (specified by name) instead of the application itself",
	}

	// TODO: Add readme

	cli.AppHelpTemplate = `Usage: {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}

{{.Usage}}
{{if .Commands}}
Commands:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
Global options:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if len .Authors}}
Authors:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
{{.Name}} - {{.Copyright}}{{end}}
`

	app := cli.NewApp()
	app.Name = "burrow"
	app.Usage = "A go build system that uses glide for dependency management."
	app.Version = "0.0.9"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Fin Christensen",
			Email: "christensen.fin@gmail.com",
		},
	}
	app.Copyright = "Copyright (c) 2017  EmbeddedEnterprises"
	app.Action = func(context *cli.Context) error {
		return cli.ShowAppHelp(context)
	}
	app.Commands = []cli.Command{
		{
			Name:        "init",
			Aliases:     []string{"create"},
			Flags:       []cli.Flag{},
			Usage:       "Create a new burrow project.",
			Description: "This action creates a new burrow project in the current directory.",
			Action:      burrow.Create,
		},
		{
			Name:        "clone",
			Aliases:     []string{},
			Flags:       []cli.Flag{},
			Usage:       "Clone a git repository into your GOPATH and create a symbolic link in your current location.",
			Description: "This action clones a git repository (go-get url scheme) into your GOPATH and creates a symbolic link in the current directory.",
			Action:      burrow.Clone,
		},
		{
			Name:        "get",
			Aliases:     []string{},
			Flags:       []cli.Flag{},
			Usage:       "Install a dependency in the vendor folder and add it to the glide yaml.",
			Description: "This runs glide get in the current directory. The first argument should be the go-get url and any argument following -- get passed directly to glide.",
			Action:      burrow.Get,
		},
		{
			Name:        "fetch",
			Aliases:     []string{"ensure", "f", "e"},
			Flags:       []cli.Flag{},
			Usage:       "Get all dependencies from the lock file to reproduce a build.",
			Description: "This runs glide install in the current directory. Any arguments following -- get passed directly to glide.",
			Action:      burrow.Fetch,
		},
		{
			Name:        "update",
			Aliases:     []string{"u", "up"},
			Flags:       []cli.Flag{},
			Usage:       "Update all dependencies from the yaml file and update the lock file.",
			Description: "This runs glide update in the current directory. Any arguments following -- get passed directly to gilde.",
			Action:      burrow.Update,
		},
		{
			Name:        "run",
			Aliases:     []string{"r"},
			Flags:       []cli.Flag{example_flag},
			Usage:       "Build and run the application.",
			Description: "This runs the compiled binary. Any arguments following -- will be directly passed to your application.",
			Action:      burrow.Run,
		},
		{
			Name:        "test",
			Aliases:     []string{"t"},
			Flags:       []cli.Flag{force_flag},
			Usage:       "Run all existing tests of the application.",
			Description: "This runs go test in the current directory. Any arguments following -- will be directly passed to go.",
			Action:      burrow.Test,
		},
		{
			Name:        "build",
			Aliases:     []string{"b"},
			Flags:       []cli.Flag{force_flag},
			Usage:       "Build the application.",
			Description: "This runs go build in the current directory for your application and all examples. Any arguments following -- will be directly passed to go.",
			Action:      burrow.Build,
		},
		{
			Name:        "install",
			Aliases:     []string{"i", "in", "inst"},
			Flags:       []cli.Flag{force_flag},
			Usage:       "Install the application in the GOPATH.",
			Description: "This runs go install in the current directory.",
			Action:      burrow.Install,
		},
		{
			Name:        "uninstall",
			Aliases:     []string{"un", "uninst"},
			Flags:       []cli.Flag{},
			Usage:       "Uninstall the application from the GOPATH.",
			Description: "This run go clean -i in the current directory.",
			Action:      burrow.Uninstall,
		},
		{
			Name:        "package",
			Aliases:     []string{"pack"},
			Flags:       []cli.Flag{force_flag},
			Usage:       "Create a .tar.gz containing the binary.",
			Description: "This runs tar to package your application.",
			Action:      burrow.Package,
		},
		{
			Name:        "publish",
			Aliases:     []string{"pub"},
			Flags:       []cli.Flag{},
			Usage:       "Publish the current version by building a package and setting a version tag in git.",
			Description: "This runs git tag -f vX.Y.Z in the current directory. Any arguments following -- will be directly passed to git.",
			Action:      burrow.Publish,
		},
		{
			Name:        "clean",
			Aliases:     []string{},
			Flags:       []cli.Flag{},
			Usage:       "Clean the project from any build artifacts.",
			Description: "This runs go clean in the current directory and removes artifacts created by burrow.",
			Action:      burrow.Clean,
		},
		{
			Name:        "doc",
			Aliases:     []string{},
			Flags:       []cli.Flag{force_flag},
			Usage:       "Host the go documentation on this machine.",
			Description: "This runs go doc in the current directory. Any arguments following -- will be directly passed to go doc.",
			Action:      burrow.Doc,
		},
		{
			Name:        "format",
			Aliases:     []string{"fmt"},
			Flags:       []cli.Flag{force_flag},
			Usage:       "Format the code of this project with gofmt.",
			Description: "This runs gofmt in the current directory. Any arguments following -- will be directly passed to gofmt.",
			Action:      burrow.Format,
		},
		{
			Name:        "check",
			Aliases:     []string{},
			Flags:       []cli.Flag{force_flag},
			Usage:       "Check the code with go vet.",
			Description: "This runs go tool vet in the current directory. Any arguments following -- will be directly passed to go.",
			Action:      burrow.Check,
		},
	}

	app.Run(os.Args)
}
