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
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/go-yaml/yaml"
	"github.com/mattn/go-shellwords"
	"github.com/urfave/cli"
)

const (
	EXIT_SUCCESS int = iota
	EXIT_CONFIG
	EXIT_ACTION
)

type Config struct {
	Name    string
	Version string
	Package struct {
		Include []string
	}
	Args struct {
		Run string
		Go  struct {
			Test  string
			Build string
			Doc   string
			Vet   string
			Fmt   string
		}
		Glide struct {
			Install string
			Update  string
			Get     string
		}
		Git struct {
			Tag string
		}
	}
}

type LogLevel int

const (
	LOG_INFO LogLevel = iota
	LOG_WARN
	LOG_ERR
)

func log(level LogLevel, target string, format string, args ...interface{}) {
	switch level {
	case LOG_INFO:
		color.Set(color.FgBlue)
	case LOG_WARN:
		color.Set(color.FgYellow)
	case LOG_ERR:
		color.Set(color.FgRed)
	}
	fmt.Fprintf(os.Stderr, "[%10s] ", target)
	color.Unset()
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

var config Config = Config{}
var isConfigLoaded bool = false

func LoadConfig() {
	if isConfigLoaded {
		return
	}

	data, err := ioutil.ReadFile("burrow.yaml")

	if err != nil {
		log(LOG_ERR, "burrow", "Not a burrow project!")
		os.Exit(EXIT_CONFIG)
	}

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		log(LOG_ERR, "burrow", "Failed to read burrow config: %v", err)
		os.Exit(EXIT_CONFIG)
	}

	isConfigLoaded = true
}

var targetState = map[string]bool{}
var projectHash string = ""

func GetCodefilesWithMtime() map[string]int64 {
	code_files := map[string]int64{}
	_ = filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".yaml") {
			code_files[path] = f.ModTime().Unix()
		}
		return nil
	})
	return code_files
}

func IsTargetUpToDate(target string) bool {
	LoadConfig()

	isTargetUpToDate, ok := targetState[target]
	if ok {
		return isTargetUpToDate
	}

	if projectHash == "" {
		project_dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		hash_source := config.Name + "_" + project_dir
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

	code_files := GetCodefilesWithMtime()

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

func UpdateTarget(target string) {
	cache := GetCodefilesWithMtime()
	ser, err := yaml.Marshal(&cache)
	if err != nil {
		log(LOG_WARN, target, "Failed to update target cache, ignoring...")
		return
	}
	usr, err := user.Current()
	if err != nil {
		log(LOG_WARN, target, "Failed to update target cache: $s", err)
		return
	}
	err = ioutil.WriteFile(usr.HomeDir+"/.cache/burrow/"+projectHash+"/"+target, ser, 0644)
	if err != nil {
		log(LOG_WARN, target, "Failed to update target cache, ignoring...")
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

func command(comm string, args ...string) error {
	cmd := exec.Command(comm, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return cli.NewExitError(fmt.Sprintf("Error running action: %v", err), EXIT_ACTION)
	}
	return nil
}

func create(context *cli.Context) error {
	return nil
}

func get(context *cli.Context) error {
	LoadConfig()

	if len(context.Args()) != 1 {
		cli.ShowCommandHelp(context, "get")
		return nil
	}

	dep := context.Args()[0]
	log(LOG_INFO, "get", "Adding new dependency %s", dep)

	args := []string{}
	args = append(args, "get")
	user_args, err := shellwords.Parse(config.Args.Glide.Get)
	if err != nil {
		log(LOG_ERR, "get", "Failed to read user arguments from config file: %s", err)
		return nil
	}
	args = append(args, user_args...)
	args = append(args, dep)

	return command("glide", args...)
}

func fetch(context *cli.Context) error {
	LoadConfig()
	log(LOG_INFO, "fetch", "Fetching dependencies from lock file")
	gopath := os.Getenv("GOPATH")
	args := []string{}
	args = append(args, "install")
	user_args, err := shellwords.Parse(config.Args.Glide.Install)
	if err != nil {
		log(LOG_ERR, "fetch", "Failed to read user arguments from config file: %s", err)
		return nil
	}
	args = append(args, user_args...)
	return command(gopath+"/bin/glide", args...)
}

func update(context *cli.Context) error {
	LoadConfig()
	log(LOG_INFO, "update", "Updating dependencies from glide yaml config")
	gopath := os.Getenv("GOPATH")
	args := []string{}
	args = append(args, "update")
	user_args, err := shellwords.Parse(config.Args.Glide.Update)
	if err != nil {
		log(LOG_ERR, "update", "Failed to read user arguments from config file: %s", err)
		return nil
	}
	args = append(args, user_args...)
	return command(gopath+"/bin/glide", args...)
}

func run(context *cli.Context) error {
	LoadConfig()

	if err := build(context); err != nil {
		return err
	}

	log(LOG_INFO, "run", "Running project")

	// add --example and -- for user args

	args := []string{}
	user_args, err := shellwords.Parse(config.Args.Run)
	if err != nil {
		log(LOG_ERR, "run", "Failed to read user arguments from config file: %s", err)
		return nil
	}
	args = append(args, user_args...)
	args = append(args, context.Args()...)
	return command("./bin/"+config.Name, args...)
}

func test(context *cli.Context) error {
	LoadConfig()

	if IsTargetUpToDate("test") {
		log(LOG_INFO, "test", "Tests are up-to-date")
		return nil
	}

	log(LOG_INFO, "test", "Running tests for project")

	args := []string{}
	args = append(args, "test")
	user_args, err := shellwords.Parse(config.Args.Go.Test)
	if err != nil {
		log(LOG_ERR, "test", "Failed to read user arguments from config file: %s", err)
		return nil
	}
	args = append(args, user_args...)
	err = command("go", args...)
	if err == nil {
		UpdateTarget("test")
	}
	return err
}

func build(context *cli.Context) error {
	LoadConfig()

	if IsTargetUpToDate("build") {
		log(LOG_INFO, "build", "Build is up-to-date")
		return nil
	}
	log(LOG_INFO, "build", "Building project")

	_ = os.Mkdir("./bin", 0755)

	args := []string{}
	args = append(args, "build", "-o", "./bin/"+config.Name)
	user_args, err := shellwords.Parse(config.Args.Go.Build)
	if err != nil {
		log(LOG_ERR, "build", "Failed to read user arguments from config file: %s", err)
		return nil
	}
	args = append(args, user_args...)
	err = command("go", args...)
	if err == nil {
		UpdateTarget("build")
	}
	return err
}

func install(context *cli.Context) error {
	LoadConfig()
	if err := format(context); err != nil {
		return err
	}
	if err := check(context); err != nil {
		return err
	}
	if err := test(context); err != nil {
		return err
	}
	if err := build(context); err != nil {
		return err
	}

	if IsTargetUpToDate("install") {
		log(LOG_INFO, "install", "Installation is up-to-date")
		return nil
	}
	log(LOG_INFO, "install", "Installing application in GOPATH")

	args := []string{}
	args = append(args, "install")
	user_args, err := shellwords.Parse(config.Args.Go.Build)
	if err != nil {
		log(LOG_ERR, "install", "Failed to read user arguments from config file: %s", err)
		return nil
	}
	args = append(args, user_args...)
	err = command("go", args...)
	if err == nil {
		UpdateTarget("install")
	}
	return err
}

func uninstall(context *cli.Context) error {
	LoadConfig()
	log(LOG_INFO, "uninstall", "Uninstalling application from GOPATH")
	return command("go", "clean", "-i")
}

func pack(context *cli.Context) error {
	LoadConfig()
	_ = os.Mkdir("./package", 0755)
	if err := format(context); err != nil {
		return err
	}
	if err := check(context); err != nil {
		return err
	}
	if err := test(context); err != nil {
		return err
	}
	if err := doc(context); err != nil {
		return err
	}
	if err := build(context); err != nil {
		return err
	}

	if IsTargetUpToDate("package") {
		log(LOG_INFO, "package", "Package is up-to-date")
		return nil
	}

	log(LOG_INFO, "package", "Packaging project")

	args := []string{}
	args = append(args,
		"czf",
		fmt.Sprintf("./package/%s-%s.tar.gz", config.Name, config.Version),
		"./bin/"+config.Name,
	)
	args = append(args, config.Package.Include...)

	err := command("tar", args...)

	if err == nil {
		UpdateTarget("package")
	}
	return err
}

func publish(context *cli.Context) error {
	LoadConfig()
	if err := pack(context); err != nil {
		return err
	}
	log(LOG_INFO, "publish", "Publishing new version tag in git")

	err := command("git", "diff-index", "--quiet", "HEAD", "--")
	if err != nil {
		log(LOG_ERR, "publish", "You have unstaged changes, commit them to proceed!")
		return cli.NewExitError("", EXIT_ACTION)
	}

	args := []string{}
	args = append(args, "tag", "-f")
	user_args, err := shellwords.Parse(config.Args.Git.Tag)
	if err != nil {
		log(LOG_ERR, "publish", "Failed to read user arguments from config file: %s", err)
		return nil
	}
	args = append(args, user_args...)
	args = append(args, config.Version)
	return command("git", args...)
}

func clean(context *cli.Context) error {
	LoadConfig()
	log(LOG_INFO, "clean", "Cleaning build artifacts")
	if err := command("go", "clean"); err != nil {
		return err
	}
	if err := command("rm", "./bin/"+config.Name); err != nil {
		return err
	}
	return command("rm", fmt.Sprintf("./package/%s-%s.tar.gz", config.Name, config.Version))
}

func doc(context *cli.Context) error {
	LoadConfig()

	if IsTargetUpToDate("doc") {
		log(LOG_INFO, "doc", "Documentation is up-to-date")
		return nil
	}

	log(LOG_INFO, "doc", "Building documentation")

	args := []string{}
	args = append(args, "doc")
	user_args, err := shellwords.Parse(config.Args.Go.Doc)
	if err != nil {
		log(LOG_ERR, "doc", "Failed to read user arguments from config file: %s", err)
		return nil
	}
	args = append(args, user_args...)
	err = command("go", args...)
	if err == nil {
		UpdateTarget("doc")
	}
	return err
}

func format(context *cli.Context) error {
	LoadConfig()

	if IsTargetUpToDate("format") {
		log(LOG_INFO, "format", "Code formatting is up-to-date")
		return nil
	}

	log(LOG_INFO, "format", "Formatting code")

	args := []string{}
	args = append(args, "-l", "-w")
	user_args, err := shellwords.Parse(config.Args.Go.Fmt)
	if err != nil {
		log(LOG_ERR, "format", "Failed to read user arguments from config file: %s", err)
		return nil
	}
	args = append(args, user_args...)
	args = append(args, GetCodefiles()...)
	err = command("gofmt", args...)
	if err == nil {
		UpdateTarget("format")
	}
	return err
}

func check(context *cli.Context) error {
	LoadConfig()

	if IsTargetUpToDate("check") {
		log(LOG_INFO, "check", "Code has already been checked")
		return nil
	}

	log(LOG_INFO, "check", "Checking code")

	args := []string{}
	args = append(args, "tool", "vet")
	user_args, err := shellwords.Parse(config.Args.Go.Vet)
	if err != nil {
		log(LOG_ERR, "check", "Failed to read user arguments from config file: %s", err)
		return nil
	}
	args = append(args, user_args...)
	args = append(args, GetCodefiles()...)
	err = command("go", args...)
	if err == nil {
		UpdateTarget("check")
	}
	return err
}

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
			Action:  create,
		},
		{
			Name:    "get",
			Aliases: []string{},
			Flags:   []cli.Flag{},
			Usage:   "Install a dependency in the vendor folder and add it to the glide yaml",
			Action:  get,
		},
		{
			Name:    "fetch",
			Aliases: []string{"ensure", "f", "e"},
			Flags:   []cli.Flag{},
			Usage:   "Get all dependencies from the lock file to reproduce a build",
			Action:  fetch,
		},
		{
			Name:    "update",
			Aliases: []string{"u", "up"},
			Flags:   []cli.Flag{},
			Usage:   "Update all dependencies from the yaml file and update the lock file",
			Action:  update,
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Flags:   []cli.Flag{},
			Usage:   "Build and run the application",
			Action:  run,
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Flags:   []cli.Flag{},
			Usage:   "Run all existing tests of the application",
			Action:  test,
		},
		{
			Name:    "build",
			Aliases: []string{"b"},
			Flags:   []cli.Flag{},
			Usage:   "Build the application",
			Action:  build,
		},
		{
			Name:    "install",
			Aliases: []string{"i", "in", "inst"},
			Flags:   []cli.Flag{},
			Usage:   "Install the application in the GOPATH",
			Action:  install,
		},
		{
			Name:    "uninstall",
			Aliases: []string{"un", "uninst"},
			Flags:   []cli.Flag{},
			Usage:   "Uninstall the application from the GOPATH",
			Action:  uninstall,
		},
		{
			Name:    "package",
			Aliases: []string{"pack"},
			Flags:   []cli.Flag{},
			Usage:   "Create a .tar.gz containing the binary",
			Action:  pack,
		},
		{
			Name:    "publish",
			Aliases: []string{"pub"},
			Flags:   []cli.Flag{},
			Usage:   "Publish the current version by building a package and setting a version tag in git",
			Action:  publish,
		},
		{
			Name:    "clean",
			Aliases: []string{},
			Flags:   []cli.Flag{},
			Usage:   "Clean the project from any build artifacts",
			Action:  clean,
		},
		{
			Name:    "doc",
			Aliases: []string{},
			Flags:   []cli.Flag{},
			Usage:   "Generate the godoc documentation for this project",
			Action:  doc,
		},
		{
			Name:    "format",
			Aliases: []string{"fmt"},
			Flags:   []cli.Flag{},
			Usage:   "Format the code of this project with gofmt",
			Action:  format,
		},
		{
			Name:    "check",
			Aliases: []string{},
			Flags:   []cli.Flag{},
			Usage:   "Check the code with go vet",
			Action:  check,
		},
	}

	app.Run(os.Args)
}
