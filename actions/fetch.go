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

// Fetch gets all dependencies from the glide lock file for reproducible builds.
func Fetch(context *cli.Context, useSecondLevelArgs bool) error {
	burrow.LoadConfig()

	burrow.Deprecation("fetch")
	burrow.Log(burrow.LOG_WARN, "fetch", "This command is not needed with the official 'go mod'! Dependencies are")
	burrow.Log(burrow.LOG_WARN, "fetch", "automatically fetched from the go.sum file with every build.")

	return nil
}
