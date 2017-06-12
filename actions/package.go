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
	"fmt"
	"os"
	"path/filepath"

	"github.com/EmbeddedEnterprises/burrow/utils"
	"github.com/urfave/cli"
)

// Package creates a .tar.gz containing the binary.
func Package(context *cli.Context) error {
	burrow.LoadConfig()
	_ = os.Mkdir("./package", 0755)
	if err := Format(context); err != nil {
		return err
	}
	if err := Check(context); err != nil {
		return err
	}
	if err := Test(context); err != nil {
		return err
	}
	if err := Build(context); err != nil {
		return err
	}

	outputs := []string{fmt.Sprintf("./package/%s-%s.tar.gz", burrow.Config.Name, burrow.Config.Version)}

	if burrow.IsTargetUpToDate("package", outputs) && !context.Bool("force") {
		burrow.Log(burrow.LOG_INFO, "package", "Package is up-to-date")
		return nil
	}

	burrow.Log(burrow.LOG_INFO, "package", "Packaging project")

	args := []string{}
	args = append(args, "czf", outputs[0])
	_ = filepath.Walk("./bin", func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			args = append(args, path)
		}
		return nil
	})
	args = append(args, burrow.Config.Package.Include...)

	err := burrow.Exec("package", "tar", args...)
	if err == nil {
		burrow.UpdateTarget("package", outputs)
	}

	return err
}
