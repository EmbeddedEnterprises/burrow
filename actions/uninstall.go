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

// Uninstall uninstalls the burrow application from the GOPATH.
func Uninstall(context *cli.Context) error {
	burrow.LoadConfig()
	burrow.Log(burrow.LOG_INFO, "uninstall", "Uninstalling application from GOPATH")
	err := burrow.Exec("uninstall", "go", "clean", "-i")

	burrow.Deprecation("uninstall", []string{"go", "clean", "-i"})

	return err
}
