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
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/urfave/cli"
)

var targetState = map[string]bool{}
var projectHash string

// IsTargetUpToDate checks whether a given build target is up-to-date. This means that all build
// artifacts of the target were created from data with the same timestamp as the currently
// available sources.
//
// The outputs parameter specifies which files are created (artifacts) by the target. If these
// files are not available all cache data is invalid.
func IsTargetUpToDate(target string, outputs []string) bool {
	LoadConfig()

	isTargetUpToDate, ok := targetState[target]
	if ok {
		return isTargetUpToDate
	}

	if projectHash == "" {
		projectDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		hashSource := Config.Name + "_" + projectDir
		sha1Hasher := sha1.New()
		sha1Hasher.Write([]byte(hashSource))
		projectHash = base64.URLEncoding.EncodeToString(sha1Hasher.Sum(nil))
	}

	usr, _ := user.Current()
	_ = os.MkdirAll(usr.HomeDir+"/.cache/burrow/", 0755)
	_ = os.MkdirAll(usr.HomeDir+"/.cache/burrow/"+projectHash, 0755)

	cache, err := ioutil.ReadFile(usr.HomeDir + "/.cache/burrow/" + projectHash + "/" + target)
	if err != nil {
		targetState[target] = false
		return false
	}

	codeFiles := GetCodefilesWithMtime(outputs)

	cachedCodeFiles := map[string]int64{}
	err = yaml.Unmarshal(cache, &cachedCodeFiles)
	if err != nil {
		targetState[target] = false
		return false
	}

	for path, mtime := range codeFiles {
		cachedMtime, ok := cachedCodeFiles[path]
		if !ok || cachedMtime < mtime {
			targetState[target] = false
			return false
		}
	}

	targetState[target] = true
	return true
}

// UpdateTarget updates the cache of a target to match the timestamps of all currently available sources.
// The outputs parameter specifies which files are created (artifacts) by the target. Timestamps of the
// artifacts will also be stored.
func UpdateTarget(target string, outputs []string) {
	cache := GetCodefilesWithMtime(outputs)
	ser, err := yaml.Marshal(&cache)
	if err != nil {
		Log(LOG_WARN, target, "Failed to update target cache, ignoring...")
		return
	}
	usr, err := user.Current()
	if err != nil {
		Log(LOG_WARN, target, "Failed to update target cache: $s", err)
		return
	}
	err = ioutil.WriteFile(usr.HomeDir+"/.cache/burrow/"+projectHash+"/"+target, ser, 0644)
	if err != nil {
		Log(LOG_WARN, target, "Failed to update target cache, ignoring...")
		return
	}
}

// GetCodefiles returns a string array containing all paths of files that contain code inside the
// current burrow project.
func GetCodefiles() []string {
	codeFiles := []string{}
	_ = filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") && !strings.Contains(path, "vendor/") {
			codeFiles = append(codeFiles, path)
		}
		return nil
	})
	return codeFiles
}

// GetCodefilesWithMtime returns a map containing all paths of files that contain code inside the
// current burrow project. The paths get mapped to Unix modification times. The outputs parameter
// should contain additional non-code files that should also be contained in the map.
func GetCodefilesWithMtime(outputs []string) map[string]int64 {
	codeFiles := map[string]int64{}
	_ = filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".yaml") {
			codeFiles[path] = f.ModTime().Unix()
		}
		return nil
	})

	for _, output := range outputs {
		info, err := os.Stat(output)
		if err != nil {
			codeFiles[output] = math.MaxInt64
		} else {
			codeFiles[output] = info.ModTime().Unix()
		}
	}

	return codeFiles
}

// GetSecondLevelArgs returns the command line arguments that are located after a double dash (--).
func GetSecondLevelArgs() cli.Args {
	args := os.Args
	second := cli.Args{}

	doubleDashFound := false
	for _, val := range args {
		if val == "--" {
			doubleDashFound = true
		} else if doubleDashFound {
			second = append(second, val)
		}
	}

	return second
}

// WrapAction sets the useSecondLevelArgs of an action to true by default.
func WrapAction(action func(*cli.Context, bool) error) func(*cli.Context) error {
	return func(c *cli.Context) error {
		return action(c, true)
	}
}
