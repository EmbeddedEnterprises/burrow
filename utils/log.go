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
	"strings"

	"github.com/fatih/color"
)

// The LogLevel type describes the log level of a log message.
type LogLevel int

// The LOG_... constants describe log level of a log message.
const (
	LOG_INFO LogLevel = iota
	LOG_WARN
	LOG_ERR
)

type targetString = string
type commandString = string

var tuple = []struct {
	string
	int
}{}

var deprecationCommands = []struct {
	targetString
	commandString
}{}

// Deprecation marks a target for deprecation and provide underlying commands as
// args to the deprecation log.
func Deprecation(target string, args ...[]string) {
	for _, command := range args {
		deprecationCommands = append(
			deprecationCommands,
			struct {
				targetString
				commandString
			}{target, strings.Join(command, " ")},
		)
	}
}

// LogDeprecationMessage logs all accumulated deprecation messages to the console.
func LogDeprecationMessage() {
	target := "burrow"

	Log(LOG_WARN, target, "")
	Log(LOG_WARN, target, "The burrow tool got deprecated in favor of the official 'go mod' tool!")

	if len(deprecationCommands) > 0 {
		Log(LOG_WARN, target, "Use the following commands instead now:")

		for _, command := range deprecationCommands {
			Log(LOG_WARN, command.targetString, "    %s", command.commandString)
		}
	}

	Log(LOG_WARN, target, "")

}

// Log writes a log message to stderr which will be colored by the given log level and prefixed
// with the given target. The format parameter is a fmt.Printf parameter following the args for
// formatting.
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

// The Logger struct holds a logger instance for output redirection.
type Logger struct {
	target string
	level  LogLevel
}

// NewLogger creates a new logger instance that can be used to redirect streams and produce log
// message.
func NewLogger(target string, level LogLevel) Logger {
	return Logger{
		target: target,
		level:  level,
	}
}

func (log Logger) Write(payload []byte) (int, error) {
	message := string(payload[:])
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		if line != "" {
			Log(log.level, log.target, line)
		}
	}
	return len(payload), nil
}
