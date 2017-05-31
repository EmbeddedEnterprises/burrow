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
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
)

var targetState = map[string]bool{}
var projectHash string = ""

func IsTargetUpToDate(target string, outputs []string) bool {
	LoadConfig()

	isTargetUpToDate, ok := targetState[target]
	if ok {
		return isTargetUpToDate
	}

	if projectHash == "" {
		project_dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		hash_source := Config.Name + "_" + project_dir
		sha1_hasher := sha1.New()
		sha1_hasher.Write([]byte(hash_source))
		projectHash = base64.URLEncoding.EncodeToString(sha1_hasher.Sum(nil))
	}

	usr, _ := user.Current()
	_ = os.Mkdir(usr.HomeDir+"/.cache/burrow/", 0755)
	_ = os.Mkdir(usr.HomeDir+"/.cache/burrow/"+projectHash, 0755)

	cache, err := ioutil.ReadFile(usr.HomeDir + "/.cache/burrow/" + projectHash + "/" + target)
	if err != nil {
		targetState[target] = false
		return false
	}

	code_files := GetCodefilesWithMtime(outputs)

	cached_code_files := map[string]int64{}
	err = yaml.Unmarshal(cache, &cached_code_files)
	if err != nil {
		targetState[target] = false
		return false
	}

	for path, mtime := range code_files {
		cached_mtime, ok := cached_code_files[path]
		if !ok || cached_mtime < mtime {
			targetState[target] = false
			return false
		}
	}

	targetState[target] = true
	return true
}

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

func GetCodefiles() []string {
	code_files := []string{}
	_ = filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") && !strings.Contains(path, "vendor/") {
			code_files = append(code_files, path)
		}
		return nil
	})
	return code_files
}

func GetCodefilesWithMtime(outputs []string) map[string]int64 {
	code_files := map[string]int64{}
	_ = filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".yaml") {
			code_files[path] = f.ModTime().Unix()
		}
		return nil
	})

	for _, output := range outputs {
		info, err := os.Stat(output)
		if err != nil {
			code_files[output] = math.MaxInt64
		} else {
			code_files[output] = info.ModTime().Unix()
		}
	}

	return code_files
}