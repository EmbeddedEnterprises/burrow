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

	"github.com/fatih/color"
)

type LogLevel int

const (
	LOG_INFO LogLevel = iota
	LOG_WARN
	LOG_ERR
)

func Log(level LogLevel, target string, format string, args ...interface{}) {
	switch level {
	case LOG_INFO:
		color.Set(color.FgWhite)
	case LOG_WARN:
		color.Set(color.FgYellow)
	case LOG_ERR:
		color.Set(color.FgRed)
	}
	fmt.Fprintf(os.Stderr, "[%10s] ", target)
	color.Unset()
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

// TODO: Add log writer for Stdout replacement
