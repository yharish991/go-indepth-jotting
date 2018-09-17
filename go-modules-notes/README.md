## Go Modules Introduced in 1.11

1. A module is a collection of related Go Packages.
2. Modules replace the OLD GOPATH based approach to specifying which source files are used in a given build.

In module-aware mode, GOPATH no longer defines the meaning of imports during a build.

The `module path` is the import path prefix corresponding to the module root.

The `go.mod` file defines the module path.

To start a new module, create a go.mod file in the root of the module's directory tree.

`go mod init example.com/m`

This declares that the directory containing it is the root of the module with path `example.com/m`
