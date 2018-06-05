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
	"path/filepath"
	"strings"

	"github.com/EmbeddedEnterprises/burrow/utils"
	"github.com/mattn/go-shellwords"
	"github.com/urfave/cli"
)

// Clone clones a git repository into the GOPATH and creates a symlink in the cwd when cwd is not in GOPATH.
func Clone(context *cli.Context) error {
	options := context.Args()

	if len(options) != 1 {
		cli.ShowCommandHelp(context, "clone")
		return nil
	}

	gopath, err := filepath.Abs(os.Getenv("GOPATH"))
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "clone", "Failed to get GOPATH!")
		return err
	}

	url := options[0]
	destination := options[0]

	if idx := strings.Index(destination, "://"); idx >= 0 {
		destination = destination[idx+3:]
	}

	if strings.HasSuffix(destination, ".git") {
		destination = destination[:len(destination)-4]
	}
	destination = gopath + "/src/" + destination

	tmp := strings.Split(destination, "/")
	link := tmp[len(tmp)-1]

	burrow.Log(burrow.LOG_INFO, "clone", "Cloning git repository into GOPATH...")

	args := []string{}
	args = append(args, "clone", url, destination)
	userArgs, err := shellwords.Parse(burrow.Config.Args.Git.Clone)
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "clone", "Failed to read user arguments from config file: %s", err)
		return err
	}
	args = append(args, userArgs...)
	if err := burrow.Exec("clone", "git", args...); err != nil {
		return err
	}

	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		burrow.Log(burrow.LOG_ERR, "clone", "Failed to get current directory!")
		return err
	}

	// only create symlink if not in GOPATH
	if !strings.HasPrefix(cwd, gopath) {
		if err = os.Symlink(destination, link); err != nil {
			burrow.Log(burrow.LOG_ERR, "clone", "Failed to create symlink!")
			return err
		}
	}

	return nil
}
