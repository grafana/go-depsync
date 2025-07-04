# `go-depsync`

`go-depsync` is a small command that identifies common dependencies with a given package, referred to as parent, and outputs a `go get` command that aligns the version of common dependencies with those of the parent package.

It is useful for extensions and plugins that are built together with a core package, and whose dependencies need to be aligned for binary compatibility.

## Usage

`go-depsync` can be installed using `go install`:

```console
$ go install github.com/grafana/go-depsync@latest
```

After that, it is ready to be used:

```bash
go-depsync --parent go.k6.io/k6
```

In some cases a package may not directly import the parent but still needs to maintain compatibility with the parent package.  This is more useful in the case of go plugin development:

```bash
go-depsync --parent go.k6.io/k6@v1.0.0-rc1
```

If the `go.mod` file for the local package is not on the working directory, a path to it can also be specified:

```bash
go-depsync --gomod /path/to/go.mod --parent go.k6.io/k6
```

`go-depsync` produces an output similar to the following:

```console
$ go-depsync --parent=go.k6.io/k6
2023/11/17 15:02:56 Found parent go.k6.io/k6@v0.46.0
  DEPENDENCY              CURRENT VERSION  NEW VERSION
  google.golang.org/grpc  v1.57.0          v1.56.1
  github.com/spf13/cobra  v1.5.0           v1.4.0
  golang.org/x/sys        v0.11.0          v0.9.0
  github.com/spf13/afero  v1.2.2           v1.1.2

go get google.golang.org/grpc@v1.56.1 github.com/spf13/cobra@v1.4.0 golang.org/x/sys@v0.9.0 github.com/spf13/afero@v1.1.2
```

The final line includes the `go get` command that, when run, will sync the versions of commons dependencies to those of the parent. `go-depsync` outputs this line to `stdout` so it can be piped to a shell, or redirected to a script for later use.

## Trivia

Go versions earlier than 1.21 have been observed to produce unexpected result when the `go get` command is run, sometimes ignoring some of the versions specified in the command, or unexpectedly upgrading/downgrading dependencies that are not present on it. It is recommended to run the `go get` commands suggested by `go-depsync` with Go >= 1.21.
