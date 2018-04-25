# burrow [![Latest Tag](https://img.shields.io/github/tag/EmbeddedEnterprises/burrow.svg)](https://github.com/EmbeddedEnterprises/burrow/releases) [![Build Status](https://travis-ci.org/EmbeddedEnterprises/burrow.svg?branch=master)](https://travis-ci.org/EmbeddedEnterprises/burrow) [![Go Report Card](https://goreportcard.com/badge/github.com/EmbeddedEnterprises/burrow)](https://goreportcard.com/report/github.com/EmbeddedEnterprises/burrow) [![GoDoc](https://godoc.org/github.com/EmbeddedEnterprises/burrow?status.svg)](https://godoc.org/github.com/EmbeddedEnterprises/burrow)
[![Docker Pulls](https://img.shields.io/docker/pulls/embeddedenterprises/burrow.svg)](https://hub.docker.com/r/embeddedenterprises/burrow/)
[![Docker Build Status](https://img.shields.io/docker/build/embeddedenterprises/burrow.svg)](https://hub.docker.com/r/embeddedenterprises/burrow/builds/)

Burrow is a go build system that uses glide for dependency management. Burrow tries to solve issues for creating reproducible builds with the `go tool`. Burrow is a wrapper around the `go tool` that enables the possibility to define default arguments for e.g. the `go build` command on a project basis. Additionally, burrow introduces a complete project lifecycle containing project creation, dependency management, building, installation, packaging, publication of version tags, documentation hosting, code formatting, and code checking.

However, every burrow project can still be built by issuing `go build` manually! Burrow is not needed to create a build of a burrow project. Burrow only simplifies the use of the `go tool` for easier build reproduction.

## But why another build system/tool for go?

![xkcd](https://imgs.xkcd.com/comics/standards.png)

So this is the next go build tool that makes the whole ecosystem even more complicated. But why are there so many go dependency management/build tools in the first place?

The `go tool` makes it really hard to let another developer reproduce your build. Glide helps but does not solve the problem of additional build parameters. A Makefile can be a solution for this but Makefiles are not that easy to update in a centralized manner. Also, the `go tool` does not model a development workflow for a project that integrates in the concept of reproducible builds. Burrow tries to solve this by providing a `publish` and a `package` command and using glide for dependency management.

## How to install burrow?

Install `glide` from [here](https://github.com/Masterminds/glide).

```
$ go get github.com/EmbeddedEnterprises/burrow
```

## How to create a project?

To create a burrow project just create an empty folder inside your `GOPATH` and run

```
$ cd "$GOPATH"/src/github.com/EmbeddedEnterprises/
$ mkdir test
$ burrow init
```

and you will be guided through the project setup. You can edit the created `burrow.yaml` manually to enter additional parameters for the `go tool` commands. When you are creating a binary project a `main.go` will be created. Otherwise a `lib.go` will be created.

You may access the project outside the `GOPATH` by symlinking to the project folder in your `GOPATH`. `burrow` will detect the symlink let you use it like if you were inside the `GOPATH`.

```
$ ln -s "$GOPATH"/src/github.com/EmbeddedEnterprises/test ~/Development/github/test
$ cd ~/Development/github.com/test
$ burrow build
```

### What about go generate?

As I currently do not use `go generate` and do not know how it would integrate in the burrow workflow, the generate command is currently not supported by burrow. Feel free to add it by yourself and make a pull request. You may also open an issue on this topic and discuss implementation approaches.

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

## Other commands

The below text can be shown by running `burrow --help`.

```
Usage: burrow [global options] command [command options] [arguments...]

A go build system that uses glide for dependency management.

Commands:
   init, create           Create a new burrow project.
   clone                  Clone a git repository into your GOPATH and create a symbolic link in your current location.
   get                    Install a dependency in the vendor folder and add it to the glide yaml.
   fetch, ensure, f, e    Get all dependencies from the lock file to reproduce a build.
   update, u, up          Update all dependencies from the yaml file and update the lock file.
   run, r                 Build and run the application.
   test, t                Run all existing tests of the application.
   build, b               Build the application.
   install, i, in, inst   Install the application in the GOPATH.
   uninstall, un, uninst  Uninstall the application from the GOPATH.
   package, pack          Create a .tar.gz containing the binary.
   publish, pub           Publish the current version by building a package and setting a version tag in git.
   clean                  Clean the project from any build artifacts.
   doc                    Host the go documentation on this machine.
   format, fmt            Format the code of this project with gofmt.
   check                  Check the code with go vet.
   help, h                Shows a list of commands or help for one command

Global options:
   --help, -h     show help
   --version, -v  print the version
   
Authors:
   Fin Christensen <christensen.fin@gmail.com>
   
burrow - Copyright (c) 2017  EmbeddedEnterprises
```

## License

This project is licensed under GPL-3.
