# burrow [![Latest Tag](https://img.shields.io/github/tag/EmbeddedEnterprises/burrow.svg)](https://github.com/EmbeddedEnterprises/burrow/releases) [![Build Status](https://travis-ci.org/EmbeddedEnterprises/burrow.svg?branch=master)](https://travis-ci.org/EmbeddedEnterprises/burrow) [![Go Report Card](https://goreportcard.com/badge/github.com/EmbeddedEnterprises/burrow)](https://goreportcard.com/report/github.com/EmbeddedEnterprises/burrow) [![GoDoc](https://godoc.org/github.com/EmbeddedEnterprises/burrow?status.svg)](https://godoc.org/github.com/EmbeddedEnterprises/burrow) [![Docker Pulls](https://img.shields.io/docker/pulls/embeddedenterprises/burrow.svg)](https://hub.docker.com/r/embeddedenterprises/burrow/) [![Docker Build Status](https://img.shields.io/docker/build/embeddedenterprises/burrow.svg)](https://hub.docker.com/r/embeddedenterprises/burrow/builds/)

> WARNING: This project got deprecated in favor of the official `go mod` tool!

Burrow is a go build system that used glide for dependency management, but now wraps the official `go mod` tool. Burrow tries to solve issues for creating reproducible builds with the `go tool`. Burrow is a wrapper around the `go tool` that enables the possibility to define default arguments for e.g. the `go build` command on a project basis. Additionally, burrow introduces a complete project lifecycle containing project creation, dependency management, building, installation, packaging, publication of version tags, documentation hosting, code formatting, and code checking.

However, every burrow project can still be built by issuing `go build` manually! Burrow is not needed to create a build of a burrow project. Burrow only simplifies the use of the `go tool` for easier build reproduction.

## But why another build system/tool for go?

![xkcd](https://imgs.xkcd.com/comics/standards.png)

So this is the next go build tool that makes the whole ecosystem even more complicated. But why are there so many go dependency management/build tools in the first place?

The `go tool` makes it really hard to let another developer reproduce your build. Glide helps but does not solve the problem of additional build parameters. A Makefile can be a solution for this but Makefiles are not that easy to update in a centralized manner. Also, the `go tool` does not model a development workflow for a project that integrates in the concept of reproducible builds. Burrow tries to solve this by providing a `publish` and a `package` command and using glide for dependency management.

## How to install burrow?

Do not use burrow for new projects! Every burrow project is buildable with the official go tooling.

## How to create a project?

Please do not use `burrow` for new projects. Instead use the official `go mod` tool for creating a new go project. Please refer to the [official documentation](https://github.com/golang/go/wiki/Modules) on go modules to create a new go project.

### The project layout of a burrow project

```
+-bin/
| +-app
+-example/
| +-api-showcase.go
+-package/
| +-app-0.1.0.tar.gz
+-vendor/
+-burrow.yaml
+-glide.lock
+-glide.yaml
+-LICENSE
+-main.go
+-README.md
```

## How to build?

To build a burrow application just run

```
$ burrow build
```

However, as `burrow` is now deprecated please use the official go tooling to build your projects now. You can run `burrow build` to see which `go build` commands are executed by `burrow` and integrate this into your own build process. Usually `burrow` runs

```
$ go build -o ./bin/<project-name> main.go
```

## Other commands

The below text can be shown by running `burrow --help`.

```
Usage: burrow [global options] command [command options] [arguments...]

A go build system that used glide for dependency management and is now deprecated in favor of 'go mod'.

Commands:
   init, create           Initializes a directory as a burrow project (deprecation stub).
   new                    Creates a new folder that contains an empty burrow project (deprecation stub).
   clone                  Clone a git repository into your current directory.
   get                    Add a dependency to the go.mod file.
   fetch, ensure, f, e    Get all dependencies from the go.sum file to reproduce a build.
   update, u, up          Update all dependencies from the go.mod file and update the go.sum file.
   run, r                 Run the application.
   test, t                Run all existing tests of the application.
   build, b               Build the application.
   install, i, in, inst   Install the application in the GOPATH.
   uninstall, un, uninst  Uninstall the application from the GOPATH.
   package, pack          Create a .tar.gz containing the binary.
   publish, pub           Publish the current version by building a package and setting a version tag in git.
   clean                  Clean the project from any build artifacts.
   doc                    Host the go documentation on this machine.
   format, fmt            Format the code of this project with 'go fmt'.
   check, vet             Check the code with 'go vet'.
   major                  Increment the major part of the version for this project.
   minor                  Increment the minor part of the version for this project.
   patch                  Increment the patch part of the version for this project.
   help, h                Shows a list of commands or help for one command

Global options:
   --help, -h     show help
   --version, -v  print the version
   
Authors:
   Fin Christensen <christensen.fin@gmail.com>
   
burrow - Copyright (c) 2017-2019  EmbeddedEnterprises
```

## License

This project is licensed under GPL-3.
